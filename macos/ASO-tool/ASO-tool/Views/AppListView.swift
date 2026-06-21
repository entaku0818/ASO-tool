//  AppListView.swift
import SwiftUI

struct AppListView: View {
    @EnvironmentObject var appState: AppState
    @Binding var selectedApp: ASOApp?
    @Environment(\.colorScheme) var colorScheme
    @State private var apps: [ASOApp] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var searchText = ""
    @State private var showingAddApp = false

    private var filtered: [ASOApp] {
        guard !searchText.isEmpty else { return apps }
        return apps.filter { $0.name.localizedCaseInsensitiveContains(searchText) || $0.bundleID.localizedCaseInsensitiveContains(searchText) }
    }

    var body: some View {
        ZStack {
            if colorScheme == .dark {
                Aurora.sidebarBg.ignoresSafeArea()
                AuroraGradient()
            }

            VStack(spacing: 0) {
                // Section header
                HStack {
                    VStack(alignment: .leading, spacing: 1) {
                        Text("ASO Studio")
                            .font(.system(size: 10, weight: .bold))
                            .foregroundStyle(colorScheme == .dark ? AnyShapeStyle(Aurora.textDim) : AnyShapeStyle(.tertiary))
                            .tracking(1.2)
                            .textCase(.uppercase)
                        Text("マイアプリ")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    }
                    Spacer()
                    Button(action: load) {
                        Image(systemName: "arrow.clockwise")
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    }
                    .buttonStyle(.plain)
                    .disabled(isLoading)
                    Button { showingAddApp = true } label: {
                        Image(systemName: "plus")
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    }
                    .buttonStyle(.plain)
                }
                .padding(.horizontal, 16)
                .padding(.top, 14)
                .padding(.bottom, 10)

                // Search
                HStack(spacing: 6) {
                    Image(systemName: "magnifyingglass")
                        .font(.system(size: 11))
                        .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    TextField("アプリを検索...", text: $searchText)
                        .font(.system(size: 12))
                        .textFieldStyle(.plain)
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                }
                .padding(.horizontal, 10)
                .padding(.vertical, 7)
                .background(
                    RoundedRectangle(cornerRadius: 8)
                        .fill(colorScheme == .dark ? Color.white.opacity(0.05) : Color.primary.opacity(0.06))
                )
                .padding(.horizontal, 12)
                .padding(.bottom, 8)

                // Section label
                Text("追跡中")
                    .font(.system(size: 10, weight: .bold))
                    .foregroundStyle(colorScheme == .dark ? AnyShapeStyle(Aurora.textDim) : AnyShapeStyle(.tertiary))
                    .tracking(0.4)
                    .textCase(.uppercase)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal, 16)
                    .padding(.bottom, 6)

                // App list
                if isLoading {
                    Spacer()
                    ProgressView()
                    Spacer()
                } else if let err = errorMessage {
                    Spacer()
                    Text(err).font(.caption).foregroundStyle(.red).padding()
                    Spacer()
                } else if filtered.isEmpty {
                    Spacer()
                    ContentUnavailableView(searchText.isEmpty ? "アプリなし" : "見つかりません", systemImage: "square.dashed")
                    Spacer()
                } else {
                    ScrollView {
                        LazyVStack(spacing: 2) {
                            ForEach(filtered, id: \.id) { app in
                                AppSidebarRow(app: app, isSelected: selectedApp?.id == app.id)
                                    .onTapGesture { selectedApp = app }
                                    .contextMenu {
                                        Button(role: .destructive) {
                                            deleteApp(app)
                                        } label: {
                                            Label("削除", systemImage: "trash")
                                        }
                                    }
                            }
                        }
                        .padding(.horizontal, 8)
                        .padding(.bottom, 12)
                    }
                }

                Divider()
                    .opacity(colorScheme == .dark ? 0.15 : 0.3)

                // Footer
                HStack(spacing: 10) {
                    ZStack {
                        Circle()
                            .fill(colorScheme == .dark ? Aurora.accentBg : Color.accentColor.opacity(0.12))
                            .frame(width: 28, height: 28)
                        Text(String(appState.email.prefix(1)).uppercased())
                            .font(.system(size: 12, weight: .bold))
                            .foregroundStyle(colorScheme == .dark ? Aurora.accent : .accentColor)
                    }
                    VStack(alignment: .leading, spacing: 1) {
                        Text(appState.email.isEmpty ? "ユーザー" : appState.email)
                            .font(.system(size: 12, weight: .medium))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                            .lineLimit(1)
                        Text("Indie プラン")
                            .font(.system(size: 10))
                            .foregroundStyle(colorScheme == .dark ? AnyShapeStyle(Aurora.textDim) : AnyShapeStyle(.tertiary))
                    }
                    Spacer()
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 10)
            }
        }
        .navigationSplitViewColumnWidth(min: 210, ideal: 248)
        .task { await fetch() }
        .sheet(isPresented: $showingAddApp) {
            AddAppView { newApp in
                apps.append(newApp)
                selectedApp = newApp
            }
            .environmentObject(appState)
        }
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

    private func deleteApp(_ app: ASOApp) {
        Task {
            try? await APIClient.shared.deleteApp(token: appState.token, appID: app.id)
            await MainActor.run {
                apps.removeAll { $0.id == app.id }
                if selectedApp?.id == app.id { selectedApp = nil }
            }
        }
    }
}

struct AppSidebarRow: View {
    let app: ASOApp
    let isSelected: Bool
    @Environment(\.colorScheme) var colorScheme
    @State private var isHovered = false

    private var appColors: [(Color, Color)] {
        [
            (Color(hex: "5B8FEF"), .white), (Color(hex: "E8624A"), .white),
            (Color(hex: "48C78E"), .white), (Color(hex: "A78BFA"), .white),
            (Color(hex: "F59E0B"), .white), (Color(hex: "06B6D4"), .white),
        ]
    }

    private func iconColor(for id: String) -> (Color, Color) {
        let idx = abs(id.hashValue) % appColors.count
        return appColors[idx]
    }

    var body: some View {
        HStack(spacing: 10) {
            // Left accent bar for selected dark mode
            if isSelected && colorScheme == .dark {
                RoundedRectangle(cornerRadius: 1)
                    .fill(Aurora.accent)
                    .frame(width: 2, height: 28)
            }

            // App icon
            let (bgColor, fgColor) = iconColor(for: app.id)
            ZStack {
                RoundedRectangle(cornerRadius: 8)
                    .fill(bgColor)
                    .frame(width: 32, height: 32)
                Text(String(app.name.prefix(1)))
                    .font(.system(size: 15, weight: .bold))
                    .foregroundStyle(fgColor)
            }

            VStack(alignment: .leading, spacing: 2) {
                Text(app.name)
                    .font(.system(size: 13, weight: isSelected ? .semibold : .medium))
                    .foregroundStyle(
                        colorScheme == .dark
                            ? (isSelected ? Aurora.rowSelText : Aurora.text)
                            : (isSelected ? Color.accentColor : .primary)
                    )
                    .lineLimit(1)
                HStack(spacing: 5) {
                    // Platform badge
                    Text(app.platform == "ios" ? "iOS" : "Android")
                        .font(.system(size: 9, weight: .semibold))
                        .padding(.horizontal, 5)
                        .padding(.vertical, 2)
                        .background(
                            colorScheme == .dark
                                ? Color.white.opacity(0.08)
                                : Color.primary.opacity(0.08)
                        )
                        .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                        .clipShape(RoundedRectangle(cornerRadius: 3))
                    Text(app.bundleID)
                        .font(.system(size: 10))
                        .foregroundStyle(colorScheme == .dark ? AnyShapeStyle(Aurora.textDim) : AnyShapeStyle(.tertiary))
                        .lineLimit(1)
                        .truncationMode(.middle)
                }
            }
        }
        .padding(.horizontal, 8)
        .padding(.vertical, 8)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(
            RoundedRectangle(cornerRadius: 9)
                .fill(
                    isSelected
                        ? (colorScheme == .dark ? Aurora.rowSel : Color.accentColor.opacity(0.12))
                        : (isHovered ? (colorScheme == .dark ? Color.white.opacity(0.04) : Color.primary.opacity(0.04)) : .clear)
                )
        )
        .onHover { isHovered = $0 }
        .animation(.easeOut(duration: 0.12), value: isHovered)
        .animation(.easeOut(duration: 0.12), value: isSelected)
    }
}
