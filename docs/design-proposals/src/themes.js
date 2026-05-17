// Two themes for the ASO redesign — Aurora (dark) + Daylight (light).
// Inspired by tryastro.app's premium-mac feel: generous whitespace,
// soft glow, sophisticated palette.

// ─────────────────────────────────────────────────────────────
// AURORA — Deep midnight + aurora glow + mint accent
// ─────────────────────────────────────────────────────────────
const AURORA_THEME = {
  name: 'aurora',
  dark: true,

  // Window
  windowBg: '#0d0f14',
  windowBorder: '1px solid rgba(255,255,255,0.08)',
  windowShadow: '0 30px 80px rgba(0,0,0,0.55), 0 0 0 1px rgba(0,0,0,0.4)',

  // Ambient aurora — soft mint/lavender/coral gradients
  aurora: `
    radial-gradient(900px 500px at 5% 0%,    rgba(102, 232, 198, 0.10), transparent 60%),
    radial-gradient(800px 600px at 95% 100%, rgba(167, 139, 250, 0.10), transparent 55%),
    radial-gradient(700px 500px at 60% 30%,  rgba(94, 178, 255, 0.04),  transparent 60%)
  `,

  // Sidebar
  sidebarBg: 'rgba(18, 22, 30, 0.7)',
  sidebarSelectedBg: 'rgba(255,255,255,0.06)',
  sidebarSelectedAccent: true,

  // Text
  text: 'rgba(240, 243, 248, 0.95)',
  textMuted: 'rgba(180, 188, 200, 0.85)',
  textDim: 'rgba(140, 148, 160, 0.7)',

  // Borders
  border: 'rgba(255, 255, 255, 0.08)',
  divider: 'rgba(255, 255, 255, 0.06)',

  // Surfaces
  surface: 'rgba(255, 255, 255, 0.03)',
  surfaceSubtle: 'rgba(255, 255, 255, 0.025)',

  // Accent — soft mint
  accent: '#66E8C6',
  accentText: '#0a1410',
  accentBg: 'rgba(102, 232, 198, 0.14)',

  // Inputs
  inputBg: 'rgba(255, 255, 255, 0.04)',

  // Rows
  rowSelectedBg: 'rgba(102, 232, 198, 0.10)',
  rowSelectedText: '#A8F0DC',
  rowHoverBg: 'rgba(255,255,255,0.04)',

  // Platform badges
  platformIosBg: 'rgba(255,255,255,0.08)',
  platformIosText: 'rgba(240,243,248,0.85)',
  platformAndBg: 'rgba(167,232,140,0.14)',
  platformAndText: '#B8E89A',

  // Status colors
  posText: '#5DDFAE',
  posBg: 'rgba(93, 223, 174, 0.14)',
  negText: '#FF8A8A',
  negBg: 'rgba(255, 138, 138, 0.14)',
  warnText: '#FFB86B',
  warnBg: 'rgba(255, 184, 107, 0.14)',

  // Segmented
  segBg: 'rgba(255,255,255,0.03)',
  segActiveBg: 'rgba(255,255,255,0.10)',
  segActiveText: 'rgba(240,243,248,0.95)',
  segActiveShadow: '0 1px 2px rgba(0,0,0,0.3), inset 0 0 0 0.5px rgba(255,255,255,0.1)',

  // Misc
  tabCountBg: 'rgba(255,255,255,0.06)',
  gridLine: 'rgba(255,255,255,0.08)',
  tableHeadBg: 'rgba(255,255,255,0.02)',
  compIcon: 'rgba(167, 139, 250, 0.7)',

  tagBg: 'rgba(102, 232, 198, 0.16)',
  tagText: '#A8F0DC',

  btnGhostBg: 'rgba(255,255,255,0.04)',
  btnDisabledBg: 'rgba(255,255,255,0.04)',
  btnPrimaryShadow: '0 2px 10px rgba(102, 232, 198, 0.25), inset 0 0 0 0.5px rgba(255,255,255,0.15)',

  modalScrim: 'rgba(5, 7, 12, 0.55)',
  modalBg: '#16191F',
  modalShadow: '0 30px 80px rgba(0,0,0,0.6), 0 0 0 1px rgba(255,255,255,0.06)',

  // Web-specific
  browserChromeBg: 'rgba(14, 17, 22, 0.85)',
  browserUrlBg: 'rgba(255,255,255,0.03)',
  brandSecondary: '#A78BFA',
};

// ─────────────────────────────────────────────────────────────
// DAYLIGHT — Warm off-white + soft coral accent
// ─────────────────────────────────────────────────────────────
const DAYLIGHT_THEME = {
  name: 'daylight',
  dark: false,

  // Window
  windowBg: '#FAF9F6',
  windowBorder: '1px solid rgba(0,0,0,0.06)',
  windowShadow: '0 24px 60px rgba(60, 50, 40, 0.18), 0 0 0 1px rgba(0,0,0,0.04)',

  // Soft daylight gradient — peach + lavender
  aurora: `
    radial-gradient(900px 500px at 5% 0%,    rgba(255, 196, 161, 0.22), transparent 65%),
    radial-gradient(800px 600px at 95% 100%, rgba(196, 173, 255, 0.18), transparent 60%),
    radial-gradient(700px 500px at 50% 50%,  rgba(255, 224, 196, 0.10), transparent 70%)
  `,

  // Sidebar
  sidebarBg: 'rgba(255, 252, 247, 0.65)',
  sidebarSelectedBg: 'rgba(255, 255, 255, 0.85)',
  sidebarSelectedAccent: true,

  // Text
  text: '#1a1816',
  textMuted: '#6b665e',
  textDim: '#9a948a',

  // Borders
  border: 'rgba(60, 50, 40, 0.10)',
  divider: 'rgba(60, 50, 40, 0.07)',

  // Surfaces
  surface: 'rgba(255, 255, 255, 0.7)',
  surfaceSubtle: 'rgba(255, 255, 255, 0.55)',

  // Accent — warm coral
  accent: '#F17E63',
  accentText: '#FFFFFF',
  accentBg: 'rgba(241, 126, 99, 0.12)',

  // Inputs
  inputBg: 'rgba(255, 255, 255, 0.7)',

  // Rows
  rowSelectedBg: 'rgba(241, 126, 99, 0.10)',
  rowSelectedText: '#B0533D',
  rowHoverBg: 'rgba(60, 50, 40, 0.04)',

  // Platform badges
  platformIosBg: 'rgba(60, 50, 40, 0.07)',
  platformIosText: '#3a3530',
  platformAndBg: 'rgba(120, 180, 100, 0.14)',
  platformAndText: '#3F7430',

  // Status
  posText: '#1F8A5B',
  posBg: 'rgba(31, 138, 91, 0.10)',
  negText: '#C24545',
  negBg: 'rgba(194, 69, 69, 0.10)',
  warnText: '#B07418',
  warnBg: 'rgba(176, 116, 24, 0.12)',

  // Segmented
  segBg: 'rgba(60, 50, 40, 0.05)',
  segActiveBg: '#FFFFFF',
  segActiveText: '#1a1816',
  segActiveShadow: '0 1px 2px rgba(60,50,40,0.10), 0 0 0 0.5px rgba(60,50,40,0.06)',

  // Misc
  tabCountBg: 'rgba(60, 50, 40, 0.06)',
  gridLine: 'rgba(60, 50, 40, 0.08)',
  tableHeadBg: 'rgba(255, 252, 247, 0.6)',
  compIcon: 'rgba(196, 173, 255, 0.85)',

  tagBg: 'rgba(60, 50, 40, 0.08)',
  tagText: '#3a3530',

  btnGhostBg: '#FFFFFF',
  btnDisabledBg: 'rgba(60, 50, 40, 0.06)',
  btnPrimaryShadow: '0 2px 8px rgba(241, 126, 99, 0.30), inset 0 -1px 0 rgba(0,0,0,0.06)',

  modalScrim: 'rgba(60, 50, 40, 0.25)',
  modalBg: '#FFFFFF',
  modalShadow: '0 30px 80px rgba(60, 50, 40, 0.25), 0 0 0 1px rgba(60, 50, 40, 0.06)',

  // Web-specific
  browserChromeBg: 'rgba(252, 248, 240, 0.85)',
  browserUrlBg: 'rgba(60, 50, 40, 0.04)',
  brandSecondary: '#FFB36B',
};

Object.assign(window, { AURORA_THEME, DAYLIGHT_THEME });
