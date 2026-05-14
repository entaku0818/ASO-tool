//  AppDetailView.swift

import SwiftUI

struct AppDetailView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp

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
            // Tab bar
            HStack(spacing: 0) {
                ForEach(Tab.allCases, id: \.self) { tab in
                    Button(tab.rawValue) { selectedTab = tab }
                        .buttonStyle(.plain)
                        .padding(.horizontal, 16)
                        .padding(.vertical, 8)
                        .background(selectedTab == tab ? Color.accentColor.opacity(0.15) : .clear)
                        .foregroundStyle(selectedTab == tab ? Color.accentColor : .primary)
                        .fontWeight(selectedTab == tab ? .semibold : .regular)
                }
                Spacer()
            }
            .background(Color(nsColor: .windowBackgroundColor))

            Divider()

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
            VStack(spacing: 0) {
                HStack {
                    Text("キーワード").font(.headline)
                    Spacer()
                    Button { showAddKeyword = true } label: {
                        Image(systemName: "plus")
                    }
                    .buttonStyle(.borderless)
                    Button { deleteSelected() } label: {
                        Image(systemName: "trash")
                    }
                    .buttonStyle(.borderless)
                    .disabled(selectedKeyword == nil)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)

                Divider()

                if isLoading {
                    ProgressView().frame(maxWidth: .infinity, maxHeight: .infinity)
                } else if keywords.isEmpty {
                    ContentUnavailableView("キーワードなし", systemImage: "magnifyingglass")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else {
                    List(keywords, id: \.id, selection: $selectedKeyword) { kw in
                        KeywordRow(keyword: kw)
                            .tag(kw)
                    }
                }
            }
            .frame(minWidth: 220, idealWidth: 240, maxWidth: 300)
            .task { await fetchKeywords() }

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

    var body: some View {
        VStack(alignment: .leading, spacing: 3) {
            Text(keyword.keyword).font(.body)
            HStack(spacing: 6) {
                Text(keyword.country).font(.caption).foregroundStyle(.secondary)
                if let score = keyword.popularityScore {
                    popularityBadge(score)
                }
            }
        }
        .padding(.vertical, 2)
    }

    @ViewBuilder
    private func popularityBadge(_ score: Int) -> some View {
        let (label, color): (String, Color) = {
            switch score {
            case 4...5: return ("高", .green)
            case 2...3: return ("中", .orange)
            default:    return ("低", .red)
            }
        }()
        Text(label)
            .font(.caption2)
            .padding(.horizontal, 5)
            .padding(.vertical, 2)
            .background(color.opacity(0.2))
            .foregroundStyle(color)
            .clipShape(Capsule())
    }
}
