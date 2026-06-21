//  KeywordSuggestionsView.swift

import SwiftUI

struct KeywordSuggestionsView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.dismiss) var dismiss
    @Environment(\.colorScheme) var colorScheme

    let app: ASOApp
    var onAdded: ([Keyword]) -> Void

    @State private var seedText = ""
    @State private var country = "JP"
    @State private var suggestions: [AutocompleteSuggestion] = []
    @State private var selected: Set<String> = []
    @State private var isLoading = false
    @State private var isAdding = false
    @State private var errorMessage: String?
    @State private var addedCount: Int?

    private let countries = ["JP", "US", "GB", "AU", "CA", "KR", "DE", "FR", "CN"]

    var body: some View {
        VStack(spacing: 0) {
            // Header
            VStack(alignment: .leading, spacing: 4) {
                Text("キーワード提案").font(.title2.bold())
                Text("App Store のオートコンプリートから関連キーワードを取得します")
                    .font(.subheadline).foregroundStyle(.secondary)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(20)

            Divider()

            // Search controls
            VStack(alignment: .leading, spacing: 10) {
                HStack(spacing: 8) {
                    TextField("キーワードを入力（例: 録音、家計簿）", text: $seedText)
                        .textFieldStyle(.roundedBorder)
                        .onSubmit { Task { await fetch() } }
                    Picker("", selection: $country) {
                        ForEach(countries, id: \.self) { Text($0).tag($0) }
                    }
                    .labelsHidden()
                    .frame(width: 70)
                    Button("検索") { Task { await fetch() } }
                        .buttonStyle(.borderedProminent)
                        .disabled(seedText.trimmingCharacters(in: .whitespaces).isEmpty || isLoading)
                }
                if let err = errorMessage {
                    Text(err).font(.caption).foregroundStyle(.red)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)

            Divider()

            // Results
            if isLoading {
                ProgressView("取得中…").frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if suggestions.isEmpty {
                ContentUnavailableView(
                    "提案なし",
                    systemImage: "lightbulb",
                    description: Text("キーワードを入力して「検索」を押してください")
                )
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                VStack(spacing: 0) {
                    // Select-all bar
                    HStack {
                        Button(selected.count == suggestions.count ? "すべて解除" : "すべて選択") {
                            if selected.count == suggestions.count {
                                selected.removeAll()
                            } else {
                                selected = Set(suggestions.map(\.term))
                            }
                        }
                        .buttonStyle(.borderless)
                        .font(.caption)
                        Text("·").foregroundStyle(.tertiary)
                        Text("\(selected.count) 件選択中").font(.caption).foregroundStyle(.secondary)
                        Spacer()
                        Text("\(suggestions.count) 件").font(.caption).foregroundStyle(.tertiary)
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 8)
                    .background(colorScheme == .dark ? Color.white.opacity(0.03) : Color.primary.opacity(0.03))

                    Divider()

                    ScrollView {
                        LazyVStack(spacing: 1) {
                            ForEach(suggestions) { s in
                                Button {
                                    if selected.contains(s.term) { selected.remove(s.term) }
                                    else { selected.insert(s.term) }
                                } label: {
                                    HStack(spacing: 10) {
                                        Image(systemName: selected.contains(s.term) ? "checkmark.circle.fill" : "circle")
                                            .foregroundStyle(selected.contains(s.term) ? Color.accentColor : Color.secondary)
                                            .font(.system(size: 16))
                                        Text(s.term)
                                            .font(.system(size: 13))
                                            .foregroundStyle(colorScheme == .dark ? .white : .primary)
                                        Spacer()
                                    }
                                    .padding(.horizontal, 14)
                                    .padding(.vertical, 9)
                                    .background(
                                        selected.contains(s.term)
                                            ? (colorScheme == .dark ? Color.accentColor.opacity(0.15) : Color.accentColor.opacity(0.08))
                                            : Color.clear
                                    )
                                }
                                .buttonStyle(.plain)
                                .animation(.easeOut(duration: 0.1), value: selected.contains(s.term))
                            }
                        }
                        .padding(.vertical, 4)
                    }
                }
            }

            Divider()

            // Footer
            HStack {
                if let count = addedCount {
                    Text("\(count) 件を追加しました").font(.caption).foregroundStyle(.green)
                }
                Button("閉じる") { dismiss() }.keyboardShortcut(.cancelAction)
                Spacer()
                Button {
                    Task { await addSelected() }
                } label: {
                    if isAdding {
                        ProgressView().controlSize(.small).frame(width: 100)
                    } else {
                        Text("\(selected.count) 件を追加")
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(selected.isEmpty || isAdding)
                .keyboardShortcut(.defaultAction)
            }
            .padding(16)
        }
        .frame(width: 480, height: 520)
    }

    private func fetch() async {
        let term = seedText.trimmingCharacters(in: .whitespaces)
        guard !term.isEmpty else { return }
        isLoading = true
        errorMessage = nil
        suggestions = []
        selected = []
        do {
            suggestions = try await APIClient.shared.getKeywordAutocomplete(
                token: appState.token, term: term, country: country)
            if suggestions.isEmpty {
                errorMessage = "提案が見つかりませんでした"
            }
        } catch {
            errorMessage = error.localizedDescription
        }
        isLoading = false
    }

    private func addSelected() async {
        isAdding = true
        var added: [Keyword] = []
        for term in selected {
            if let kw = try? await APIClient.shared.createKeyword(
                token: appState.token, appID: app.id, keyword: term, country: country) {
                added.append(kw)
            }
        }
        await MainActor.run {
            addedCount = added.count
            selected.removeAll()
            isAdding = false
            if !added.isEmpty { onAdded(added) }
        }
    }
}
