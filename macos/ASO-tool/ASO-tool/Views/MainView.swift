//  MainView.swift

import SwiftUI

struct MainView: View {
    @State private var selectedApp: ASOApp?

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
            }
        }
    }
}
