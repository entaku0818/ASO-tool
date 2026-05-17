// AuroraTheme.swift
// NOTE: Add this file to the Xcode project target (drag into Xcode navigator)

import SwiftUI

// MARK: - Hex color init
extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default: (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(.sRGB, red: Double(r)/255, green: Double(g)/255, blue: Double(b)/255, opacity: Double(a)/255)
    }
}

// MARK: - Aurora (dark) tokens
enum Aurora {
    static let windowBg   = Color(hex: "0d0f14")
    static let sidebarBg  = Color(hex: "12161e")
    static let accent     = Color(hex: "66E8C6")
    static let accentText = Color(hex: "0a1410")
    static let accentBg   = Color(hex: "66E8C6").opacity(0.14)
    static let surface    = Color.white.opacity(0.03)
    static let border     = Color.white.opacity(0.08)
    static let divider    = Color.white.opacity(0.06)
    static let text       = Color(hex: "F0F3F8")
    static let textMuted  = Color(hex: "B4BCC8")
    static let textDim    = Color(hex: "8C94A0")
    static let rowSel     = Color(hex: "66E8C6").opacity(0.10)
    static let rowSelText = Color(hex: "A8F0DC")
    static let posText    = Color(hex: "5DDFAE")
    static let negText    = Color(hex: "FF8A8A")
    static let warnText   = Color(hex: "FFB86B")
    static let tagBg      = Color(hex: "66E8C6").opacity(0.16)
    static let tagText    = Color(hex: "A8F0DC")
}

// MARK: - Daylight (light) tokens
enum Daylight {
    static let accent     = Color(hex: "F17E63")
    static let accentText = Color.white
    static let accentBg   = Color(hex: "F17E63").opacity(0.12)
    static let posText    = Color(hex: "1F8A5B")
    static let negText    = Color(hex: "C24545")
    static let warnText   = Color(hex: "B07418")
}

// MARK: - Adaptive helpers
extension Color {
    static func asoAccent(_ s: ColorScheme) -> Color      { s == .dark ? Aurora.accent : Daylight.accent }
    static func asoAccentBg(_ s: ColorScheme) -> Color    { s == .dark ? Aurora.accentBg : Daylight.accentBg }
    static func asoAccentText(_ s: ColorScheme) -> Color  { s == .dark ? Aurora.accentText : Daylight.accentText }
    static func asoPosText(_ s: ColorScheme) -> Color     { s == .dark ? Aurora.posText : Daylight.posText }
    static func asoNegText(_ s: ColorScheme) -> Color     { s == .dark ? Aurora.negText : Daylight.negText }
    static func asoWarnText(_ s: ColorScheme) -> Color    { s == .dark ? Aurora.warnText : Daylight.warnText }
    static func asoTagBg(_ s: ColorScheme) -> Color       { s == .dark ? Aurora.tagBg : Color.primary.opacity(0.08) }
    static func asoTagText(_ s: ColorScheme) -> Color     { s == .dark ? Aurora.tagText : Color.primary.opacity(0.75) }
}

// MARK: - Aurora gradient overlay
struct AuroraGradient: View {
    var body: some View {
        ZStack {
            RadialGradient(
                colors: [Color(hex: "66E8C6").opacity(0.10), .clear],
                center: UnitPoint(x: 0.05, y: 0), startRadius: 0, endRadius: 500
            )
            RadialGradient(
                colors: [Color(hex: "A78BFA").opacity(0.10), .clear],
                center: UnitPoint(x: 0.95, y: 1.0), startRadius: 0, endRadius: 450
            )
            RadialGradient(
                colors: [Color(hex: "5EB2FF").opacity(0.04), .clear],
                center: UnitPoint(x: 0.60, y: 0.30), startRadius: 0, endRadius: 400
            )
        }
        .allowsHitTesting(false)
        .ignoresSafeArea()
    }
}
