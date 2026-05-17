//  MainView.swift

import SwiftUI

struct MainView: View {
    @State private var selectedApp: ASOApp?
    @Environment(\.colorScheme) var colorScheme

    var body: some View {
        NavigationSplitView {
            AppListView(selectedApp: $selectedApp)
        } detail: {
            if let app = selectedApp {
                AppDetailView(app: app)
            } else {
                ContentUnavailableView(
                    "アプリを選択",
                    systemImage: "square.dashed",
                    description: Text("左のリストからアプリを選んでください")
                )
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
                .overlay(colorScheme == .dark ? AuroraGradient() : nil)
            }
        }
        .tint(colorScheme == .dark ? Aurora.accent : Daylight.accent)
    }
}
