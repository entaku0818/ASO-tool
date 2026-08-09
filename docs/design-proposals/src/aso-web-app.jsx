// AsoWebApp — web (browser) variant of the ASO Compass redesign.
// Reuses the panes (KeywordsPane/GapPane/MetadataPane) from aso-app.jsx;
// replaces macOS chrome+sidebar with web-native top nav + app switcher.

function BrowserChrome({ t, url }) {
  const dot = (bg) => (
    <div style={{
      width: 11, height: 11, borderRadius: '50%', background: bg,
      boxShadow: 'inset 0 0 0 0.5px rgba(0,0,0,0.18)',
    }} />
  );
  return (
    <div style={{
      height: 38, padding: '0 14px',
      display: 'flex', alignItems: 'center', gap: 10,
      background: t.browserChromeBg,
      borderBottom: `1px solid ${t.divider}`,
      flexShrink: 0,
    }}>
      <div style={{ display: 'flex', gap: 8 }}>
        {dot('#ff5f57')}{dot('#febc2e')}{dot('#28c840')}
      </div>
      {/* Nav arrows */}
      <div style={{ display: 'flex', gap: 2, marginLeft: 12, color: t.textDim }}>
        {[
          'M9 3L4 8l5 5',
          'M5 3l5 5-5 5',
          'M3 8c0-3 3-5 5-5s5 2 5 5',
        ].map((d, i) => (
          <div key={i} style={{
            width: 22, height: 22, borderRadius: 5,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            color: t.textDim,
          }}>
            <svg width="13" height="13" viewBox="0 0 16 16" fill="none"
              stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round">
              <path d={d} />
            </svg>
          </div>
        ))}
      </div>
      {/* URL bar */}
      <div style={{
        flex: 1, height: 24, borderRadius: 12,
        background: t.browserUrlBg, border: `1px solid ${t.border}`,
        display: 'flex', alignItems: 'center', gap: 6,
        padding: '0 12px', fontSize: 11.5, color: t.textMuted,
        fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
        maxWidth: 600, margin: '0 auto',
      }}>
        <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke={t.textDim} strokeWidth="1.4">
          <rect x="3" y="5" width="6" height="5" rx="0.8"/>
          <path d="M4.5 5V3.5a1.5 1.5 0 013 0V5"/>
        </svg>
        <span style={{ color: t.textDim }}>aso-tool.entaku.app</span>
        <span style={{ color: t.text }}>{url}</span>
      </div>
      <div style={{ width: 60 }} />
    </div>
  );
}

// App switcher dropdown
function AppSwitcher({ apps, selectedId, onSelect, t }) {
  const [open, setOpen] = React.useState(false);
  const ref = React.useRef();
  React.useEffect(() => {
    const off = (e) => { if (open && ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener('pointerdown', off, true);
    return () => document.removeEventListener('pointerdown', off, true);
  }, [open]);

  const app = apps.find(a => a.id === selectedId);
  return (
    <div ref={ref} style={{ position: 'relative' }}>
      <button onClick={() => setOpen(o => !o)} style={{
        display: 'flex', alignItems: 'center', gap: 8,
        padding: '5px 10px 5px 5px', borderRadius: 8,
        background: t.surfaceSubtle, border: `1px solid ${t.border}`,
        color: t.text, fontSize: 12.5, fontWeight: 500,
        fontFamily: FONT_STACK, cursor: 'pointer',
        transition: 'background 120ms',
      }}>
        <div style={{
          width: 24, height: 24, borderRadius: 6, background: app.iconBg,
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: 12, color: '#fff', fontWeight: 700,
          boxShadow: 'inset 0 0 0 0.5px rgba(255,255,255,0.2)',
        }}>{app.icon}</div>
        <span style={{ letterSpacing: -0.1 }}>{app.name}</span>
        <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" style={{ color: t.textDim, marginLeft: 2 }}>
          <path d="M2.5 4l2.5 2.5L7.5 4"/>
        </svg>
      </button>
      {open && (
        <div style={{
          position: 'absolute', top: 'calc(100% + 6px)', left: 0,
          width: 280, zIndex: 50,
          background: t.modalBg, border: `1px solid ${t.border}`,
          borderRadius: 10, boxShadow: t.modalShadow,
          padding: 6, animation: 'asoFadeIn 140ms',
        }}>
          <div style={{
            padding: '8px 10px 6px', fontSize: 10, fontWeight: 700,
            color: t.textDim, letterSpacing: 0.4, textTransform: 'uppercase',
          }}>マイアプリ ({apps.length})</div>
          {apps.map(a => {
            const sel = a.id === selectedId;
            return (
              <div key={a.id} onClick={() => { onSelect(a.id); setOpen(false); }} style={{
                display: 'flex', alignItems: 'center', gap: 10,
                padding: '8px 10px', borderRadius: 7, cursor: 'pointer',
                background: sel ? t.rowSelectedBg : 'transparent',
                transition: 'background 120ms',
              }}
              onMouseEnter={(e) => { if (!sel) e.currentTarget.style.background = t.rowHoverBg; }}
              onMouseLeave={(e) => { if (!sel) e.currentTarget.style.background = 'transparent'; }}
              >
                <div style={{
                  width: 28, height: 28, borderRadius: 7, background: a.iconBg,
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontSize: 13, color: '#fff', fontWeight: 700,
                  boxShadow: 'inset 0 0 0 0.5px rgba(255,255,255,0.2)',
                }}>{a.icon}</div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{
                    fontSize: 12.5, fontWeight: sel ? 600 : 500,
                    color: sel ? t.rowSelectedText : t.text, letterSpacing: -0.1,
                  }}>{a.name}</div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 5, marginTop: 2 }}>
                    <PlatformBadge platform={a.platform} t={t} />
                    <span style={{ fontSize: 10, color: t.textDim }}>{a.category}</span>
                  </div>
                </div>
                {sel && (
                  <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke={t.accent} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M2.5 7.5l3 3 6-7"/>
                  </svg>
                )}
              </div>
            );
          })}
          <div style={{
            margin: '6px -6px -6px', padding: '8px 16px',
            borderTop: `1px solid ${t.divider}`,
            display: 'flex', alignItems: 'center', gap: 6,
            fontSize: 12, color: t.accent, cursor: 'pointer', fontWeight: 500,
          }}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round"><path d="M6 2v8M2 6h8"/></svg>
            新規アプリを追加
          </div>
        </div>
      )}
    </div>
  );
}

function AsoWebApp({ theme: t, width = 1440, height = 900 }) {
  const apps = window.ASO_APPS;
  const [selectedAppId, setSelectedAppId] = React.useState(apps[0].id);
  const [tab, setTab] = React.useState('keywords');
  const [showAdd, setShowAdd] = React.useState(false);

  const app = apps.find(a => a.id === selectedAppId);

  return (
    <div style={{
      width, height,
      borderRadius: 12, overflow: 'hidden',
      background: t.windowBg, color: t.text,
      fontFamily: FONT_STACK,
      boxShadow: t.windowShadow,
      border: t.windowBorder,
      display: 'flex', flexDirection: 'column',
      position: 'relative', isolation: 'isolate',
    }}>
      {/* Ambient gradient behind everything */}
      {t.aurora && <div style={{
        position: 'absolute', inset: 0, pointerEvents: 'none',
        background: t.aurora, zIndex: 0,
      }} />}

      <div style={{ position: 'relative', zIndex: 1, display: 'flex', flexDirection: 'column', height: '100%' }}>
        <BrowserChrome t={t} url={`/apps/${app.id}`} />

        {/* App nav bar (web-native) */}
        <header style={{
          height: 56, padding: '0 32px',
          display: 'flex', alignItems: 'center', gap: 16,
          background: t.sidebarBg,
          backdropFilter: 'blur(40px) saturate(180%)',
          WebkitBackdropFilter: 'blur(40px) saturate(180%)',
          borderBottom: `1px solid ${t.divider}`,
          flexShrink: 0,
        }}>
          {/* Brand */}
          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <div style={{
              width: 28, height: 28, borderRadius: 8,
              background: `linear-gradient(135deg, ${t.accent}, ${t.brandSecondary || t.accent})`,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              boxShadow: `0 1px 4px ${t.accent}44, inset 0 0 0 0.5px rgba(255,255,255,0.15)`,
            }}>
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke={t.accentText} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                <path d="M3 11l3-6 2 4 3-5"/>
                <circle cx="3" cy="11" r="0.8" fill="currentColor"/>
                <circle cx="6" cy="5" r="0.8" fill="currentColor"/>
                <circle cx="8" cy="9" r="0.8" fill="currentColor"/>
                <circle cx="11" cy="4" r="0.8" fill="currentColor"/>
              </svg>
            </div>
            <div>
              <div style={{
                fontSize: 14.5, fontWeight: 700, color: t.text,
                letterSpacing: -0.3, lineHeight: 1.1,
              }}>ASO Studio</div>
            </div>
          </div>

          {/* Divider */}
          <div style={{ width: 1, height: 22, background: t.divider }} />

          {/* App switcher */}
          <AppSwitcher apps={apps} selectedId={selectedAppId} onSelect={setSelectedAppId} t={t} />

          <div style={{ flex: 1 }} />

          {/* Global search */}
          <div style={{ width: 280 }}>
            <SearchInput value="" onChange={() => {}} placeholder="キーワード・アプリを検索..." t={t} />
          </div>

          {/* Quick actions */}
          <button title="通知" style={{
            ...iconBtnStyle(t), width: 32, height: 32,
            background: t.surfaceSubtle, border: `1px solid ${t.border}`,
            position: 'relative',
          }}>
            <svg width="15" height="15" viewBox="0 0 15 15" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round">
              <path d="M3.5 6.5c0-2.5 1.5-4 4-4s4 1.5 4 4v3l1 1.5h-10l1-1.5z"/>
              <path d="M6 12.5c0 .8.7 1.5 1.5 1.5s1.5-.7 1.5-1.5"/>
            </svg>
            <div style={{
              position: 'absolute', top: 6, right: 6,
              width: 6, height: 6, borderRadius: '50%',
              background: t.negText, border: `1.5px solid ${t.windowBg}`,
            }} />
          </button>

          {/* User menu */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: 8,
            padding: '4px 12px 4px 4px', borderRadius: 18,
            background: t.surfaceSubtle, border: `1px solid ${t.border}`,
            cursor: 'pointer',
          }}>
            <div style={{
              width: 24, height: 24, borderRadius: '50%',
              background: t.accentBg, color: t.accent,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 11, fontWeight: 700,
            }}>E</div>
            <span style={{ fontSize: 12, fontWeight: 500, color: t.text }}>entaku</span>
            <svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke={t.textDim} strokeWidth="1.6" strokeLinecap="round"><path d="M2.5 3.5L4.5 5.5 6.5 3.5"/></svg>
          </div>
        </header>

        {/* Page header / app card */}
        <div style={{ padding: '20px 32px 0' }}>
          {/* Breadcrumb */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: 8,
            fontSize: 12, color: t.textMuted, marginBottom: 14,
          }}>
            <span>マイアプリ</span>
            <svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round"><path d="M3 1.5L6 4.5 3 7.5"/></svg>
            <span style={{ color: t.text, fontWeight: 500 }}>{app.name}</span>
          </div>

          {/* App hero card */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: 20,
            padding: '18px 22px', borderRadius: 14,
            background: t.surface, border: `1px solid ${t.border}`,
            marginBottom: 16,
          }}>
            <div style={{
              width: 56, height: 56, borderRadius: 13, background: app.iconBg,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 26, color: '#fff', fontWeight: 700, flexShrink: 0,
              boxShadow: '0 2px 8px rgba(0,0,0,0.10), inset 0 0 0 0.5px rgba(255,255,255,0.2)',
            }}>{app.icon}</div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{
                display: 'flex', alignItems: 'center', gap: 10,
              }}>
                <h1 style={{
                  fontSize: 22, fontWeight: 600, color: t.text,
                  margin: 0, letterSpacing: -0.4,
                }}>{app.name}</h1>
                <PlatformBadge platform={app.platform} t={t} />
              </div>
              <div style={{
                display: 'flex', alignItems: 'center', gap: 10, marginTop: 5,
                fontSize: 12, color: t.textMuted,
              }}>
                <span style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace' }}>{app.bundleID}</span>
                <span style={{ color: t.divider }}>·</span>
                <span>{app.category}</span>
              </div>
            </div>

            {/* Quick stats */}
            <div style={{ display: 'flex', gap: 0, alignItems: 'stretch' }}>
              {(() => {
                const kwCount = (window.ASO_KEYWORDS_BY_APP[app.id] || []).length;
                const ranked = (window.ASO_KEYWORDS_BY_APP[app.id] || []).filter(k => k.currentRank != null).length;
                const top10 = (window.ASO_KEYWORDS_BY_APP[app.id] || []).filter(k => k.currentRank != null && k.currentRank <= 10).length;
                return [
                  { label: 'キーワード', value: kwCount },
                  { label: 'ランク内', value: ranked },
                  { label: 'TOP10', value: top10 },
                ].map((s, i) => (
                  <React.Fragment key={i}>
                    {i > 0 && <div style={{ width: 1, background: t.divider, margin: '4px 0' }} />}
                    <div style={{ padding: '0 18px', textAlign: 'center' }}>
                      <div style={{
                        fontSize: 20, fontWeight: 700, color: t.text,
                        fontVariantNumeric: 'tabular-nums', letterSpacing: -0.5,
                      }}>{s.value}</div>
                      <div style={{ fontSize: 10.5, color: t.textDim, fontWeight: 500, marginTop: 2, letterSpacing: 0.3, textTransform: 'uppercase' }}>{s.label}</div>
                    </div>
                  </React.Fragment>
                ));
              })()}
            </div>

            {/* Rating */}
            <div style={{
              padding: '10px 14px', borderRadius: 10,
              background: t.surfaceSubtle, border: `1px solid ${t.border}`,
              minWidth: 90,
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
                <svg width="13" height="13" viewBox="0 0 12 12" fill="#febc2e">
                  <path d="M6 1l1.5 3.2 3.5.4-2.6 2.4.7 3.5L6 8.8 2.9 10.5l.7-3.5L1 4.6l3.5-.4z"/>
                </svg>
                <span style={{
                  fontSize: 15, fontWeight: 700, color: t.text, fontVariantNumeric: 'tabular-nums',
                }}>{app.rating}</span>
              </div>
              <div style={{ fontSize: 10.5, color: t.textDim, marginTop: 2, fontVariantNumeric: 'tabular-nums' }}>
                {app.reviews.toLocaleString()} 件
              </div>
            </div>

            <button style={{
              padding: '9px 16px', borderRadius: 9, border: 'none',
              background: t.accent, color: t.accentText,
              fontSize: 12.5, fontWeight: 600, fontFamily: FONT_STACK,
              cursor: 'pointer', boxShadow: t.btnPrimaryShadow,
              display: 'inline-flex', alignItems: 'center', gap: 6,
            }}>
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
                <path d="M2 6h8M6 2v8"/>
              </svg>
              キーワード追加
            </button>
          </div>

          {/* Tab bar */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: 4,
            borderBottom: `1px solid ${t.divider}`,
            marginBottom: 0,
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
                  padding: '12px 16px', position: 'relative',
                  fontFamily: FONT_STACK, fontSize: 13.5,
                  fontWeight: active ? 600 : 500,
                  color: active ? t.text : t.textMuted,
                  display: 'inline-flex', alignItems: 'center', gap: 6,
                  letterSpacing: -0.1, transition: 'color 120ms',
                }}>
                  {item.label}
                  <span style={{
                    fontSize: 10.5, fontWeight: 600, padding: '1px 6px', borderRadius: 3,
                    background: active ? t.accentBg : t.tabCountBg,
                    color: active ? t.accent : t.textDim,
                    fontVariantNumeric: 'tabular-nums',
                  }}>{item.count}</span>
                  {active && <div style={{
                    position: 'absolute', left: 12, right: 12, bottom: -1, height: 2,
                    background: t.accent, borderRadius: 1,
                  }} />}
                </button>
              );
            })}
          </div>
        </div>

        {/* Tab content */}
        <div style={{
          flex: 1, minHeight: 0, overflow: 'hidden',
          margin: '0 32px 24px', borderRadius: '0 0 14px 14px',
          background: t.surface, border: `1px solid ${t.border}`,
          borderTopLeftRadius: 0, borderTopRightRadius: 14,
          borderTop: 'none',
          marginTop: -1,
          display: 'flex', flexDirection: 'column',
        }}>
          {tab === 'keywords' && <KeywordsPane app={app} t={t} onOpenAdd={() => setShowAdd(true)} />}
          {tab === 'gap'      && <GapPane      app={app} t={t} />}
          {tab === 'metadata' && <MetadataPane app={app} t={t} />}
        </div>
      </div>

      {showAdd && <AddKeywordSheet app={app} t={t} onClose={() => setShowAdd(false)} />}
    </div>
  );
}

Object.assign(window, { AsoWebApp, BrowserChrome, AppSwitcher });
