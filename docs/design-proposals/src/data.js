// Shared mock data for ASO tool variants

const APPS = [
  {
    id: 'app-1',
    name: 'ハビットガーデン',
    bundleID: 'app.habitgarden.ios',
    platform: 'ios',
    icon: '🌱',
    iconBg: 'linear-gradient(135deg, #6FCF97, #27AE60)',
    category: '健康＆フィットネス',
    rating: 4.7,
    reviews: 1284,
  },
  {
    id: 'app-2',
    name: 'マネーフロー',
    bundleID: 'app.moneyflow.android',
    platform: 'android',
    icon: '¥',
    iconBg: 'linear-gradient(135deg, #5B8DEF, #2D5BFF)',
    category: 'ファイナンス',
    rating: 4.5,
    reviews: 3812,
  },
  {
    id: 'app-3',
    name: 'スリープスカイ',
    bundleID: 'app.sleepsky.ios',
    platform: 'ios',
    icon: '☾',
    iconBg: 'linear-gradient(135deg, #8E7CFF, #5B47E6)',
    category: 'ヘルスケア',
    rating: 4.8,
    reviews: 921,
  },
  {
    id: 'app-4',
    name: 'レシピボックス',
    bundleID: 'app.recipebox.ios',
    platform: 'ios',
    icon: '🍳',
    iconBg: 'linear-gradient(135deg, #FFB36B, #F17E63)',
    category: 'フード＆ドリンク',
    rating: 4.6,
    reviews: 2456,
  },
];

// Keywords per app
const KEYWORDS_BY_APP = {
  'app-1': [
    { id: 'k1', keyword: '習慣 アプリ',     country: 'JP', popularityScore: 5, currentRank: 4,    delta: +2 },
    { id: 'k2', keyword: '習慣化',           country: 'JP', popularityScore: 5, currentRank: 12,   delta: -1 },
    { id: 'k3', keyword: 'todo',             country: 'JP', popularityScore: 4, currentRank: 28,   delta: +5 },
    { id: 'k4', keyword: '目標管理',         country: 'JP', popularityScore: 3, currentRank: 7,    delta: 0  },
    { id: 'k5', keyword: 'タスク 無料',      country: 'JP', popularityScore: 4, currentRank: 45,   delta: -3 },
    { id: 'k6', keyword: 'habit tracker',    country: 'US', popularityScore: 5, currentRank: 33,   delta: +8 },
    { id: 'k7', keyword: 'ルーティン',       country: 'JP', popularityScore: 3, currentRank: 9,    delta: +1 },
    { id: 'k8', keyword: '生活改善',         country: 'JP', popularityScore: 2, currentRank: null, delta: 0  },
    { id: 'k9', keyword: 'モチベーション',   country: 'JP', popularityScore: 3, currentRank: 22,   delta: +4 },
    { id: 'k10', keyword: '朝活',            country: 'JP', popularityScore: 4, currentRank: 6,    delta: +3 },
  ],
  'app-2': [
    { id: 'k11', keyword: '家計簿',          country: 'JP', popularityScore: 5, currentRank: 8,    delta: +1 },
    { id: 'k12', keyword: '家計簿 簡単',     country: 'JP', popularityScore: 4, currentRank: 14,   delta: 0  },
    { id: 'k13', keyword: '節約',            country: 'JP', popularityScore: 5, currentRank: 23,   delta: -2 },
    { id: 'k14', keyword: '予算管理',        country: 'JP', popularityScore: 3, currentRank: 4,    delta: +6 },
  ],
  'app-3': [
    { id: 'k15', keyword: '睡眠 記録',       country: 'JP', popularityScore: 4, currentRank: 5,    delta: +2 },
    { id: 'k16', keyword: '快眠',            country: 'JP', popularityScore: 3, currentRank: 18,   delta: +1 },
    { id: 'k17', keyword: 'いびき',          country: 'JP', popularityScore: 3, currentRank: null, delta: 0  },
  ],
  'app-4': [
    { id: 'k18', keyword: 'レシピ',          country: 'JP', popularityScore: 5, currentRank: 31,   delta: +4 },
    { id: 'k19', keyword: '献立',            country: 'JP', popularityScore: 4, currentRank: 17,   delta: -1 },
    { id: 'k20', keyword: '料理 簡単',       country: 'JP', popularityScore: 4, currentRank: 9,    delta: +3 },
  ],
};

// Generate ranking history (30 days) per keyword
function makeRankingHistory(seedRank, volatility = 5) {
  const days = 30;
  const out = [];
  let r = seedRank + Math.floor(Math.random() * 10);
  for (let i = days - 1; i >= 0; i--) {
    const drift = (Math.random() - 0.5) * volatility * 2;
    r = Math.max(1, Math.min(120, Math.round(r + drift)));
    const d = new Date();
    d.setDate(d.getDate() - i);
    out.push({ date: d, rank: r });
  }
  // Force last point to current rank for consistency
  if (out.length > 0) out[out.length - 1].rank = seedRank;
  return out;
}

// Pre-generate per keyword (deterministic enough; user won't compare)
const RANKING_HISTORY = {};
Object.values(KEYWORDS_BY_APP).flat().forEach((kw) => {
  RANKING_HISTORY[kw.id] = makeRankingHistory(kw.currentRank ?? 80, 4);
});

// Competitor gap per app
const GAP_BY_APP = {
  'app-1': [
    { id: 'g1', keyword: '日記 アプリ',     country: 'JP', competitorName: 'DailyLog',     competitorRank: 3,  ourRank: 45 },
    { id: 'g2', keyword: 'マインドフル',    country: 'JP', competitorName: 'Calm Garden',  competitorRank: 7,  ourRank: null },
    { id: 'g3', keyword: '禁煙',            country: 'JP', competitorName: 'Quit Smoke',   competitorRank: 5,  ourRank: 88 },
    { id: 'g4', keyword: 'ダイエット 記録', country: 'JP', competitorName: 'BodyTrack',    competitorRank: 11, ourRank: null },
    { id: 'g5', keyword: '英語学習',        country: 'JP', competitorName: 'LearnFlow',    competitorRank: 9,  ourRank: 67 },
    { id: 'g6', keyword: '読書 記録',       country: 'JP', competitorName: 'BookShelf',    competitorRank: 14, ourRank: null },
    { id: 'g7', keyword: '瞑想',            country: 'JP', competitorName: 'Calm Garden',  competitorRank: 6,  ourRank: 52 },
    { id: 'g8', keyword: 'ヨガ',            country: 'JP', competitorName: 'YogaFlow',     competitorRank: 17, ourRank: null },
  ],
  'app-2': [
    { id: 'g9',  keyword: '貯金',           country: 'JP', competitorName: 'SaveUp',       competitorRank: 4,  ourRank: 38 },
    { id: 'g10', keyword: '投資 初心者',    country: 'JP', competitorName: 'InvestEZ',     competitorRank: 8,  ourRank: null },
    { id: 'g11', keyword: 'ポイント管理',   country: 'JP', competitorName: 'PointKeeper',  competitorRank: 12, ourRank: 71 },
  ],
  'app-3': [],
  'app-4': [
    { id: 'g12', keyword: '弁当',           country: 'JP', competitorName: 'BentoMaker',   competitorRank: 6,  ourRank: 49 },
    { id: 'g13', keyword: 'ダイエット 料理', country: 'JP', competitorName: 'DietChef',    competitorRank: 10, ourRank: null },
  ],
};

// Metadata versions per app
const METADATA_BY_APP = {
  'app-1': [
    {
      id: 'm1', locale: 'ja', versionTag: 'v2.4',
      title: 'ハビットガーデン - 習慣を育てる',
      subtitle: '小さな積み重ねで、人生を変える',
      keywords: '習慣,習慣化,todo,タスク,目標,ルーティン,朝活,記録,生活,改善',
      promotionalText: '新機能：AIコーチがあなたの習慣を分析し、最適な時間帯を提案します。今すぐ試してみよう。',
      description: '「ハビットガーデン」は、毎日の小さな習慣を可視化して、人生をゆっくり変えていくアプリです。\n\n■主な機能\n・シンプルな習慣チェック\n・連続記録の可視化\n・週次・月次レポート\n・AIコーチによる改善提案\n・Apple Watch対応\n\n小さく始めて、大きく育てる。今日から、あなたの習慣を育ててみませんか？',
      updatedAt: '2026-05-12',
    },
    {
      id: 'm2', locale: 'ja', versionTag: 'draft',
      title: 'ハビットガーデン｜習慣化アプリ',
      subtitle: '毎朝5分で、新しい自分に',
      keywords: '習慣,習慣化,朝活,タスク管理,目標達成,セルフケア,日記,健康,記録',
      promotionalText: '夏のリニューアル記念で、プレミアム機能が30日間無料！',
      description: '習慣を育てる、毎日の小さな一歩。\n\nハビットガーデンは...',
      updatedAt: '2026-05-15',
    },
    {
      id: 'm3', locale: 'en-US', versionTag: 'v2.4',
      title: 'Habit Garden - Grow Daily',
      subtitle: 'Small steps, big change',
      keywords: 'habit,tracker,routine,goal,daily,morning,journal,wellness',
      promotionalText: 'New: AI Coach analyzes your habits and suggests the best time to act.',
      description: 'Habit Garden helps you build small daily habits that change your life over time...',
      updatedAt: '2026-05-12',
    },
  ],
  'app-2': [
    {
      id: 'm4', locale: 'ja', versionTag: 'v3.1',
      title: 'マネーフロー｜家計簿アプリ',
      subtitle: '見える化で、お金を味方に',
      keywords: '家計簿,節約,予算,貯金,投資,お金,管理,自動,シンプル',
      promotionalText: '銀行・カード連携で、入力ゼロの家計簿。',
      description: 'マネーフローは、あなたのお金の流れを自動で見える化する家計簿アプリです...',
      updatedAt: '2026-05-08',
    },
  ],
  'app-3': [],
  'app-4': [],
};

const LOCALES = [
  { code: 'ja',      name: '日本語',         flag: '🇯🇵' },
  { code: 'en-US',   name: 'English (US)',    flag: '🇺🇸' },
  { code: 'zh-Hans', name: '中文 (简体)',      flag: '🇨🇳' },
  { code: 'ko',      name: '한국어',           flag: '🇰🇷' },
];

const CHAR_LIMITS = {
  title: 30,
  subtitle: 30,
  keywords: 100,
  promotionalText: 170,
  description: 4000,
};

Object.assign(window, {
  ASO_APPS: APPS,
  ASO_KEYWORDS_BY_APP: KEYWORDS_BY_APP,
  ASO_RANKING_HISTORY: RANKING_HISTORY,
  ASO_GAP_BY_APP: GAP_BY_APP,
  ASO_METADATA_BY_APP: METADATA_BY_APP,
  ASO_LOCALES: LOCALES,
  ASO_CHAR_LIMITS: CHAR_LIMITS,
});
