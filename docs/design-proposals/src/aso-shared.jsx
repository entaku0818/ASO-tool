// AsoApp — shared component for the ASO macOS-native redesign.
// Themed by props.theme; same layout/logic for Aurora (dark) + Daylight (light).
// Fully interactive: app/tab/keyword selection, search, sort, locale switch,
// char counters, add-keyword sheet, chart period switch.

const FONT_STACK = '-apple-system, BlinkMacSystemFont, "SF Pro Display", "SF Pro Text", "Inter", "Hiragino Sans", "Noto Sans JP", sans-serif';

// ─────────────────────────────────────────────────────────────
// Traffic lights (macOS) — colored, with hover-symbol affordance
// ─────────────────────────────────────────────────────────────
function TrafficLights({ dark }) {
  const dot = (bg) => (
    <div style={{
      width: 12, height: 12, borderRadius: '50%', background: bg,
      boxShadow: dark
        ? 'inset 0 0 0 0.5px rgba(0,0,0,0.4)'
        : 'inset 0 0 0 0.5px rgba(0,0,0,0.18)',
    }} />
  );
  return (
    <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
      {dot('#ff5f57')}{dot('#febc2e')}{dot('#28c840')}
    </div>
  );
}

// Platform badge (iOS / Android pill)
function PlatformBadge({ platform, t }) {
  const isIOS = platform === 'ios';
  return (
    <span style={{
      fontSize: 9.5, fontWeight: 600, letterSpacing: 0.3,
      padding: '2px 7px', borderRadius: 4,
      background: isIOS ? t.platformIosBg : t.platformAndBg,
      color: isIOS ? t.platformIosText : t.platformAndText,
      textTransform: 'uppercase',
    }}>{isIOS ? 'iOS' : 'Android'}</span>
  );
}

// Popularity badge — 高/中/低
function PopularityBadge({ score, t }) {
  if (score == null) return null;
  let label, color, bg;
  if (score >= 4) { label = '高'; color = t.posText; bg = t.posBg; }
  else if (score >= 2) { label = '中'; color = t.warnText; bg = t.warnBg; }
  else { label = '低'; color = t.negText; bg = t.negBg; }
  return (
    <span style={{
      fontSize: 9.5, fontWeight: 700, padding: '1px 5px', borderRadius: 3,
      background: bg, color, lineHeight: 1.4,
    }}>{label}</span>
  );
}

// Delta indicator: ▲+2 / ▼-3 / 0
function DeltaIndicator({ delta, t }) {
  if (delta == null || delta === 0) {
    return <span style={{ color: t.textDim, fontSize: 11, fontVariantNumeric: 'tabular-nums' }}>—</span>;
  }
  const up = delta > 0;
  return (
    <span style={{
      color: up ? t.posText : t.negText,
      fontSize: 11, fontWeight: 600, fontVariantNumeric: 'tabular-nums',
      display: 'inline-flex', alignItems: 'center', gap: 2,
    }}>
      <span style={{ fontSize: 8 }}>{up ? '▲' : '▼'}</span>
      {up ? `+${delta}` : delta}
    </span>
  );
}

// ─────────────────────────────────────────────────────────────
// Ranking chart — inline SVG line/area
// ─────────────────────────────────────────────────────────────
function RankingChart({ history, t, period = 30 }) {
  const data = history.slice(-period);
  if (data.length === 0) return null;
  const W = 720, H = 280, PAD_L = 44, PAD_R = 16, PAD_T = 12, PAD_B = 28;

  const ranks = data.map(d => d.rank);
  const minRank = Math.max(1, Math.min(...ranks) - 3);
  const maxRank = Math.min(150, Math.max(...ranks) + 3);

  const x = (i) => PAD_L + (i / (data.length - 1 || 1)) * (W - PAD_L - PAD_R);
  // Y inverted: rank 1 at top, big rank at bottom
  const y = (r) => PAD_T + ((r - minRank) / (maxRank - minRank || 1)) * (H - PAD_T - PAD_B);

  const linePts = data.map((d, i) => `${x(i)},${y(d.rank)}`).join(' ');
  const areaPts = `${x(0)},${H - PAD_B} ${linePts} ${x(data.length - 1)},${H - PAD_B}`;

  // Y axis grid (4 lines)
  const yTicks = 4;
  const yLabels = [];
  for (let i = 0; i <= yTicks; i++) {
    const v = Math.round(minRank + (maxRank - minRank) * (i / yTicks));
    yLabels.push({ y: y(v), label: `${v}位` });
  }

  // X axis: 5 date ticks
  const xTickIdx = [0, Math.floor((data.length - 1) / 4), Math.floor((data.length - 1) / 2), Math.floor(3 * (data.length - 1) / 4), data.length - 1];
  const xLabels = xTickIdx.map(i => {
    const d = data[i].date;
    return { x: x(i), label: `${d.getMonth() + 1}/${d.getDate()}` };
  });

  const latest = data[data.length - 1];
  const prev = data[data.length - 2];
  const trend = prev ? latest.rank - prev.rank : 0;

  return (
    <svg viewBox={`0 0 ${W} ${H}`} width="100%" style={{ display: 'block' }}>
      <defs>
        <linearGradient id={`grad-${t.name}`} x1="0" x2="0" y1="0" y2="1">
          <stop offset="0%" stopColor={t.accent} stopOpacity="0.35" />
          <stop offset="100%" stopColor={t.accent} stopOpacity="0" />
        </linearGradient>
      </defs>
      {/* Grid */}
      {yLabels.map((g, i) => (
        <line key={i} x1={PAD_L} x2={W - PAD_R} y1={g.y} y2={g.y}
          stroke={t.gridLine} strokeWidth="0.5" strokeDasharray={i === 0 ? '' : '2 4'} />
      ))}
      {/* Y labels */}
      {yLabels.map((g, i) => (
        <text key={i} x={PAD_L - 8} y={g.y + 3} fill={t.textDim} fontSize="10"
          textAnchor="end" fontFamily={FONT_STACK} style={{ fontVariantNumeric: 'tabular-nums' }}>
          {g.label}
        </text>
      ))}
      {/* X labels */}
      {xLabels.map((g, i) => (
        <text key={i} x={g.x} y={H - 8} fill={t.textDim} fontSize="10"
          textAnchor="middle" fontFamily={FONT_STACK} style={{ fontVariantNumeric: 'tabular-nums' }}>
          {g.label}
        </text>
      ))}
      {/* Area */}
      <polygon points={areaPts} fill={`url(#grad-${t.name})`} />
      {/* Line */}
      <polyline points={linePts} fill="none"
        stroke={t.accent} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
      {/* Points */}
      {data.map((d, i) => (
        <circle key={i} cx={x(i)} cy={y(d.rank)} r={i === data.length - 1 ? 4 : 2.2}
          fill={i === data.length - 1 ? t.accent : t.surface}
          stroke={t.accent} strokeWidth={i === data.length - 1 ? 0 : 1.4} />
      ))}
      {/* Latest pill */}
      <g transform={`translate(${x(data.length - 1) - 78}, ${y(latest.rank) - 30})`}>
        <rect width="74" height="22" rx="11" fill={t.accent} />
        <text x="37" y="15" fill={t.accentText} fontSize="11" fontWeight="700"
          textAnchor="middle" fontFamily={FONT_STACK} style={{ fontVariantNumeric: 'tabular-nums' }}>
          {`最新 ${latest.rank}位`}
        </text>
      </g>
    </svg>
  );
}

// ─────────────────────────────────────────────────────────────
// Keywords pane
// ─────────────────────────────────────────────────────────────
function KeywordsPane({ app, t, onOpenAdd }) {
  const allKeywords = window.ASO_KEYWORDS_BY_APP[app.id] || [];
  const [selectedKwId, setSelectedKwId] = React.useState(allKeywords[0]?.id);
  const [search, setSearch] = React.useState('');
  const [period, setPeriod] = React.useState(30);

  React.useEffect(() => {
    // Reset selection when app changes
    setSelectedKwId(allKeywords[0]?.id);
    setSearch('');
  }, [app.id]);

  const filtered = allKeywords.filter(k =>
    !search || k.keyword.toLowerCase().includes(search.toLowerCase())
  );
  const selected = allKeywords.find(k => k.id === selectedKwId);
  const history = selected ? window.ASO_RANKING_HISTORY[selected.id] || [] : [];

  return (
    <div style={{ display: 'flex', height: '100%' }}>
      {/* Left: keyword list */}
      <div style={{
        width: 280, display: 'flex', flexDirection: 'column',
        borderRight: `1px solid ${t.divider}`,
      }}>
        <div style={{ padding: '14px 16px 10px', display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <div style={{ flex: 1, fontSize: 13, fontWeight: 600, color: t.text, letterSpacing: -0.1 }}>
              キーワード <span style={{ color: t.textDim, fontWeight: 500, marginLeft: 4 }}>{allKeywords.length}</span>
            </div>
            <button onClick={onOpenAdd} title="追加" style={{
              ...iconBtnStyle(t),
            }}>
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round"><path d="M7 2v10M2 7h10"/></svg>
            </button>
          </div>
          <SearchInput value={search} onChange={setSearch} placeholder="検索..." t={t} />
        </div>
        <div style={{ flex: 1, overflowY: 'auto', padding: '4px 8px 12px' }}>
          {filtered.length === 0 ? (
            <div style={{ padding: 24, textAlign: 'center', color: t.textDim, fontSize: 12 }}>
              該当なし
            </div>
          ) : filtered.map(kw => (
            <KeywordRow key={kw.id} kw={kw} t={t}
              selected={kw.id === selectedKwId}
              onClick={() => setSelectedKwId(kw.id)} />
          ))}
        </div>
      </div>

      {/* Right: ranking chart */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        {!selected ? (
          <EmptyState t={t} icon="chart" title="キーワードを選択"
            desc="左のリストから選ぶと、ランキングの推移を表示します。" />
        ) : (
          <>
            <div style={{
              padding: '20px 28px 12px',
              display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 16,
            }}>
              <div style={{ minWidth: 0 }}>
                <div style={{ fontSize: 11, fontWeight: 600, color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase' }}>
                  ランキング推移
                </div>
                <div style={{ fontSize: 22, fontWeight: 600, color: t.text, marginTop: 4, letterSpacing: -0.4 }}>
                  {selected.keyword}
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginTop: 8, color: t.textMuted, fontSize: 12 }}>
                  <span>国: <span style={{ color: t.text, fontWeight: 500 }}>{selected.country}</span></span>
                  <span style={{ color: t.divider }}>·</span>
                  <span style={{ display: 'inline-flex', alignItems: 'center', gap: 5 }}>
                    人気度 <PopularityBadge score={selected.popularityScore} t={t} />
                  </span>
                  <span style={{ color: t.divider }}>·</span>
                  <span>30日変化 <DeltaIndicator delta={selected.delta} t={t} /></span>
                </div>
              </div>
              {/* Period segmented */}
              <Segmented value={period} onChange={setPeriod} t={t}
                options={[{ v: 7, l: '7日' }, { v: 14, l: '14日' }, { v: 30, l: '30日' }]} />
            </div>

            {/* Chart */}
            <div style={{ padding: '4px 28px 16px' }}>
              <div style={{
                padding: 18, borderRadius: 12, background: t.surfaceSubtle,
                border: `1px solid ${t.border}`,
              }}>
                <RankingChart history={history} t={t} period={period} />
              </div>
            </div>

            {/* Stats row */}
            <div style={{
              padding: '0 28px 24px', display: 'grid',
              gridTemplateColumns: 'repeat(4, 1fr)', gap: 12,
            }}>
              {(() => {
                const data = history.slice(-period);
                const ranks = data.map(d => d.rank);
                const cur = selected.currentRank;
                const best = ranks.length ? Math.min(...ranks) : null;
                const worst = ranks.length ? Math.max(...ranks) : null;
                const avg = ranks.length ? Math.round(ranks.reduce((a, b) => a + b, 0) / ranks.length) : null;
                const stats = [
                  { label: '現在', value: cur != null ? `${cur}位` : '圏外' },
                  { label: 'ベスト', value: best != null ? `${best}位` : '—' },
                  { label: '平均', value: avg != null ? `${avg}位` : '—' },
                  { label: 'ワースト', value: worst != null ? `${worst}位` : '—' },
                ];
                return stats.map((s, i) => (
                  <StatCard key={i} t={t} label={s.label} value={s.value} />
                ));
              })()}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

function KeywordRow({ kw, selected, onClick, t }) {
  return (
    <div onClick={onClick} style={{
      padding: '9px 10px', borderRadius: 8, cursor: 'pointer',
      background: selected ? t.rowSelectedBg : 'transparent',
      transition: 'background 120ms',
      marginBottom: 2,
    }}
    onMouseEnter={(e) => { if (!selected) e.currentTarget.style.background = t.rowHoverBg; }}
    onMouseLeave={(e) => { if (!selected) e.currentTarget.style.background = 'transparent'; }}
    >
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 8 }}>
        <span style={{
          fontSize: 13, color: selected ? t.rowSelectedText : t.text,
          fontWeight: selected ? 600 : 500, letterSpacing: -0.1,
          overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1, minWidth: 0,
        }}>{kw.keyword}</span>
        <span style={{
          fontSize: 12, fontWeight: 600,
          color: kw.currentRank == null ? t.textDim : t.text,
          fontVariantNumeric: 'tabular-nums', flexShrink: 0,
        }}>
          {kw.currentRank == null ? '—' : `${kw.currentRank}位`}
        </span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginTop: 4 }}>
        <span style={{ fontSize: 10, color: t.textDim, fontWeight: 500 }}>{kw.country}</span>
        <PopularityBadge score={kw.popularityScore} t={t} />
        <div style={{ flex: 1 }} />
        <DeltaIndicator delta={kw.delta} t={t} />
      </div>
    </div>
  );
}

function StatCard({ label, value, t }) {
  return (
    <div style={{
      padding: '12px 14px', borderRadius: 10,
      background: t.surfaceSubtle, border: `1px solid ${t.border}`,
    }}>
      <div style={{ fontSize: 10, fontWeight: 600, color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase' }}>
        {label}
      </div>
      <div style={{
        fontSize: 20, fontWeight: 600, color: t.text, marginTop: 4,
        letterSpacing: -0.3, fontVariantNumeric: 'tabular-nums',
      }}>
        {value}
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Search input
// ─────────────────────────────────────────────────────────────
function SearchInput({ value, onChange, placeholder, t, width }) {
  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 6,
      padding: '6px 10px', borderRadius: 8,
      background: t.inputBg, border: `1px solid ${t.border}`,
      width: width || '100%',
    }}>
      <svg width="12" height="12" viewBox="0 0 13 13" fill="none" stroke={t.textDim} strokeWidth="1.5">
        <circle cx="5.5" cy="5.5" r="4" /><path d="M8.5 8.5l3 3" strokeLinecap="round"/>
      </svg>
      <input type="text" value={value} onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        style={{
          flex: 1, minWidth: 0, border: 'none', outline: 'none', background: 'transparent',
          fontSize: 12, color: t.text, fontFamily: FONT_STACK,
        }} />
    </div>
  );
}

// Segmented control
function Segmented({ value, onChange, options, t }) {
  return (
    <div style={{
      display: 'inline-flex', padding: 2, borderRadius: 8,
      background: t.segBg, border: `1px solid ${t.border}`,
    }}>
      {options.map(opt => {
        const sel = opt.v === value;
        return (
          <button key={opt.v} onClick={() => onChange(opt.v)}
            style={{
              border: 'none', cursor: 'pointer',
              padding: '5px 12px', borderRadius: 6,
              fontSize: 11.5, fontWeight: 600, fontFamily: FONT_STACK,
              color: sel ? t.segActiveText : t.textMuted,
              background: sel ? t.segActiveBg : 'transparent',
              boxShadow: sel ? t.segActiveShadow : 'none',
              transition: 'all 120ms',
            }}>
            {opt.l}
          </button>
        );
      })}
    </div>
  );
}

// Empty state
function EmptyState({ icon, title, desc, t }) {
  const icons = {
    chart: (<path d="M3 17l6-6 4 4 8-8" strokeLinecap="round" strokeLinejoin="round"/>),
    select: (<><rect x="3" y="3" width="18" height="18" rx="3" strokeDasharray="3 3"/><path d="M9 12l2 2 4-4" strokeLinecap="round" strokeLinejoin="round"/></>),
    check: (<><circle cx="12" cy="12" r="9"/><path d="M8 12l3 3 5-6" strokeLinecap="round" strokeLinejoin="round"/></>),
  };
  return (
    <div style={{
      flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center',
      justifyContent: 'center', padding: 40, gap: 12,
    }}>
      <div style={{
        width: 56, height: 56, borderRadius: '50%',
        background: t.surfaceSubtle, border: `1px solid ${t.border}`,
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        color: t.textDim,
      }}>
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6">
          {icons[icon] || icons.select}
        </svg>
      </div>
      <div style={{ fontSize: 15, fontWeight: 600, color: t.text }}>{title}</div>
      <div style={{ fontSize: 12, color: t.textMuted, textAlign: 'center', maxWidth: 280, lineHeight: 1.5 }}>
        {desc}
      </div>
    </div>
  );
}

function iconBtnStyle(t) {
  return {
    width: 24, height: 24, borderRadius: 6, border: 'none', cursor: 'pointer',
    background: 'transparent', color: t.textMuted,
    display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
    transition: 'all 120ms',
  };
}

Object.assign(window, {
  TrafficLights, PlatformBadge, PopularityBadge, DeltaIndicator,
  RankingChart, KeywordsPane, KeywordRow, SearchInput, Segmented,
  EmptyState, StatCard, iconBtnStyle, FONT_STACK,
});
