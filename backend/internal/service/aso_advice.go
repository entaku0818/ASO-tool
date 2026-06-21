package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

// ASOAdvice holds a single actionable advice item.
type ASOAdvice struct {
	Priority    string `json:"priority"`    // high / medium / low
	Category    string `json:"category"`    // keyword / competitor / metadata / ranking
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// ASOAdviceResponse is the full response from the advice endpoint.
type ASOAdviceResponse struct {
	Summary string      `json:"summary"`
	Advice  []ASOAdvice `json:"advice"`
}

type ASOAdviceService struct {
	keywordRepo    *repository.KeywordRepository
	rankingRepo    *repository.RankingRepository
	competitorRepo *repository.CompetitorRepository
	metadataRepo   *repository.MetadataRepository
	appRepo        *repository.AppRepository
	httpClient     *http.Client
	apiKey         string
}

func NewASOAdviceService(
	keywordRepo *repository.KeywordRepository,
	rankingRepo *repository.RankingRepository,
	competitorRepo *repository.CompetitorRepository,
	metadataRepo *repository.MetadataRepository,
	appRepo *repository.AppRepository,
) *ASOAdviceService {
	return &ASOAdviceService{
		keywordRepo:    keywordRepo,
		rankingRepo:    rankingRepo,
		competitorRepo: competitorRepo,
		metadataRepo:   metadataRepo,
		appRepo:        appRepo,
		httpClient:     &http.Client{Timeout: 60 * time.Second},
		apiKey:         os.Getenv("ANTHROPIC_API_KEY"),
	}
}

func (s *ASOAdviceService) IsConfigured() bool { return s.apiKey != "" }

// Generate collects app data and returns AI-generated ASO advice.
func (s *ASOAdviceService) Generate(ctx context.Context, userID, appID string) (*ASOAdviceResponse, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	// Fetch app info
	app, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, fmt.Errorf("get app: %w", err)
	}

	// Fetch keywords + ranks
	keywords, _ := s.keywordRepo.ListByApp(ctx, appID)
	ranks, _ := s.rankingRepo.GetAllKeywordRanks(ctx, appID)

	rankMap := map[string]*repository.KeywordRankSummary{}
	for i := range ranks {
		rankMap[ranks[i].KeywordID] = &ranks[i]
	}

	// Fetch competitor gap
	gaps, _ := s.competitorRepo.GetKeywordGap(ctx, appID, 20, 30)

	// Fetch metadata
	metaList, _ := s.metadataRepo.List(ctx, appID)

	// Fetch rising keywords
	rising, _ := s.rankingRepo.GetRisingKeywords(ctx, appID, 7)

	prompt := buildASOAdvicePrompt(app, keywords, rankMap, gaps, metaList, rising)

	advice, err := s.callClaude(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return advice, nil
}

func buildASOAdvicePrompt(
	app *model.App,
	keywords []*model.Keyword,
	rankMap map[string]*repository.KeywordRankSummary,
	gaps []*model.KeywordGap,
	metaList []*model.AppMetadataVersion,
	rising []repository.RisingKeyword,
) string {
	var sb strings.Builder

	sb.WriteString("あなたはApp Store最適化（ASO）の専門家です。以下のデータを分析して、具体的で実行可能なASOアドバイスを日本語で提供してください。\n\n")
	sb.WriteString(fmt.Sprintf("## アプリ情報\n- 名前: %s\n- プラットフォーム: %s\n\n", app.Name, app.Platform))

	// Keywords & ranks (top 30)
	sb.WriteString("## 現在のキーワードと順位\n")
	count := 0
	for _, kw := range keywords {
		if count >= 30 {
			break
		}
		rankStr := "圏外"
		changeStr := ""
		if r, ok := rankMap[kw.ID]; ok && r.CurrentRank != nil {
			rankStr = fmt.Sprintf("%d位", *r.CurrentRank)
			if r.Change != nil && *r.Change != 0 {
				if *r.Change > 0 {
					changeStr = fmt.Sprintf(" (↑%d)", *r.Change)
				} else {
					changeStr = fmt.Sprintf(" (↓%d)", -*r.Change)
				}
			}
		}
		sb.WriteString(fmt.Sprintf("- %s [%s]: %s%s\n", kw.Keyword, kw.Country, rankStr, changeStr))
		count++
	}

	// Rising
	if len(rising) > 0 {
		sb.WriteString("\n## 急上昇キーワード（過去7日）\n")
		for _, r := range rising {
			sb.WriteString(fmt.Sprintf("- %s: %d位→%d位 (+%d)\n", r.Keyword, r.PreviousRank, r.CurrentRank, r.Improvement))
		}
	}

	// Competitor gaps (top 10)
	if len(gaps) > 0 {
		sb.WriteString("\n## 競合ギャップ（競合が上位・自社が圏外）\n")
		for i, g := range gaps {
			if i >= 10 {
				break
			}
			sb.WriteString(fmt.Sprintf("- %s [%s]: 競合%d位、自社圏外 (%s)\n", g.Keyword, g.Country, g.CompetitorRank, g.CompetitorName))
		}
	}

	// Metadata
	if len(metaList) > 0 {
		sb.WriteString("\n## メタデータ\n")
		for _, m := range metaList {
			sb.WriteString(fmt.Sprintf("### %s (%s)\n", m.Locale, m.VersionTag))
			if m.Title != nil {
				sb.WriteString(fmt.Sprintf("- タイトル: %s\n", *m.Title))
			}
			if m.Subtitle != nil {
				sb.WriteString(fmt.Sprintf("- サブタイトル: %s\n", *m.Subtitle))
			}
			if m.Keywords != nil {
				sb.WriteString(fmt.Sprintf("- キーワードフィールド: %s\n", *m.Keywords))
			}
		}
	}

	sb.WriteString(`
## 出力形式
以下のJSON形式のみで回答してください（余計なテキスト不要）:
{
  "summary": "全体的なASO状況の2-3文の要約",
  "advice": [
    {
      "priority": "high|medium|low",
      "category": "keyword|competitor|metadata|ranking",
      "title": "アドバイスのタイトル（20文字以内）",
      "description": "詳細な説明（100文字以内）",
      "action": "具体的なアクション（50文字以内）"
    }
  ]
}
優先度highを2-3件、mediumを2-3件、lowを1-2件、合計5-8件のアドバイスを返してください。`)

	return sb.String()
}

type claudeAdviceResp struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (s *ASOAdviceService) callClaude(ctx context.Context, prompt string) (*ASOAdviceResponse, error) {
	reqBody, err := json.Marshal(map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1500,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, claudeAPIURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude api: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude api error %d: %s", resp.StatusCode, body)
	}

	var cr claudeAdviceResp
	if err := json.Unmarshal(body, &cr); err != nil || len(cr.Content) == 0 {
		return nil, fmt.Errorf("parse claude response: %w", err)
	}

	text := cr.Content[0].Text
	// Extract JSON block
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("no JSON in response: %s", text)
	}
	text = text[start : end+1]

	var result ASOAdviceResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("unmarshal advice: %w", err)
	}
	return &result, nil
}
