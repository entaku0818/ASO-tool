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
                AppDetailView(app: app)
            } else {
                DashboardView(selectedApp: $selectedApp)
                    .environmentObject(appState)
            }
        }
        .tint(colorScheme == .dark ? Aurora.accent : Daylight.accent)
    }
}
