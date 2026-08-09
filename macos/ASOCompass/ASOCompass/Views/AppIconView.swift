//  AppIconView.swift

import Combine
import SwiftUI

/// bundle ID から App Store のアイコンURLを解決してキャッシュする。
///
/// URL は `/api/scraper/app-info` から取る。アイコン画像そのものは AsyncImage が
/// URLCache 経由でキャッシュするので、ここでキャッシュするのは URL だけでよい。
/// 解決済みURLは UserDefaults に残すので、2回目以降の起動では即座に出る。
@MainActor
final class AppIconStore: ObservableObject {
    static let shared = AppIconStore()

    private static let cacheKey = "app_icon_urls"

    @Published private(set) var urls: [String: String] = [:]

    /// 同じ bundle ID を並行して何度も引かないための進行中マーカー。
    private var inFlight: Set<String> = []
    private let defaults: UserDefaults

    init(defaults: UserDefaults = .standard) {
        self.defaults = defaults
        urls = defaults.dictionary(forKey: Self.cacheKey) as? [String: String] ?? [:]
    }

    func url(for bundleID: String) -> URL? {
        guard let s = urls[bundleID] else { return nil }
        return URL(string: s)
    }

    /// 未解決の bundle ID だけを取りに行く。解決済み・進行中のものは何もしない。
    func resolve(bundleID: String, platform: String, token: String) {
        guard platform == "ios" else { return }  // Google Play 側は未対応
        guard urls[bundleID] == nil, !inFlight.contains(bundleID), !token.isEmpty else { return }
        inFlight.insert(bundleID)

        Task { [weak self] in
            defer { Task { @MainActor in self?.inFlight.remove(bundleID) } }
            guard let info = try? await APIClient.shared.fetchAppInfo(
                token: token, bundleID: bundleID, platform: platform
            ), !info.iconURL.isEmpty else { return }

            await MainActor.run {
                guard let self else { return }
                self.urls[bundleID] = info.iconURL
                self.defaults.set(self.urls, forKey: Self.cacheKey)
            }
        }
    }
}

/// アプリアイコン。取得できるまで（および取得に失敗したとき）は
/// 従来どおり頭文字を色付きの角丸に載せたプレースホルダを出す。
struct AppIconView: View {
    let app: ASOApp
    let size: CGFloat
    var cornerRadius: CGFloat { size / 4.4 }

    @EnvironmentObject var appState: AppState
    @ObservedObject private var store = AppIconStore.shared

    /// プレースホルダの色。bundle ID から決めるので、同じアプリなら毎回同じ色になる。
    private static let palette: [Color] = [
        Color(hex: "5B8FEF"), Color(hex: "E8624A"), Color(hex: "48C78E"),
        Color(hex: "A78BFA"), Color(hex: "F59E0B"), Color(hex: "06B6D4")
    ]

    private var placeholderColor: Color {
        Self.palette[abs(app.bundleID.hashValue) % Self.palette.count]
    }

    var body: some View {
        Group {
            if let url = store.url(for: app.bundleID) {
                AsyncImage(url: url) { phase in
                    if let image = phase.image {
                        image.resizable().interpolation(.high)
                    } else {
                        placeholder
                    }
                }
            } else {
                placeholder
            }
        }
        .frame(width: size, height: size)
        .clipShape(RoundedRectangle(cornerRadius: cornerRadius, style: .continuous))
        .overlay(
            // ダークモードで白基調のアイコンが背景に溶けないよう、薄い縁を足す
            RoundedRectangle(cornerRadius: cornerRadius, style: .continuous)
                .strokeBorder(Color.primary.opacity(0.10), lineWidth: 0.5)
        )
        .task(id: app.bundleID) {
            store.resolve(bundleID: app.bundleID, platform: app.platform, token: appState.token)
        }
    }

    private var placeholder: some View {
        ZStack {
            placeholderColor
            Text(String(app.name.prefix(1)))
                .font(.system(size: size * 0.47, weight: .bold))
                .foregroundStyle(.white)
        }
    }
}
