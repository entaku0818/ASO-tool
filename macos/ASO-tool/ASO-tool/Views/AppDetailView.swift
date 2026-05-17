//  AppDetailView.swift

import SwiftUI

struct AppDetailView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp
    @Environment(\.colorScheme) var colorScheme

    enum Tab: String, CaseIterable {
        case keywords = "キーワード"
        case gap      = "競合ギャップ"
        case metadata = "メタデータ"
    }

    @State private var selectedTab: Tab = .keywords
    @State private var keywords: [Keyword] = []
    @State private var selectedKeyword: Keyword?
    @State private var showAddKeyword = false
    @State private var isLoading = false

    var body: some View {
        VStack(spacing: 0) {
            // App header bar
            HStack(spacing: 14) {
                // App icon placeholder
                ZStack {
                    RoundedRectangle(cornerRadius: 9)
                        .fill(Color(hex: "5B8FEF"))
                        .frame(width: 36, height: 36)
                    Text(String(app.name.prefix(1)))
                        .font(.system(size: 17, weight: .bold))
                        .foregroundStyle(.white)
                }

                VStack(alignment: .leading, spacing: 2) {
                    Text(app.name)
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    HStack(spacing: 6) {
                        Text(app.platform == "ios" ? "iOS" : "Android")
                            .font(.system(size: 9, weight: .semibold))
                            .padding(.horizontal, 5)
                            .padding(.vertical, 2)
                            .background(colorScheme == .dark ? Color.white.opacity(0.08) : Color.primary.opacity(0.08))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                            .clipShape(RoundedRectangle(cornerRadius: 3))
                        Text(app.bundleID)
                            .font(.system(size: 11, design: .monospaced))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                }
                Spacer()
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 12)
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))

            // Tab bar — underline indicator style
            HStack(spacing: 0) {
                ForEach(Tab.allCases, id: \.self) { tab in
                    let isActive = selectedTab == tab
                    let count: Int = {
                        switch tab {
                        case .keywords: return keywords.count
                        default: return 0
                        }
                    }()

                    Button {
                        selectedTab = tab
                    } label: {
                        VStack(spacing: 0) {
                            HStack(spacing: 6) {
                                Text(tab.rawValue)
                                    .font(.system(size: 13, weight: isActive ? .semibold : .medium))
                                    .foregroundStyle(
                                        colorScheme == .dark
                                            ? (isActive ? Aurora.text : Aurora.textMuted)
                                            : (isActive ? Color.primary : Color.secondary)
                                    )

                                if tab == .keywords {
                                    Text("\(count)")
                                        .font(.system(size: 10, weight: .semibold))
                                        .monospacedDigit()
                                        .padding(.horizontal, 5)
                                        .padding(.vertical, 1)
                                        .background(
                                            RoundedRectangle(cornerRadius: 3)
                                                .fill(
                                                    isActive
                                                        ? (colorScheme == .dark ? Aurora.accentBg : Color.accentColor.opacity(0.12))
                                                        : (colorScheme == .dark ? Color.white.opacity(0.06) : Color.primary.opacity(0.06))
                                                )
                                        )
                                        .foregroundStyle(
                                            isActive
                                                ? (colorScheme == .dark ? Aurora.accent : Color.accentColor)
                                                : (colorScheme == .dark ? Aurora.textDim : Color.secondary)
                                        )
                                }
                            }
                            .padding(.horizontal, 14)
                            .padding(.vertical, 11)

                            // Active underline
                            Rectangle()
                                .fill(
                                    isActive
                                        ? (colorScheme == .dark ? Aurora.accent : Color.accentColor)
                                        : Color.clear
                                )
                                .frame(height: 2)
                        }
                    }
                    .buttonStyle(.plain)
                    .animation(.easeOut(duration: 0.12), value: isActive)
                }
                Spacer()
            }
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))

            // Divider below tabs
            Rectangle()
                .fill(colorScheme == .dark ? Aurora.divider : Color.primary.opacity(0.08))
                .frame(height: 1)

            // Tab content
            switch selectedTab {
            case .keywords:
                keywordsPane
            case .gap:
                KeywordGapView(app: app)
                    .environmentObject(appState)
            case .metadata:
                MetadataView(app: app)
                    .environmentObject(appState)
            }
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
        .navigationTitle(app.name)
        .sheet(isPresented: $showAddKeyword) {
            AddKeywordView(app: app) { newKw in
                keywords.append(newKw)
            }
            .environmentObject(appState)
        }
    }

    // MARK: - Keywords pane

    private var keywordsPane: some View {
        HSplitView {
            // Left panel
            VStack(spacing: 0) {
                // Header
                HStack {
                    Text("キーワード")
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    Spacer()
                    Button {
                        showAddKeyword = true
                    } label: {
                        Image(systemName: "plus")
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    }
                    .buttonStyle(.plain)

                    Button {
                        deleteSelected()
                    } label: {
                        Image(systemName: "trash")
                            .font(.system(size: 12))
                            .foregroundStyle(selectedKeyword != nil ? (colorScheme == .dark ? Aurora.negText : .red) : (colorScheme == .dark ? Aurora.textDim : Color.secondary))
                    }
                    .buttonStyle(.plain)
                    .disabled(selectedKeyword == nil)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 9)
                .background(colorScheme == .dark ? Aurora.surface : Color.primary.opacity(0.02))

                Rectangle()
                    .fill(colorScheme == .dark ? Aurora.divider : Color.primary.opacity(0.08))
                    .frame(height: 1)

                if isLoading {
                    ProgressView().frame(maxWidth: .infinity, maxHeight: .infinity)
                } else if keywords.isEmpty {
                    ContentUnavailableView("キーワードなし", systemImage: "magnifyingglass")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else {
                    ScrollView {
                        LazyVStack(spacing: 1) {
                            ForEach(keywords, id: \.id) { kw in
                                KeywordRow(keyword: kw, isSelected: selectedKeyword?.id == kw.id)
                                    .onTapGesture { selectedKeyword = kw }
                            }
                        }
                        .padding(.vertical, 6)
                        .padding(.horizontal, 6)
                    }
                }
            }
            .frame(minWidth: 220, idealWidth: 240, maxWidth: 300)
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
            .task { await fetchKeywords() }

            // Right panel
            if let kw = selectedKeyword {
                RankingChartView(app: app, keyword: kw)
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ContentUnavailableView(
                    "キーワードを選択",
                    systemImage: "chart.line.uptrend.xyaxis",
                    description: Text("左からキーワードを選ぶとランキング推移を表示します")
                )
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
            }
        }
    }

    private func fetchKeywords() async {
        isLoading = true
        keywords = (try? await APIClient.shared.getKeywords(token: appState.token, appID: app.id)) ?? []
        isLoading = false
    }

    private func deleteSelected() {
        guard let kw = selectedKeyword else { return }
        Task {
            try? await APIClient.shared.deleteKeyword(token: appState.token, appID: app.id, keywordID: kw.id)
            keywords.removeAll { $0.id == kw.id }
            selectedKeyword = nil
        }
    }
}

struct KeywordRow: View {
    let keyword: Keyword
    let isSelected: Bool
    @Environment(\.colorScheme) var colorScheme
    @State private var isHovered = false

    var body: some View {
        HStack(spacing: 8) {
            VStack(alignment: .leading, spacing: 3) {
                Text(keyword.keyword)
                    .font(.system(size: 13, weight: isSelected ? .semibold : .regular))
                    .foregroundStyle(
                        colorScheme == .dark
                            ? (isSelected ? Aurora.rowSelText : Aurora.text)
                            : .primary
                    )
                HStack(spacing: 6) {
                    Text(keyword.country)
                        .font(.system(size: 11))
                        .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    if let score = keyword.popularityScore {
                        popularityBadge(score)
                    }
                }
            }
            Spacer()
        }
        .padding(.horizontal, 10)
        .padding(.vertical, 8)
        .background(
            RoundedRectangle(cornerRadius: 7)
                .fill(
                    isSelected
                        ? (colorScheme == .dark ? Aurora.rowSel : Color.accentColor.opacity(0.10))
                        : (isHovered ? (colorScheme == .dark ? Color.white.opacity(0.04) : Color.primary.opacity(0.04)) : .clear)
                )
        )
        .onHover { isHovered = $0 }
        .animation(.easeOut(duration: 0.12), value: isHovered)
        .animation(.easeOut(duration: 0.12), value: isSelected)
    }

    @ViewBuilder
    private func popularityBadge(_ score: Int) -> some View {
        let (label, color): (String, Color) = {
            switch score {
            case 4...5: return ("高", colorScheme == .dark ? Aurora.posText : Daylight.posText)
            case 2...3: return ("中", colorScheme == .dark ? Aurora.warnText : Daylight.warnText)
            default:    return ("低", colorScheme == .dark ? Aurora.negText : Daylight.negText)
            }
        }()
        Text(label)
            .font(.system(size: 9, weight: .semibold))
            .padding(.horizontal, 5)
            .padding(.vertical, 2)
            .background(color.opacity(0.16))
            .foregroundStyle(color)
            .clipShape(RoundedRectangle(cornerRadius: 3))
    }
}
