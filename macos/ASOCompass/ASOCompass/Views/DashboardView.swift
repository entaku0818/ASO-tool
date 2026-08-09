//  DashboardView.swift

import SwiftUI

struct DashboardView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.colorScheme) var colorScheme
    @Binding var selectedApp: ASOApp?

    @State private var apps: [ASOApp] = []
    @State private var isLoading = false

    private let columns = [GridItem(.adaptive(minimum: 200, maximum: 280), spacing: 12)]

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 24) {
                // Welcome header
                VStack(alignment: .leading, spacing: 4) {
                    Text("ダッシュボード")
                        .font(.system(size: 22, weight: .bold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    Text("\(appState.email)")
                        .font(.subheadline)
                        .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                }

                // Summary cards
                HStack(spacing: 12) {
                    SummaryCard(title: "アプリ数", value: "\(apps.count)", icon: "square.stack.3d.up", color: .blue)
                    SummaryCard(title: "iOS", value: "\(apps.filter { $0.platform == "ios" }.count)", icon: "iphone", color: .indigo)
                    SummaryCard(title: "Android", value: "\(apps.filter { $0.platform == "android" }.count)", icon: "android.circle", color: .green)
                }

                // App grid
                if !apps.isEmpty {
                    VStack(alignment: .leading, spacing: 10) {
                        Text("追跡中のアプリ")
                            .font(.system(size: 13, weight: .semibold))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)

                        LazyVGrid(columns: columns, spacing: 12) {
                            ForEach(apps, id: \.id) { app in
                                AppDashboardCard(app: app, colorScheme: colorScheme) {
                                    selectedApp = app
                                }
                            }
                        }
                    }
                }
            }
            .padding(24)
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
        .overlay(colorScheme == .dark ? AuroraGradient() : nil)
        .task { await fetch() }
    }

    private func fetch() async {
        isLoading = true
        apps = (try? await APIClient.shared.getApps(token: appState.token)) ?? []
        isLoading = false
    }
}

private struct SummaryCard: View {
    let title: String
    let value: String
    let icon: String
    let color: Color
    @Environment(\.colorScheme) var colorScheme

    var body: some View {
        HStack(spacing: 10) {
            ZStack {
                RoundedRectangle(cornerRadius: 8)
                    .fill(color.opacity(colorScheme == .dark ? 0.2 : 0.12))
                    .frame(width: 36, height: 36)
                Image(systemName: icon)
                    .font(.system(size: 16))
                    .foregroundStyle(color)
            }
            VStack(alignment: .leading, spacing: 1) {
                Text(value)
                    .font(.system(size: 18, weight: .bold, design: .monospaced))
                    .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                Text(title)
                    .font(.caption)
                    .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
            }
            Spacer()
        }
        .padding(12)
        .background(
            RoundedRectangle(cornerRadius: 10)
                .fill(colorScheme == .dark ? Aurora.surface : Color.primary.opacity(0.04))
        )
    }
}

private struct AppDashboardCard: View {
    let app: ASOApp
    let colorScheme: ColorScheme
    let onOpen: () -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack(spacing: 10) {
                AppIconView(app: app, size: 36)
                VStack(alignment: .leading, spacing: 2) {
                    Text(app.name)
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                        .lineLimit(1)
                    Text(app.platform == "ios" ? "iOS" : "Android")
                        .font(.system(size: 9, weight: .semibold))
                        .padding(.horizontal, 5).padding(.vertical, 2)
                        .background(colorScheme == .dark ? Color.white.opacity(0.08) : Color.primary.opacity(0.08))
                        .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                        .clipShape(RoundedRectangle(cornerRadius: 3))
                }
                Spacer()
            }

            Text(app.bundleID)
                .font(.system(size: 10, design: .monospaced))
                .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                .lineLimit(1)

            Button("開く") { onOpen() }
                .buttonStyle(.borderless)
                .font(.system(size: 12, weight: .medium))
                .foregroundStyle(Color.accentColor)
        }
        .padding(12)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(colorScheme == .dark ? Aurora.surface : Color.primary.opacity(0.04))
        )
        .overlay(
            RoundedRectangle(cornerRadius: 12)
                .stroke(colorScheme == .dark ? Color.white.opacity(0.06) : Color.primary.opacity(0.06), lineWidth: 1)
        )
    }
}
