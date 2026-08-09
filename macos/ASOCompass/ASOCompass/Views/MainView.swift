//  MainView.swift

import SwiftUI

struct MainView: View {
    @EnvironmentObject var appState: AppState
    @State private var selectedApp: ASOApp?
    @Environment(\.colorScheme) var colorScheme

    var body: some View {
        NavigationSplitView {
            AppListView(selectedApp: $selectedApp)
        } detail: {
            if let app = selectedApp {
                // アプリを切り替えたらビューを作り直す。
                // これがないと SwiftUI が同一ビューを再利用し、keywords などの
                // @State と .task が前のアプリのまま残る。
                AppDetailView(app: app)
                    .id(app.id)
            } else {
                DashboardView(selectedApp: $selectedApp)
                    .environmentObject(appState)
            }
        }
        .tint(colorScheme == .dark ? Aurora.accent : Daylight.accent)
    }
}
