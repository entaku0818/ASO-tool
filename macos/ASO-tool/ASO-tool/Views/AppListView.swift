//  AppListView.swift

import SwiftUI

struct AppListView: View {
    @EnvironmentObject var appState: AppState
    @Binding var selectedApp: ASOApp?
    @State private var apps: [ASOApp] = []
    @State private var isLoading = false
    @State private var errorMessage: String?

    var body: some View {
        List(apps, id: \.id, selection: $selectedApp) { app in
            VStack(alignment: .leading, spacing: 4) {
                Text(app.name).font(.headline)
                Text(app.bundleID).font(.caption).foregroundStyle(.secondary)
                Text(app.platform == "ios" ? "iOS" : "Android")
                    .font(.caption2)
                    .padding(.horizontal, 6)
                    .padding(.vertical, 2)
                    .background(.secondary.opacity(0.15))
                    .clipShape(Capsule())
            }
            .padding(.vertical, 4)
            .tag(app)
        }
        .navigationSplitViewColumnWidth(min: 200, ideal: 240)
        .navigationTitle("アプリ")
        .toolbar {
            ToolbarItem {
                Button(action: load) {
                    Label("更新", systemImage: "arrow.clockwise")
                }
                .disabled(isLoading)
            }
        }
        .overlay {
            if isLoading {
                ProgressView()
            } else if let err = errorMessage {
                ContentUnavailableView(err, systemImage: "exclamationmark.triangle")
            } else if apps.isEmpty {
                ContentUnavailableView("アプリなし", systemImage: "square.dashed")
            }
        }
        .task { await fetch() }
    }

    private func load() { Task { await fetch() } }

    private func fetch() async {
        isLoading = true
        errorMessage = nil
        do {
            apps = try await APIClient.shared.getApps(token: appState.token)
        } catch {
            errorMessage = error.localizedDescription
        }
        isLoading = false
    }
}
