package repository

import (
	"testing"
	"time"
)

func ptr(i int) *int { return &i }

// at builds a UTC timestamp from a JST wall-clock time, so the tests read in the
// timezone the daily batch actually runs in.
func at(y int, m time.Month, d, hh, mm int) *time.Time {
	t := time.Date(y, m, d, hh, mm, 0, 0, jst).UTC()
	return &t
}

func row(id string, rank *int, recordedAt *time.Time) keywordRankRow {
	return keywordRankRow{
		KeywordID:  id,
		Keyword:    "録音",
		Country:    "JP",
		Rank:       rank,
		RecordedAt: recordedAt,
	}
}

func TestSummarizeKeywordRanks(t *testing.T) {
	tests := []struct {
		name         string
		rows         []keywordRankRow
		wantCurrent  *int
		wantPrevious *int
		wantChange   *int
	}{
		{
			// 本 issue の再現ケース: 同じ日に日次バッチ分と、macOSアプリを開いたことによる
			// 臨時スクレイプ分の2行がある。前日(08-01)と比較しなければならない。
			name: "同日に複数スクレイプがあっても前日と比較する",
			rows: []keywordRankRow{
				row("k1", ptr(15), at(2026, 8, 1, 11, 23)),
				row("k1", ptr(12), at(2026, 8, 2, 11, 14)),
				row("k1", ptr(12), at(2026, 8, 2, 21, 47)), // 臨時スクレイプ
			},
			wantCurrent:  ptr(12),
			wantPrevious: ptr(15),
			wantChange:   ptr(3), // 15 -> 12 なので3つ改善
		},
		{
			name: "同日の複数行のうち最新のスクレイプを当日の代表値にする",
			rows: []keywordRankRow{
				row("k1", ptr(20), at(2026, 8, 1, 11, 0)),
				row("k1", ptr(18), at(2026, 8, 2, 11, 0)),
				row("k1", ptr(9), at(2026, 8, 2, 23, 30)), // これが当日の代表値
			},
			wantCurrent:  ptr(9),
			wantPrevious: ptr(20),
			wantChange:   ptr(11),
		},
		{
			name: "前日データが無い(登録直後で当日分しかない)場合は previous/change とも nil",
			rows: []keywordRankRow{
				row("k1", ptr(12), at(2026, 8, 2, 11, 14)),
				row("k1", ptr(12), at(2026, 8, 2, 21, 47)),
			},
			wantCurrent:  ptr(12),
			wantPrevious: nil,
			wantChange:   nil,
		},
		{
			name:         "順位履歴が1件も無いキーワードも結果に含まれる",
			rows:         []keywordRankRow{row("k1", nil, nil)},
			wantCurrent:  nil,
			wantPrevious: nil,
			wantChange:   nil,
		},
		{
			name: "前日が欠けていても直近の別日と比較する",
			rows: []keywordRankRow{
				row("k1", ptr(30), at(2026, 7, 28, 11, 0)), // 5日前
				row("k1", ptr(25), at(2026, 8, 2, 11, 0)),
			},
			wantCurrent:  ptr(25),
			wantPrevious: ptr(30),
			wantChange:   ptr(5),
		},
		{
			name: "当日が圏外(rank=nil)なら change は出さない",
			rows: []keywordRankRow{
				row("k1", ptr(40), at(2026, 8, 1, 11, 0)),
				row("k1", nil, at(2026, 8, 2, 11, 0)),
			},
			wantCurrent:  nil,
			wantPrevious: ptr(40),
			wantChange:   nil,
		},
		{
			name: "前日が圏外(rank=nil)なら change は出さない",
			rows: []keywordRankRow{
				row("k1", nil, at(2026, 8, 1, 11, 0)),
				row("k1", ptr(40), at(2026, 8, 2, 11, 0)),
			},
			wantCurrent:  ptr(40),
			wantPrevious: nil,
			wantChange:   nil,
		},
		{
			// UTC では 08-01T23:00Z と 08-02T02:00Z で日付が変わってしまうが、
			// JST ではどちらも 08-02。同日として畳まれなければならない。
			name: "日付の境界は JST で判定する",
			rows: []keywordRankRow{
				row("k1", ptr(50), at(2026, 8, 1, 10, 0)),
				row("k1", ptr(20), at(2026, 8, 2, 8, 0)),  // UTC では 08-01T23:00Z
				row("k1", ptr(10), at(2026, 8, 2, 11, 0)), // UTC では 08-02T02:00Z
			},
			wantCurrent:  ptr(10),
			wantPrevious: ptr(50),
			wantChange:   ptr(40),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summarizeKeywordRanks(tt.rows)
			if len(got) != 1 {
				t.Fatalf("結果は1件のはず: got %d", len(got))
			}
			assertIntPtr(t, "current_rank", got[0].CurrentRank, tt.wantCurrent)
			assertIntPtr(t, "previous_rank", got[0].PreviousRank, tt.wantPrevious)
			assertIntPtr(t, "change", got[0].Change, tt.wantChange)
		})
	}
}

func TestSummarizeKeywordRanksMultipleKeywords(t *testing.T) {
	rows := []keywordRankRow{
		{KeywordID: "k1", Keyword: "録音", Country: "JP", Rank: ptr(15), RecordedAt: at(2026, 8, 1, 11, 0)},
		{KeywordID: "k1", Keyword: "録音", Country: "JP", Rank: ptr(12), RecordedAt: at(2026, 8, 2, 11, 0)},
		{KeywordID: "k2", Keyword: "議事録", Country: "JP", Rank: ptr(5), RecordedAt: at(2026, 8, 1, 11, 0)},
		{KeywordID: "k2", Keyword: "議事録", Country: "JP", Rank: ptr(8), RecordedAt: at(2026, 8, 2, 11, 0)},
		{KeywordID: "k3", Keyword: "取材", Country: "JP", Rank: nil, RecordedAt: nil},
	}

	got := summarizeKeywordRanks(rows)
	if len(got) != 3 {
		t.Fatalf("キーワード3件のはず: got %d", len(got))
	}

	// 入力の出現順（= SQL の ORDER BY）が保たれること
	wantOrder := []string{"k1", "k2", "k3"}
	for i, id := range wantOrder {
		if got[i].KeywordID != id {
			t.Errorf("順序が違う: got[%d].KeywordID = %q, want %q", i, got[i].KeywordID, id)
		}
	}

	assertIntPtr(t, "k1 change", got[0].Change, ptr(3))  // 15 -> 12 改善
	assertIntPtr(t, "k2 change", got[1].Change, ptr(-3)) // 5 -> 8 悪化
	assertIntPtr(t, "k3 current", got[2].CurrentRank, nil)

	if got[1].Keyword != "議事録" || got[1].Country != "JP" {
		t.Errorf("keyword/country が引き継がれていない: %+v", got[1])
	}
}

func assertIntPtr(t *testing.T, label string, got, want *int) {
	t.Helper()
	switch {
	case want == nil && got == nil:
	case want == nil:
		t.Errorf("%s = %d, want nil", label, *got)
	case got == nil:
		t.Errorf("%s = nil, want %d", label, *want)
	case *got != *want:
		t.Errorf("%s = %d, want %d", label, *got, *want)
	}
}
