//  AppDetailView.swift

import SwiftUI

struct AppDetailView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp
    @Environment(\.colorScheme) var colorScheme

    enum Tab: String, CaseIterable {
        case keywords = "キーワード"
        case rising   = "急上昇"
        case gap      = "競合ギャップ"
        case metadata = "メタデータ"
    }

    @State private var selectedTab: Tab = .keywords
    @State private var keywords: [Keyword] = []
    @State private var selectedKeyword: Keyword?
    @State private var showAddKeyword = false
    @State private var showImportKeywords = false
    @State private var showSuggestions = false
    @State private var isLoading = false
    @State private var isUpdatingRankings = false
    @State private var lastUpdated: Date?
    // Multi-select
    @State private var multiSelectMode = false
    @State private var selectedKeywords: Set<String> = []
    // Sort, Filter & Search
    enum SortOrder: String, CaseIterable { case alpha = "五十音", rank = "順位", popularity = "人気度" }
    @State private var sortOrder: SortOrder = .rank
    @State private var filterCountry = "すべて"
    @State private var searchText = ""
    // Ranks & Rising
    @State private var rankSummaries: [String: KeywordRankSummary] = [:]
    @State private var risingKeywords: [RisingKeyword] = []

    private var availableCountries: [String] {
        let countries = Set(keywords.map(\.country)).sorted()
        return ["すべて"] + countries
    }

    private var displayedKeywords: [Keyword] {
        var list = keywords
        if filterCountry != "すべて" { list = list.filter { $0.country == filterCountry } }
        if !searchText.isEmpty { list = list.filter { $0.keyword.localizedCaseInsensitiveContains(searchText) } }
        switch sortOrder {
        case .alpha:
            list.sort { $0.keyword < $1.keyword }
        case .rank:
            list.sort {
                let r0 = rankSummaries[$0.id]?.currentRank ?? Int.max
                let r1 = rankSummaries[$1.id]?.currentRank ?? Int.max
                return r0 < r1
            }
        case .popularity:
            list.sort { ($0.popularityScore ?? -1) > ($1.popularityScore ?? -1) }
        }
        return list
    }

    private var statsTop10: Int { rankSummaries.values.filter { ($0.currentRank ?? Int.max) <= 10 }.count }
    private var statsTop50: Int { rankSummaries.values.filter { ($0.currentRank ?? Int.max) <= 50 }.count }
    private var statsAvgRank: Double? {
        let ranks = rankSummaries.values.compactMap(\.currentRank)
        guard !ranks.isEmpty else { return nil }
        return Double(ranks.reduce(0, +)) / Double(ranks.count)
    }

    var body: some View {
        VStack(spacing: 0) {
            // App header bar
            HStack(spacing: 14) {
                AppIconView(app: app, size: 36)
                VStack(alignment: .leading, spacing: 2) {
                    Text(app.name)
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    HStack(spacing: 6) {
                        Text(app.platform == "ios" ? "iOS" : "Android")
                            .font(.system(size: 9, weight: .semibold))
                            .padding(.horizontal, 5).padding(.vertical, 2)
                            .background(colorScheme == .dark ? Color.white.opacity(0.08) : Color.primary.opacity(0.08))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                            .clipShape(RoundedRectangle(cornerRadius: 3))
                        Text(app.bundleID)
                            .font(.system(size: 11, design: .monospaced))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                }
                Spacer()
                // Auto-update status
                if isUpdatingRankings {
                    HStack(spacing: 4) {
                        ProgressView().controlSize(.mini)
                        Text("順位更新中…").font(.caption2).foregroundStyle(.secondary)
                    }
                } else if let d = lastUpdated {
                    Text("更新: \(d, formatter: timeFormatter)")
                        .font(.caption2).foregroundStyle(.tertiary)
                }
            }
            .padding(.horizontal, 20).padding(.vertical, 12)
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))

            // Tab bar
            HStack(spacing: 0) {
                ForEach(Tab.allCases, id: \.self) { tab in
                    let isActive = selectedTab == tab
                    Button { selectedTab = tab } label: {
                        VStack(spacing: 0) {
                            HStack(spacing: 6) {
                                Text(tab.rawValue)
                                    .font(.system(size: 13, weight: isActive ? .semibold : .medium))
                                    .foregroundStyle(colorScheme == .dark
                                        ? (isActive ? Aurora.text : Aurora.textMuted)
                                        : (isActive ? Color.primary : Color.secondary))
                                let badgeCount: Int? = {
                                    switch tab {
                                    case .keywords: return keywords.count
                                    case .rising:   return risingKeywords.isEmpty ? nil : risingKeywords.count
                                    default:        return nil
                                    }
                                }()
                                if let n = badgeCount {
                                    Text("\(n)")
                                        .font(.system(size: 10, weight: .semibold)).monospacedDigit()
                                        .padding(.horizontal, 5).padding(.vertical, 1)
                                        .background(RoundedRectangle(cornerRadius: 3).fill(
                                            isActive
                                                ? (colorScheme == .dark ? Aurora.accentBg : Color.accentColor.opacity(0.12))
                                                : (colorScheme == .dark ? Color.white.opacity(0.06) : Color.primary.opacity(0.06))
                                        ))
                                        .foregroundStyle(isActive
                                            ? (colorScheme == .dark ? Aurora.accent : Color.accentColor)
                                            : (colorScheme == .dark ? Aurora.textDim : Color.secondary))
                                }
                            }
                            .padding(.horizontal, 14).padding(.vertical, 11)
                            Rectangle()
                                .fill(isActive ? (colorScheme == .dark ? Aurora.accent : Color.accentColor) : Color.clear)
                                .frame(height: 2)
                        }
                    }
                    .buttonStyle(.plain)
                    .animation(.easeOut(duration: 0.12), value: isActive)
                }
                Spacer()
            }
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))

            Rectangle()
                .fill(colorScheme == .dark ? Aurora.divider : Color.primary.opacity(0.08))
                .frame(height: 1)

            switch selectedTab {
            case .keywords: keywordsPane
            case .rising:   risingPane
            case .gap:      KeywordGapView(app: app).environmentObject(appState)
            case .metadata: MetadataView(app: app).environmentObject(appState)
            }
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
        .navigationTitle(app.name)
        .sheet(isPresented: $showAddKeyword) {
            AddKeywordView(app: app) { newKw in keywords.append(newKw) }
                .environmentObject(appState)
        }
        .sheet(isPresented: $showImportKeywords) {
            ImportKeywordsView(app: app) { Task { await fetchKeywords() } }
                .environmentObject(appState)
        }
        .sheet(isPresented: $showSuggestions) {
            KeywordSuggestionsView(app: app) { newKws in keywords.append(contentsOf: newKws) }
                .environmentObject(appState)
        }
    }

    // MARK: - Keywords pane

    private var keywordsPane: some View {
        HSplitView {
            VStack(spacing: 0) {
                // Header
                HStack(spacing: 6) {
                    Text("キーワード")
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    Spacer()

                    if multiSelectMode {
                        Text("\(selectedKeywords.count)件選択")
                            .font(.caption).foregroundStyle(.secondary)
                        Button {
                            bulkDelete()
                        } label: {
                            Image(systemName: "trash").font(.system(size: 12))
                                .foregroundStyle(selectedKeywords.isEmpty
                                    ? (colorScheme == .dark ? Aurora.textDim : Color.secondary)
                                    : (colorScheme == .dark ? Aurora.negText : .red))
                        }
                        .buttonStyle(.plain).disabled(selectedKeywords.isEmpty)
                        Button("完了") { multiSelectMode = false; selectedKeywords.removeAll() }
                            .buttonStyle(.borderless)
                            .font(.system(size: 12, weight: .medium))
                            .foregroundStyle(Color.accentColor)
                    } else {
                        Button { showSuggestions = true } label: {
                            Image(systemName: "lightbulb").font(.system(size: 12))
                                .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                        }
                        .buttonStyle(.plain).help("キーワード提案")

                        Button { showImportKeywords = true } label: {
                            Image(systemName: "square.and.arrow.down").font(.system(size: 12))
                                .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                        }
                        .buttonStyle(.plain).help("CSV からインポート")

                        Button { showAddKeyword = true } label: {
                            Image(systemName: "plus").font(.system(size: 12))
                                .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                        }
                        .buttonStyle(.plain)

                        Button {
                            multiSelectMode = true
                            selectedKeywords.removeAll()
                        } label: {
                            Image(systemName: "checkmark.circle").font(.system(size: 12))
                                .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                        }
                        .buttonStyle(.plain).help("複数選択して削除")
                    }
                }
                .padding(.horizontal, 12).padding(.vertical, 9)
                .background(colorScheme == .dark ? Aurora.surface : Color.primary.opacity(0.02))

                // Stats bar
                if !keywords.isEmpty && !multiSelectMode {
                    HStack(spacing: 10) {
                        statChip(label: "Top10", value: "\(statsTop10)", color: .green)
                        statChip(label: "Top50", value: "\(statsTop50)", color: .orange)
                        if let avg = statsAvgRank {
                            statChip(label: "平均", value: String(format: "%.0f位", avg), color: .blue)
                        }
                        Spacer()
                    }
                    .padding(.horizontal, 10).padding(.vertical, 6)
                    .background(colorScheme == .dark ? Color.white.opacity(0.03) : Color.primary.opacity(0.02))
                }

                // Search bar
                if !keywords.isEmpty && !multiSelectMode {
                    HStack(spacing: 6) {
                        Image(systemName: "magnifyingglass").font(.system(size: 10)).foregroundStyle(.secondary)
                        TextField("キーワードを検索", text: $searchText)
                            .font(.system(size: 12)).textFieldStyle(.plain)
                        if !searchText.isEmpty {
                            Button { searchText = "" } label: {
                                Image(systemName: "xmark.circle.fill").font(.system(size: 11)).foregroundStyle(.secondary)
                            }.buttonStyle(.plain)
                        }
                    }
                    .padding(.horizontal, 8).padding(.vertical, 5)
                    .background(RoundedRectangle(cornerRadius: 6)
                        .fill(colorScheme == .dark ? Color.white.opacity(0.05) : Color.primary.opacity(0.06)))
                    .padding(.horizontal, 8).padding(.vertical, 4)
                }

                // Sort & Filter bar
                if !keywords.isEmpty && !multiSelectMode {
                    HStack(spacing: 6) {
                        Picker("", selection: $sortOrder) {
                            ForEach(SortOrder.allCases, id: \.self) { Text($0.rawValue).tag($0) }
                        }
                        .labelsHidden().frame(width: 80).controlSize(.small)

                        if availableCountries.count > 2 {
                            Picker("", selection: $filterCountry) {
                                ForEach(availableCountries, id: \.self) { Text($0).tag($0) }
                            }
                            .labelsHidden().frame(width: 70).controlSize(.small)
                        }
                        Spacer()
                    }
                    .padding(.horizontal, 8).padding(.vertical, 4)
                    .background(colorScheme == .dark ? Color.white.opacity(0.02) : Color.primary.opacity(0.02))
                }

                Rectangle()
                    .fill(colorScheme == .dark ? Aurora.divider : Color.primary.opacity(0.08))
                    .frame(height: 1)

                if isLoading {
                    ProgressView().frame(maxWidth: .infinity, maxHeight: .infinity)
                } else if displayedKeywords.isEmpty {
                    ContentUnavailableView(keywords.isEmpty ? "キーワードなし" : "該当なし", systemImage: "magnifyingglass")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else {
                    ScrollView {
                        LazyVStack(spacing: 1) {
                            ForEach(displayedKeywords, id: \.id) { kw in
                                if multiSelectMode {
                                    KeywordRow(keyword: kw, rankSummary: rankSummaries[kw.id], isSelected: selectedKeywords.contains(kw.id), isMultiSelect: true)
                                        .onTapGesture {
                                            if selectedKeywords.contains(kw.id) { selectedKeywords.remove(kw.id) }
                                            else { selectedKeywords.insert(kw.id) }
                                        }
                                } else {
                                    KeywordRow(keyword: kw, rankSummary: rankSummaries[kw.id], isSelected: selectedKeyword?.id == kw.id, isMultiSelect: false)
                                        .onTapGesture { selectedKeyword = kw }
                                }
                            }
                        }
                        .padding(.vertical, 6).padding(.horizontal, 6)
                    }
                }
            }
            .frame(minWidth: 220, idealWidth: 240, maxWidth: 300)
            .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
            .task { await loadKeywordsAndUpdate() }

            if multiSelectMode {
                VStack {
                    Spacer()
                    VStack(spacing: 8) {
                        Image(systemName: "checkmark.circle.fill")
                            .font(.system(size: 36)).foregroundStyle(.secondary)
                        Text("削除するキーワードを選択してください")
                            .font(.subheadline).foregroundStyle(.secondary)
                        if !keywords.isEmpty {
                            Button("すべて選択") {
                                selectedKeywords = Set(keywords.map(\.id))
                            }
                            .buttonStyle(.borderless).foregroundStyle(Color.accentColor)
                        }
                    }
                    Spacer()
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
            } else if let kw = selectedKeyword {
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

    private func loadKeywordsAndUpdate() async {
        await fetchKeywords()
        await fetchRankSummaries()
        await fetchRising()
        await autoUpdateRankings()
    }

    private func fetchKeywords() async {
        isLoading = true
        keywords = (try? await APIClient.shared.getKeywords(token: appState.token, appID: app.id)) ?? []
        isLoading = false
    }

    private func fetchRankSummaries() async {
        if let summaries = try? await APIClient.shared.getKeywordRanks(token: appState.token, appID: app.id) {
            var dict = [String: KeywordRankSummary]()
            for s in summaries { dict[s.keywordID] = s }
            await MainActor.run { rankSummaries = dict }
        }
    }

    private func fetchRising() async {
        if let rising = try? await APIClient.shared.getRisingKeywords(token: appState.token, appID: app.id) {
            await MainActor.run { risingKeywords = rising }
        }
    }

    private func autoUpdateRankings() async {
        guard !keywords.isEmpty else { return }
        isUpdatingRankings = true

        // Snapshot rankings before scrape
        let kwSnapshot = keywords
        var snapshot = [String: Int?]()
        await withTaskGroup(of: (String, Int?).self) { group in
            for kw in kwSnapshot {
                group.addTask {
                    let r = try? await APIClient.shared.getLatestRanking(
                        token: appState.token, appID: app.id, keywordID: kw.id)
                    return (kw.id, r?.rank)
                }
            }
            for await (id, rank) in group { snapshot[id] = rank }
        }

        do {
            try await APIClient.shared.scrapeRankings(token: appState.token, appID: app.id)
            await MainActor.run { lastUpdated = Date() }

            // Compare after scrape and fire notifications
            await withTaskGroup(of: Void.self) { group in
                for kw in kwSnapshot {
                    group.addTask {
                        let newRanking = try? await APIClient.shared.getLatestRanking(
                            token: appState.token, appID: app.id, keywordID: kw.id)
                        let newRank = newRanking?.rank
                        let oldRank = snapshot[kw.id] ?? nil
                        if let old = oldRank, let new = newRank, abs(old - new) >= 3 {
                            AppState.sendRankingNotification(
                                keyword: kw.keyword, appName: app.name, oldRank: old, newRank: new)
                        }
                    }
                }
            }

            keywords = (try? await APIClient.shared.getKeywords(token: appState.token, appID: app.id)) ?? keywords
            await fetchRankSummaries()
            await fetchRising()
        } catch { }

        await MainActor.run { isUpdatingRankings = false }
    }

    @ViewBuilder
    private func statChip(label: String, value: String, color: Color) -> some View {
        HStack(spacing: 3) {
            Text(value).font(.system(size: 11, weight: .bold, design: .monospaced)).foregroundStyle(color)
            Text(label).font(.system(size: 9)).foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
        }
        .padding(.horizontal, 6).padding(.vertical, 3)
        .background(color.opacity(0.10))
        .clipShape(RoundedRectangle(cornerRadius: 5))
    }

    private var risingPane: some View {
        VStack(spacing: 0) {
            HStack {
                Text("急上昇キーワード")
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                Spacer()
                Button { Task { await fetchRising() } } label: {
                    Image(systemName: "arrow.clockwise").font(.system(size: 12))
                        .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                }
                .buttonStyle(.plain)
            }
            .padding(.horizontal, 16).padding(.vertical, 10)
            .background(colorScheme == .dark ? Aurora.surface : Color.primary.opacity(0.02))

            Divider().opacity(0.5)

            if risingKeywords.isEmpty {
                ContentUnavailableView("急上昇なし", systemImage: "arrow.up.right",
                    description: Text("過去7日間で3位以上改善したキーワードがありません"))
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                Table(risingKeywords) {
                    TableColumn("キーワード") { r in
                        Text(r.keyword)
                            .font(.system(size: 13, weight: .medium))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    }
                    TableColumn("国") { r in
                        Text(r.country).font(.system(size: 11)).foregroundStyle(.secondary)
                    }.width(40)
                    TableColumn("現在") { r in
                        Text("\(r.currentRank)位")
                            .font(.system(size: 12, weight: .semibold)).monospacedDigit()
                            .foregroundStyle(Color.asoPosText(colorScheme))
                    }.width(60)
                    TableColumn("前回") { r in
                        Text("\(r.previousRank)位")
                            .font(.system(size: 11)).monospacedDigit().foregroundStyle(.secondary)
                    }.width(60)
                    TableColumn("改善") { r in
                        HStack(spacing: 3) {
                            Image(systemName: "arrow.up").font(.system(size: 9)).foregroundStyle(.green)
                            Text("\(r.improvement)")
                                .font(.system(size: 11, weight: .bold)).foregroundStyle(.green)
                        }
                    }.width(50)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            }
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
    }

    private func bulkDelete() {
        let ids = selectedKeywords
        Task {
            for id in ids {
                try? await APIClient.shared.deleteKeyword(token: appState.token, appID: app.id, keywordID: id)
            }
            await MainActor.run {
                keywords.removeAll { ids.contains($0.id) }
                selectedKeywords.removeAll()
                if let sel = selectedKeyword, ids.contains(sel.id) { selectedKeyword = nil }
                multiSelectMode = false
            }
        }
    }
}

private let timeFormatter: DateFormatter = {
    let f = DateFormatter()
    f.dateFormat = "HH:mm"
    return f
}()

struct KeywordRow: View {
    let keyword: Keyword
    let rankSummary: KeywordRankSummary?
    let isSelected: Bool
    let isMultiSelect: Bool
    @Environment(\.colorScheme) var colorScheme
    @State private var isHovered = false

    var body: some View {
        HStack(spacing: 8) {
            if isMultiSelect {
                Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
                    .foregroundStyle(isSelected ? Color.accentColor : Color.secondary)
                    .font(.system(size: 15))
            }
            VStack(alignment: .leading, spacing: 3) {
                Text(keyword.keyword)
                    .font(.system(size: 13, weight: isSelected && !isMultiSelect ? .semibold : .regular))
                    .foregroundStyle(colorScheme == .dark
                        ? (isSelected && !isMultiSelect ? Aurora.rowSelText : Aurora.text)
                        : .primary)
                HStack(spacing: 6) {
                    Text(keyword.country)
                        .font(.system(size: 11))
                        .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    if let score = keyword.popularityScore { popularityBadge(score) }
                }
            }
            Spacer()
            // Rank + change
            if let summary = rankSummary {
                VStack(alignment: .trailing, spacing: 1) {
                    if let rank = summary.currentRank {
                        Text("\(rank)").font(.system(size: 13, weight: .bold, design: .monospaced))
                            .foregroundStyle(rankColor(rank))
                    } else {
                        Text("圏外").font(.system(size: 9, weight: .semibold)).foregroundStyle(.secondary)
                    }
                    if let change = summary.change, change != 0 {
                        HStack(spacing: 1) {
                            Image(systemName: change > 0 ? "arrow.up" : "arrow.down")
                                .font(.system(size: 7))
                            Text("\(abs(change))").font(.system(size: 9, weight: .semibold))
                        }
                        .foregroundStyle(change > 0 ? Color.green : Color.red)
                    }
                }
            }
        }
        .padding(.horizontal, 10).padding(.vertical, 8)
        .background(
            RoundedRectangle(cornerRadius: 7).fill(
                isSelected
                    ? (colorScheme == .dark
                        ? (isMultiSelect ? Color.accentColor.opacity(0.12) : Aurora.rowSel)
                        : Color.accentColor.opacity(0.10))
                    : (isHovered
                        ? (colorScheme == .dark ? Color.white.opacity(0.04) : Color.primary.opacity(0.04))
                        : .clear)
            )
        )
        .onHover { isHovered = $0 }
        .animation(.easeOut(duration: 0.12), value: isHovered)
        .animation(.easeOut(duration: 0.12), value: isSelected)
    }

    private func rankColor(_ rank: Int) -> Color {
        switch rank {
        case 1...10:  return colorScheme == .dark ? Aurora.posText : Daylight.posText
        case 11...50: return colorScheme == .dark ? Aurora.warnText : Daylight.warnText
        default:      return colorScheme == .dark ? Aurora.textDim : Color.secondary
        }
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
            .padding(.horizontal, 5).padding(.vertical, 2)
            .background(color.opacity(0.16))
            .foregroundStyle(color)
            .clipShape(RoundedRectangle(cornerRadius: 3))
    }
}
