// AsoApp panes — Gap, Metadata, AddKeyword sheet — plus the main AsoApp shell.

// ─────────────────────────────────────────────────────────────
// Gap pane — table of competitor keywords
// ─────────────────────────────────────────────────────────────
function GapPane({ app, t }) {
  const all = window.ASO_GAP_BY_APP[app.id] || [];
  const [sortKey, setSortKey] = React.useState('competitorRank');
  const [sortDir, setSortDir] = React.useState('asc');
  const [filter, setFilter] = React.useState('all'); // all | unranked

  const filtered = all.filter(g => filter === 'all' ? true : g.ourRank == null);
  const sorted = [...filtered].sort((a, b) => {
    const av = a[sortKey], bv = b[sortKey];
    if (av == null && bv == null) return 0;
    if (av == null) return 1;
    if (bv == null) return -1;
    if (typeof av === 'string') return sortDir === 'asc' ? av.localeCompare(bv) : bv.localeCompare(av);
    return sortDir === 'asc' ? av - bv : bv - av;
  });

  const toggleSort = (k) => {
    if (sortKey === k) setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
    else { setSortKey(k); setSortDir('asc'); }
  };

  if (all.length === 0) {
    return (
      <EmptyState t={t} icon="check" title="ギャップなし"
        desc="競合が上位かつ自社が圏外のキーワードはありません。素晴らしい！" />
    );
  }

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Header row */}
      <div style={{
        padding: '20px 28px 14px', display: 'flex', alignItems: 'flex-end',
        justifyContent: 'space-between', gap: 16,
      }}>
        <div>
          <div style={{ fontSize: 11, fontWeight: 600, color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase' }}>
            キーワードギャップ
          </div>
          <div style={{ fontSize: 20, fontWeight: 600, color: t.text, marginTop: 4, letterSpacing: -0.3 }}>
            競合が上位、自社が圏外
          </div>
          <div style={{ fontSize: 12, color: t.textMuted, marginTop: 6, maxWidth: 480 }}>
            競合が20位以内、自社が30位以下またはランク外のキーワード。攻めるべき領域です。
          </div>
        </div>
        <div style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
          <Segmented value={filter} onChange={setFilter} t={t}
            options={[{ v: 'all', l: 'すべて' }, { v: 'unranked', l: '圏外のみ' }]} />
          <button style={{
            ...iconBtnStyle(t), width: 'auto', padding: '6px 12px', gap: 6,
            background: t.btnGhostBg, border: `1px solid ${t.border}`, color: t.text,
            fontSize: 12, fontWeight: 500, fontFamily: FONT_STACK,
          }}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
              <path d="M6 1.5v9M2.5 7L6 10.5 9.5 7"/>
            </svg>
            CSV出力
          </button>
        </div>
      </div>

      {/* Table */}
      <div style={{ flex: 1, overflow: 'auto', padding: '0 28px 24px' }}>
        <div style={{
          border: `1px solid ${t.border}`, borderRadius: 12, overflow: 'hidden',
          background: t.surface,
        }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontFamily: FONT_STACK }}>
            <thead>
              <tr style={{ background: t.tableHeadBg }}>
                <Th t={t} label="キーワード" sortKey="keyword" current={sortKey} dir={sortDir} onClick={toggleSort} />
                <Th t={t} label="国" sortKey="country" current={sortKey} dir={sortDir} onClick={toggleSort} width={70} />
                <Th t={t} label="競合アプリ" sortKey="competitorName" current={sortKey} dir={sortDir} onClick={toggleSort} />
                <Th t={t} label="競合順位" sortKey="competitorRank" current={sortKey} dir={sortDir} onClick={toggleSort} width={110} align="right" />
                <Th t={t} label="自社順位" sortKey="ourRank" current={sortKey} dir={sortDir} onClick={toggleSort} width={110} align="right" />
                <Th t={t} label="差" sortKey="gap" current={sortKey} dir={sortDir} onClick={toggleSort} width={90} align="right" />
              </tr>
            </thead>
            <tbody>
              {sorted.map((g, i) => {
                const gap = g.ourRank == null ? '∞' : (g.ourRank - g.competitorRank);
                return (
                  <tr key={g.id} style={{
                    borderTop: `1px solid ${t.divider}`,
                  }}>
                    <Td t={t}>
                      <span style={{ fontWeight: 500, color: t.text }}>{g.keyword}</span>
                    </Td>
                    <Td t={t} muted>{g.country}</Td>
                    <Td t={t}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                        <div style={{
                          width: 18, height: 18, borderRadius: 4,
                          background: t.compIcon, display: 'flex',
                          alignItems: 'center', justifyContent: 'center',
                          fontSize: 9, fontWeight: 700, color: t.accentText,
                        }}>{g.competitorName[0]}</div>
                        <span style={{ color: t.text }}>{g.competitorName}</span>
                      </div>
                    </Td>
                    <Td t={t} align="right">
                      <span style={{
                        color: t.posText, fontWeight: 600, fontVariantNumeric: 'tabular-nums',
                      }}>{g.competitorRank}位</span>
                    </Td>
                    <Td t={t} align="right">
                      {g.ourRank == null ? (
                        <span style={{
                          fontSize: 10, fontWeight: 700, padding: '2px 7px', borderRadius: 4,
                          background: t.negBg, color: t.negText,
                        }}>圏外</span>
                      ) : (
                        <span style={{
                          color: t.negText, fontWeight: 600, fontVariantNumeric: 'tabular-nums',
                        }}>{g.ourRank}位</span>
                      )}
                    </Td>
                    <Td t={t} align="right">
                      <span style={{
                        color: t.textMuted, fontVariantNumeric: 'tabular-nums', fontSize: 11,
                      }}>{gap}</span>
                    </Td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          {sorted.length === 0 && (
            <div style={{ padding: 40, textAlign: 'center', color: t.textDim, fontSize: 12 }}>
              該当なし
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function Th({ label, sortKey, current, dir, onClick, t, width, align }) {
  const active = current === sortKey;
  return (
    <th onClick={() => onClick(sortKey)} style={{
      padding: '10px 14px',
      textAlign: align || 'left',
      fontSize: 11, fontWeight: 600, color: active ? t.text : t.textMuted,
      letterSpacing: 0.2,
      cursor: 'pointer', userSelect: 'none',
      width, borderBottom: `1px solid ${t.border}`,
    }}>
      <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
        {label}
        {active && (
          <svg width="8" height="8" viewBox="0 0 8 8" fill={t.text}>
            {dir === 'asc' ? <path d="M4 2L7 6H1z"/> : <path d="M4 6L1 2h6z"/>}
          </svg>
        )}
      </span>
    </th>
  );
}
function Td({ children, align, muted, t }) {
  return (
    <td style={{
      padding: '11px 14px', fontSize: 12.5, color: muted ? t.textMuted : t.text,
      textAlign: align || 'left', verticalAlign: 'middle',
    }}>{children}</td>
  );
}

// ─────────────────────────────────────────────────────────────
// Metadata pane
// ─────────────────────────────────────────────────────────────
function MetadataPane({ app, t }) {
  const all = window.ASO_METADATA_BY_APP[app.id] || [];
  const locales = window.ASO_LOCALES;
  const charLimits = window.ASO_CHAR_LIMITS;

  const [selectedId, setSelectedId] = React.useState(all[0]?.id || null);
  const selected = all.find(m => m.id === selectedId);

  // Editable buffer
  const [buf, setBuf] = React.useState(() => selected ? { ...selected } : emptyMeta());
  React.useEffect(() => {
    setBuf(selected ? { ...selected } : emptyMeta());
  }, [selected?.id, app.id]);

  function emptyMeta() {
    return { locale: 'ja', versionTag: 'draft', title: '', subtitle: '', keywords: '', promotionalText: '', description: '' };
  }

  const dirty = selected && (
    buf.title !== (selected.title || '') ||
    buf.subtitle !== (selected.subtitle || '') ||
    buf.keywords !== (selected.keywords || '') ||
    buf.promotionalText !== (selected.promotionalText || '') ||
    buf.description !== (selected.description || '')
  );

  return (
    <div style={{ display: 'flex', height: '100%' }}>
      {/* Left: version list */}
      <div style={{
        width: 240, display: 'flex', flexDirection: 'column',
        borderRight: `1px solid ${t.divider}`,
      }}>
        <div style={{ padding: '14px 16px 10px', display: 'flex', alignItems: 'center', gap: 8 }}>
          <div style={{ flex: 1, fontSize: 13, fontWeight: 600, color: t.text }}>
            保存済み <span style={{ color: t.textDim, fontWeight: 500, marginLeft: 4 }}>{all.length}</span>
          </div>
          <button title="新規" style={iconBtnStyle(t)}>
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round"><path d="M7 2v10M2 7h10"/></svg>
          </button>
        </div>
        <div style={{ flex: 1, overflowY: 'auto', padding: '4px 8px 12px' }}>
          {all.length === 0 ? (
            <div style={{ padding: 24, textAlign: 'center', color: t.textDim, fontSize: 12 }}>
              バージョンなし
            </div>
          ) : all.map(m => {
            const sel = m.id === selectedId;
            const loc = locales.find(l => l.code === m.locale);
            return (
              <div key={m.id} onClick={() => setSelectedId(m.id)} style={{
                padding: '10px 10px', borderRadius: 8, cursor: 'pointer',
                background: sel ? t.rowSelectedBg : 'transparent',
                marginBottom: 2, transition: 'background 120ms',
              }}
              onMouseEnter={(e) => { if (!sel) e.currentTarget.style.background = t.rowHoverBg; }}
              onMouseLeave={(e) => { if (!sel) e.currentTarget.style.background = 'transparent'; }}
              >
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 8 }}>
                  <div style={{
                    display: 'flex', alignItems: 'center', gap: 6,
                    fontSize: 12.5, color: sel ? t.rowSelectedText : t.text,
                    fontWeight: sel ? 600 : 500,
                  }}>
                    <span style={{ fontSize: 13 }}>{loc?.flag}</span>
                    {loc?.name || m.locale}
                  </div>
                  <span style={{
                    fontSize: 9.5, fontWeight: 700, padding: '2px 6px', borderRadius: 3,
                    background: m.versionTag === 'draft' ? t.warnBg : t.tagBg,
                    color: m.versionTag === 'draft' ? t.warnText : t.tagText,
                    letterSpacing: 0.3, textTransform: 'uppercase',
                  }}>{m.versionTag}</span>
                </div>
                <div style={{
                  fontSize: 11, color: t.textMuted, marginTop: 4, lineHeight: 1.4,
                  overflow: 'hidden', textOverflow: 'ellipsis', display: '-webkit-box',
                  WebkitLineClamp: 1, WebkitBoxOrient: 'vertical',
                }}>{m.title}</div>
                <div style={{ fontSize: 10, color: t.textDim, marginTop: 3 }}>
                  {m.updatedAt} 更新
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Right: editor */}
      {!selected ? (
        <EmptyState t={t} icon="select" title="バージョンを選択"
          desc="左のリストから編集したいメタデータを選んでください。" />
      ) : (
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
          {/* Toolbar */}
          <div style={{
            padding: '14px 28px', display: 'flex', alignItems: 'center', gap: 12,
            borderBottom: `1px solid ${t.divider}`,
          }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
              <span style={{ fontSize: 10, fontWeight: 600, color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase' }}>言語</span>
              <select value={buf.locale} onChange={(e) => setBuf({ ...buf, locale: e.target.value })}
                style={selectStyle(t)}>
                {locales.map(l => <option key={l.code} value={l.code}>{l.flag} {l.name}</option>)}
              </select>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
              <span style={{ fontSize: 10, fontWeight: 600, color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase' }}>バージョンタグ</span>
              <input type="text" value={buf.versionTag} onChange={(e) => setBuf({ ...buf, versionTag: e.target.value })}
                style={textInputStyle(t, 110)} />
            </div>
            <div style={{ flex: 1 }} />
            {dirty && (
              <span style={{
                fontSize: 11, color: t.warnText, fontWeight: 500,
                display: 'inline-flex', alignItems: 'center', gap: 5,
              }}>
                <span style={{ width: 6, height: 6, borderRadius: '50%', background: t.warnText }} />
                未保存の変更
              </span>
            )}
            <button style={{
              padding: '7px 14px', borderRadius: 7, border: `1px solid ${t.border}`,
              background: t.btnGhostBg, color: t.text, fontSize: 12, fontWeight: 500,
              cursor: 'pointer', fontFamily: FONT_STACK,
            }}>プレビュー</button>
            <button style={{
              padding: '7px 14px', borderRadius: 7, border: 'none',
              background: t.accent, color: t.accentText, fontSize: 12, fontWeight: 600,
              cursor: 'pointer', fontFamily: FONT_STACK,
              boxShadow: t.btnPrimaryShadow,
            }}>{dirty ? '更新' : '保存済み'}</button>
          </div>

          {/* Form */}
          <div style={{ flex: 1, overflowY: 'auto', padding: '20px 28px' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16, maxWidth: 720 }}>
              <Field t={t} label="タイトル" limit={charLimits.title} value={buf.title}
                onChange={(v) => setBuf({ ...buf, title: v })} />
              <Field t={t} label="サブタイトル" limit={charLimits.subtitle} value={buf.subtitle}
                onChange={(v) => setBuf({ ...buf, subtitle: v })} />
              <Field t={t} label="キーワード" limit={charLimits.keywords} value={buf.keywords}
                onChange={(v) => setBuf({ ...buf, keywords: v })}
                hint="カンマ区切り（半角）" />
              <TextArea t={t} label="プロモーションテキスト" limit={charLimits.promotionalText}
                value={buf.promotionalText} onChange={(v) => setBuf({ ...buf, promotionalText: v })}
                rows={3} />
              <TextArea t={t} label="説明文" limit={charLimits.description}
                value={buf.description} onChange={(v) => setBuf({ ...buf, description: v })}
                rows={9} />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function Field({ label, value, onChange, limit, hint, t }) {
  const over = value.length > limit;
  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', marginBottom: 6 }}>
        <label style={{ fontSize: 11.5, fontWeight: 600, color: t.text }}>
          {label}
          {hint && <span style={{ color: t.textDim, fontWeight: 400, marginLeft: 6 }}>{hint}</span>}
        </label>
        <span style={{
          fontSize: 10.5, fontWeight: over ? 700 : 500,
          color: over ? t.negText : t.textDim, fontVariantNumeric: 'tabular-nums',
        }}>{value.length}/{limit}</span>
      </div>
      <input type="text" value={value} onChange={(e) => onChange(e.target.value)}
        style={textInputStyle(t, '100%', over)} />
    </div>
  );
}

function TextArea({ label, value, onChange, limit, rows, t }) {
  const over = value.length > limit;
  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', marginBottom: 6 }}>
        <label style={{ fontSize: 11.5, fontWeight: 600, color: t.text }}>{label}</label>
        <span style={{
          fontSize: 10.5, fontWeight: over ? 700 : 500,
          color: over ? t.negText : t.textDim, fontVariantNumeric: 'tabular-nums',
        }}>{value.length}/{limit}</span>
      </div>
      <textarea value={value} onChange={(e) => onChange(e.target.value)} rows={rows}
        style={{
          ...textInputStyle(t, '100%', over),
          resize: 'vertical', fontFamily: FONT_STACK, lineHeight: 1.5,
          padding: '10px 12px',
        }} />
    </div>
  );
}

function textInputStyle(t, width, error) {
  return {
    width, padding: '7px 11px', borderRadius: 7,
    border: `1px solid ${error ? t.negText : t.border}`,
    background: t.inputBg, color: t.text,
    fontSize: 12.5, fontFamily: FONT_STACK, outline: 'none',
    transition: 'border-color 120ms, box-shadow 120ms',
  };
}
function selectStyle(t) {
  return {
    ...textInputStyle(t, 'auto'),
    cursor: 'pointer', paddingRight: 28,
    appearance: 'none',
    backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='10' viewBox='0 0 10 10'%3E%3Cpath d='M2 4l3 3 3-3' stroke='%23${t.textMuted.replace('#', '').slice(0, 6) || '888888'}' stroke-width='1.5' fill='none' stroke-linecap='round'/%3E%3C/svg%3E")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'right 8px center',
  };
}

// ─────────────────────────────────────────────────────────────
// Add keyword sheet (modal)
// ─────────────────────────────────────────────────────────────
function AddKeywordSheet({ app, t, onClose }) {
  const [keyword, setKeyword] = React.useState('');
  const [country, setCountry] = React.useState('JP');
  const countries = ['JP', 'US', 'GB', 'AU', 'CA', 'KR', 'DE', 'FR', 'CN'];

  return (
    <div onClick={onClose} style={{
      position: 'absolute', inset: 0, zIndex: 100,
      background: t.modalScrim,
      backdropFilter: 'blur(8px)', WebkitBackdropFilter: 'blur(8px)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      animation: 'asoFadeIn 180ms',
    }}>
      <div onClick={(e) => e.stopPropagation()} style={{
        width: 420, background: t.modalBg, borderRadius: 14,
        border: `1px solid ${t.border}`,
        boxShadow: t.modalShadow,
        padding: 28, fontFamily: FONT_STACK,
        animation: 'asoSheetIn 220ms cubic-bezier(0.22, 1, 0.36, 1)',
      }}>
        <div style={{
          display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20,
        }}>
          <div style={{
            width: 36, height: 36, borderRadius: 9, background: t.accentBg,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            color: t.accent,
          }}>
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
              <path d="M9 3v12M3 9h12"/>
            </svg>
          </div>
          <div>
            <div style={{ fontSize: 15, fontWeight: 600, color: t.text }}>キーワード追加</div>
            <div style={{ fontSize: 12, color: t.textMuted, marginTop: 2 }}>
              {app.name} に追跡対象を追加します
            </div>
          </div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
          <div>
            <label style={{ fontSize: 11.5, fontWeight: 600, color: t.text, marginBottom: 6, display: 'block' }}>キーワード</label>
            <input type="text" value={keyword} onChange={(e) => setKeyword(e.target.value)}
              autoFocus placeholder="例: 習慣 アプリ"
              style={textInputStyle(t, '100%')} />
          </div>
          <div>
            <label style={{ fontSize: 11.5, fontWeight: 600, color: t.text, marginBottom: 6, display: 'block' }}>国</label>
            <select value={country} onChange={(e) => setCountry(e.target.value)} style={{ ...selectStyle(t), width: 120 }}>
              {countries.map(c => <option key={c} value={c}>{c}</option>)}
            </select>
          </div>
        </div>

        <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 24 }}>
          <button onClick={onClose} style={{
            padding: '8px 16px', borderRadius: 7, border: `1px solid ${t.border}`,
            background: 'transparent', color: t.text, fontSize: 12.5, fontWeight: 500,
            cursor: 'pointer', fontFamily: FONT_STACK,
          }}>キャンセル</button>
          <button onClick={onClose} disabled={!keyword.trim()} style={{
            padding: '8px 18px', borderRadius: 7, border: 'none',
            background: keyword.trim() ? t.accent : t.btnDisabledBg,
            color: keyword.trim() ? t.accentText : t.textDim,
            fontSize: 12.5, fontWeight: 600,
            cursor: keyword.trim() ? 'pointer' : 'not-allowed',
            fontFamily: FONT_STACK,
            boxShadow: keyword.trim() ? t.btnPrimaryShadow : 'none',
          }}>追加</button>
        </div>
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Main AsoApp
// ─────────────────────────────────────────────────────────────
function AsoApp({ theme: t, width = 1280, height = 800 }) {
  const apps = window.ASO_APPS;
  const [selectedAppId, setSelectedAppId] = React.useState(apps[0].id);
  const [tab, setTab] = React.useState('keywords');
  const [showAdd, setShowAdd] = React.useState(false);
  const [sidebarSearch, setSidebarSearch] = React.useState('');

  const app = apps.find(a => a.id === selectedAppId);
  const filteredApps = apps.filter(a =>
    !sidebarSearch || a.name.includes(sidebarSearch) || a.bundleID.includes(sidebarSearch)
  );

  return (
    <div style={{
      width, height, position: 'relative',
      borderRadius: 12, overflow: 'hidden',
      background: t.windowBg, color: t.text,
      fontFamily: FONT_STACK,
      boxShadow: t.windowShadow,
      border: t.windowBorder,
      isolation: 'isolate',
    }}>
      {/* Ambient aurora gradient behind everything */}
      {t.aurora && <div style={{
        position: 'absolute', inset: 0, pointerEvents: 'none',
        background: t.aurora, zIndex: 0,
      }} />}

      <div style={{ display: 'flex', height: '100%', position: 'relative', zIndex: 1 }}>
        {/* Sidebar */}
        <aside style={{
          width: 248, display: 'flex', flexDirection: 'column',
          background: t.sidebarBg,
          backdropFilter: 'blur(40px) saturate(180%)',
          WebkitBackdropFilter: 'blur(40px) saturate(180%)',
          borderRight: `1px solid ${t.divider}`,
        }}>
          {/* Titlebar with traffic lights */}
          <div style={{
            height: 44, padding: '0 16px',
            display: 'flex', alignItems: 'center', gap: 12,
            WebkitAppRegion: 'drag',
          }}>
            <TrafficLights dark={t.dark} />
          </div>

          {/* App name */}
          <div style={{ padding: '6px 16px 14px' }}>
            <div style={{
              fontSize: 11, fontWeight: 700, color: t.textDim,
              letterSpacing: 1.2, textTransform: 'uppercase',
            }}>ASO Studio</div>
            <div style={{
              fontSize: 17, fontWeight: 600, color: t.text,
              letterSpacing: -0.3, marginTop: 2,
            }}>マイアプリ</div>
          </div>

          {/* Search */}
          <div style={{ padding: '0 12px 10px' }}>
            <SearchInput value={sidebarSearch} onChange={setSidebarSearch}
              placeholder="アプリを検索..." t={t} />
          </div>

          {/* App list */}
          <div style={{ flex: 1, overflowY: 'auto', padding: '4px 8px 12px' }}>
            <div style={{
              fontSize: 10, fontWeight: 700, color: t.textDim,
              padding: '6px 8px 8px', letterSpacing: 0.4, textTransform: 'uppercase',
            }}>追跡中</div>
            {filteredApps.map(a => {
              const sel = a.id === selectedAppId;
              return (
                <div key={a.id} onClick={() => setSelectedAppId(a.id)} style={{
                  padding: '8px 8px', borderRadius: 9, cursor: 'pointer',
                  display: 'flex', alignItems: 'center', gap: 10,
                  background: sel ? t.sidebarSelectedBg : 'transparent',
                  marginBottom: 2, transition: 'background 120ms',
                  position: 'relative',
                }}
                onMouseEnter={(e) => { if (!sel) e.currentTarget.style.background = t.rowHoverBg; }}
                onMouseLeave={(e) => { if (!sel) e.currentTarget.style.background = 'transparent'; }}
                >
                  {sel && t.sidebarSelectedAccent && (
                    <div style={{
                      position: 'absolute', left: -8, top: 8, bottom: 8, width: 2,
                      background: t.accent, borderRadius: 1,
                    }} />
                  )}
                  <div style={{
                    width: 32, height: 32, borderRadius: 8,
                    background: a.iconBg,
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    fontSize: 16, fontWeight: 700, color: '#fff',
                    flexShrink: 0,
                    boxShadow: '0 1px 2px rgba(0,0,0,0.1), inset 0 0 0 0.5px rgba(255,255,255,0.2)',
                  }}>{a.icon}</div>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div style={{
                      fontSize: 13, fontWeight: sel ? 600 : 500,
                      color: sel ? t.rowSelectedText : t.text,
                      overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                      letterSpacing: -0.1,
                    }}>{a.name}</div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 5, marginTop: 2 }}>
                      <PlatformBadge platform={a.platform} t={t} />
                      <span style={{
                        fontSize: 10, color: t.textDim,
                        overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                      }}>{a.category}</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Footer: user / settings */}
          <div style={{
            padding: '10px 12px', borderTop: `1px solid ${t.divider}`,
            display: 'flex', alignItems: 'center', gap: 10,
          }}>
            <div style={{
              width: 28, height: 28, borderRadius: '50%',
              background: t.accentBg,
              color: t.accent,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 12, fontWeight: 700,
            }}>E</div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ fontSize: 12, fontWeight: 500, color: t.text, letterSpacing: -0.1 }}>
                entaku
              </div>
              <div style={{ fontSize: 10, color: t.textDim }}>Indie プラン</div>
            </div>
            <button style={iconBtnStyle(t)} title="設定">
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.4">
                <circle cx="7" cy="7" r="2.2"/>
                <path d="M7 1v1.5M7 11.5V13M13 7h-1.5M2.5 7H1M11.24 2.76l-1.06 1.06M3.82 10.18l-1.06 1.06M11.24 11.24l-1.06-1.06M3.82 3.82L2.76 2.76" strokeLinecap="round"/>
              </svg>
            </button>
          </div>
        </aside>

        {/* Main pane */}
        <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
          {/* App header bar */}
          <div style={{
            padding: '14px 28px',
            borderBottom: `1px solid ${t.divider}`,
            display: 'flex', alignItems: 'center', gap: 16,
            WebkitAppRegion: 'drag',
          }}>
            <div style={{
              width: 38, height: 38, borderRadius: 9, background: app.iconBg,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 18, fontWeight: 700, color: '#fff', flexShrink: 0,
              boxShadow: '0 1px 3px rgba(0,0,0,0.12), inset 0 0 0 0.5px rgba(255,255,255,0.2)',
            }}>{app.icon}</div>
            <div style={{ minWidth: 0 }}>
              <div style={{ fontSize: 16, fontWeight: 600, color: t.text, letterSpacing: -0.2 }}>{app.name}</div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginTop: 2 }}>
                <PlatformBadge platform={app.platform} t={t} />
                <span style={{ fontSize: 11, color: t.textMuted, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace' }}>
                  {app.bundleID}
                </span>
              </div>
            </div>
            <div style={{ flex: 1 }} />
            {/* Rating */}
            <div style={{
              display: 'flex', alignItems: 'center', gap: 6,
              padding: '6px 10px', borderRadius: 8,
              background: t.surfaceSubtle, border: `1px solid ${t.border}`,
            }}>
              <svg width="12" height="12" viewBox="0 0 12 12" fill="#febc2e">
                <path d="M6 1l1.5 3.2 3.5.4-2.6 2.4.7 3.5L6 8.8 2.9 10.5l.7-3.5L1 4.6l3.5-.4z"/>
              </svg>
              <span style={{
                fontSize: 12, fontWeight: 600, color: t.text, fontVariantNumeric: 'tabular-nums',
              }}>{app.rating}</span>
              <span style={{ fontSize: 11, color: t.textDim, fontVariantNumeric: 'tabular-nums' }}>
                ({app.reviews.toLocaleString()})
              </span>
            </div>
            <button style={{
              padding: '7px 14px', borderRadius: 8, border: `1px solid ${t.border}`,
              background: t.btnGhostBg, color: t.text, fontSize: 12, fontWeight: 500,
              cursor: 'pointer', fontFamily: FONT_STACK,
              display: 'inline-flex', alignItems: 'center', gap: 5,
            }}>
              <svg width="11" height="11" viewBox="0 0 11 11" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
                <circle cx="5.5" cy="5.5" r="4"/>
                <path d="M5.5 3v2.5l1.5 1"/>
              </svg>
              5分前に更新
            </button>
          </div>

          {/* Tab bar */}
          <div style={{
            padding: '0 28px', borderBottom: `1px solid ${t.divider}`,
            display: 'flex', alignItems: 'center', gap: 4,
          }}>
            {[
              { id: 'keywords', label: 'キーワード', count: (window.ASO_KEYWORDS_BY_APP[app.id] || []).length },
              { id: 'gap',      label: '競合ギャップ', count: (window.ASO_GAP_BY_APP[app.id] || []).length },
              { id: 'metadata', label: 'メタデータ', count: (window.ASO_METADATA_BY_APP[app.id] || []).length },
            ].map(item => {
              const active = item.id === tab;
              return (
                <button key={item.id} onClick={() => setTab(item.id)} style={{
                  border: 'none', background: 'transparent', cursor: 'pointer',
                  padding: '12px 14px', position: 'relative',
                  fontFamily: FONT_STACK, fontSize: 13,
                  fontWeight: active ? 600 : 500,
                  color: active ? t.text : t.textMuted,
                  display: 'inline-flex', alignItems: 'center', gap: 6,
                  letterSpacing: -0.1,
                  transition: 'color 120ms',
                }}>
                  {item.label}
                  <span style={{
                    fontSize: 10, fontWeight: 600,
                    padding: '1px 5px', borderRadius: 3,
                    background: active ? t.accentBg : t.tabCountBg,
                    color: active ? t.accent : t.textDim,
                    fontVariantNumeric: 'tabular-nums',
                  }}>{item.count}</span>
                  {active && <div style={{
                    position: 'absolute', left: 14, right: 14, bottom: -1, height: 2,
                    background: t.accent, borderRadius: 1,
                  }} />}
                </button>
              );
            })}
          </div>

          {/* Tab content */}
          <div style={{ flex: 1, minHeight: 0, overflow: 'hidden' }}>
            {tab === 'keywords' && <KeywordsPane app={app} t={t} onOpenAdd={() => setShowAdd(true)} />}
            {tab === 'gap' && <GapPane app={app} t={t} />}
            {tab === 'metadata' && <MetadataPane app={app} t={t} />}
          </div>
        </main>
      </div>

      {/* Modal */}
      {showAdd && <AddKeywordSheet app={app} t={t} onClose={() => setShowAdd(false)} />}
    </div>
  );
}

Object.assign(window, { GapPane, MetadataPane, AddKeywordSheet, AsoApp, textInputStyle, selectStyle });
