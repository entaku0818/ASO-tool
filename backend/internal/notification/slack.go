package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
)

// SlackNotifier sends notifications to Slack via webhook
type SlackNotifier struct {
	webhookURL string
	httpClient *http.Client
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Blocks      []SlackBlock `json:"blocks,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// SlackBlock represents a Slack block element
type SlackBlock struct {
	Type     string      `json:"type"`
	Text     *TextObject `json:"text,omitempty"`
	Elements []Element   `json:"elements,omitempty"`
}

// TextObject represents a text object in Slack blocks
type TextObject struct {
	Type  string `json:"type"` // "plain_text" or "mrkdwn"
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

// Element represents an element in Slack blocks
type Element struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Attachment represents a Slack attachment
type Attachment struct {
	Color  string `json:"color,omitempty"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}

// Field represents a field in Slack attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SendBatchResult sends a batch job result notification to Slack
func (s *SlackNotifier) SendBatchResult(ctx context.Context, result *model.BatchResult) error {
	if s.webhookURL == "" {
		return nil // Skip if webhook URL is not configured
	}

	msg := s.formatBatchResultMessage(result)
	return s.send(ctx, msg)
}

// SendRankingChanges sends ranking change notifications to Slack
func (s *SlackNotifier) SendRankingChanges(ctx context.Context, changes []*model.RankingChange) error {
	if s.webhookURL == "" || len(changes) == 0 {
		return nil
	}

	msg := s.formatRankingChangesMessage(changes)
	return s.send(ctx, msg)
}

func (s *SlackNotifier) formatBatchResultMessage(result *model.BatchResult) *SlackMessage {
	statusEmoji := ":white_check_mark:"
	statusColor := "#36a64f" // green
	statusText := "成功"

	if !result.Success {
		statusEmoji = ":x:"
		statusColor = "#dc3545" // red
		statusText = "失敗"
	}

	// Build blocks
	blocks := []SlackBlock{
		{
			Type: "header",
			Text: &TextObject{
				Type:  "plain_text",
				Text:  fmt.Sprintf("%s ASO Tool バッチ実行結果", statusEmoji),
				Emoji: true,
			},
		},
		{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*ジョブ名:* `%s`\n*ステータス:* %s\n*実行時間:* %s",
					result.JobName, statusText, result.Duration.Round(time.Second)),
			},
		},
		{
			Type: "divider",
		},
	}

	// Add stats
	statsText := fmt.Sprintf(
		"*処理したアプリ数:* %d\n*更新したキーワード数:* %d\n*検知した変動数:* %d",
		result.AppsProcessed,
		result.KeywordsUpdated,
		len(result.Changes),
	)

	blocks = append(blocks, SlackBlock{
		Type: "section",
		Text: &TextObject{
			Type: "mrkdwn",
			Text: statsText,
		},
	})

	// Add errors if any
	if len(result.Errors) > 0 {
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf(":warning: *エラー (%d件):*\n```%s```",
					len(result.Errors),
					strings.Join(result.Errors, "\n")),
			},
		})
	}

	// Add significant ranking changes summary
	significantChanges := filterSignificantChanges(result.Changes)
	if len(significantChanges) > 0 {
		blocks = append(blocks, SlackBlock{
			Type: "divider",
		})
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf(":chart_with_upwards_trend: *重要な順位変動 (%d件):*", len(significantChanges)),
			},
		})

		// Show top 5 changes
		changeText := formatTopChanges(significantChanges, 5)
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: changeText,
			},
		})
	}

	// Add timestamp
	blocks = append(blocks, SlackBlock{
		Type: "context",
		Elements: []Element{
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("実行完了: %s", result.CompletedAt.Format("2006-01-02 15:04:05 JST")),
			},
		},
	})

	return &SlackMessage{
		Blocks: blocks,
		Attachments: []Attachment{
			{
				Color: statusColor,
			},
		},
	}
}

func (s *SlackNotifier) formatRankingChangesMessage(changes []*model.RankingChange) *SlackMessage {
	blocks := []SlackBlock{
		{
			Type: "header",
			Text: &TextObject{
				Type:  "plain_text",
				Text:  ":bell: ランキング変動アラート",
				Emoji: true,
			},
		},
	}

	// Group changes by app
	changesByApp := make(map[string][]*model.RankingChange)
	for _, change := range changes {
		changesByApp[change.AppName] = append(changesByApp[change.AppName], change)
	}

	for appName, appChanges := range changesByApp {
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*%s*", appName),
			},
		})

		changeText := formatChangesForApp(appChanges)
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: changeText,
			},
		})
	}

	return &SlackMessage{
		Blocks: blocks,
	}
}

func filterSignificantChanges(changes []*model.RankingChange) []*model.RankingChange {
	var significant []*model.RankingChange
	for _, c := range changes {
		if c.IsSignificant {
			significant = append(significant, c)
		}
	}
	return significant
}

func formatTopChanges(changes []*model.RankingChange, limit int) string {
	var lines []string
	count := 0
	for _, c := range changes {
		if count >= limit {
			break
		}

		emoji := ":arrow_up:"
		if c.Change < 0 {
			emoji = ":arrow_down:"
		}

		changeType := c.GetChangeType()
		var rankInfo string
		switch changeType {
		case model.RankingChangeTypeNew:
			rankInfo = fmt.Sprintf("圏外 -> %d位", *c.NewRank)
		case model.RankingChangeTypeLost:
			rankInfo = fmt.Sprintf("%d位 -> 圏外", *c.OldRank)
		default:
			if c.OldRank != nil && c.NewRank != nil {
				rankInfo = fmt.Sprintf("%d位 -> %d位 (%+d)", *c.OldRank, *c.NewRank, c.Change)
			}
		}

		line := fmt.Sprintf("%s *%s* (%s): %s", emoji, c.Keyword, c.Country, rankInfo)
		lines = append(lines, line)
		count++
	}

	if len(changes) > limit {
		lines = append(lines, fmt.Sprintf("...他 %d件", len(changes)-limit))
	}

	return strings.Join(lines, "\n")
}

func formatChangesForApp(changes []*model.RankingChange) string {
	var lines []string
	for _, c := range changes {
		emoji := ":arrow_up:"
		if c.Change < 0 {
			emoji = ":arrow_down:"
		}

		changeType := c.GetChangeType()
		var rankInfo string
		switch changeType {
		case model.RankingChangeTypeNew:
			emoji = ":new:"
			rankInfo = fmt.Sprintf("圏外 -> %d位", *c.NewRank)
		case model.RankingChangeTypeLost:
			emoji = ":warning:"
			rankInfo = fmt.Sprintf("%d位 -> 圏外", *c.OldRank)
		default:
			if c.OldRank != nil && c.NewRank != nil {
				rankInfo = fmt.Sprintf("%d位 -> %d位 (%+d)", *c.OldRank, *c.NewRank, c.Change)
			}
		}

		line := fmt.Sprintf("%s `%s` (%s): %s", emoji, c.Keyword, c.Country, rankInfo)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (s *SlackNotifier) send(ctx context.Context, msg *SlackMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed with status: %d", resp.StatusCode)
	}

	return nil
}
