//  AddKeywordView.swift

import SwiftUI

struct AddKeywordView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.dismiss) private var dismiss

    let app: ASOApp
    let onAdded: (Keyword) -> Void

    @State private var keywordText = ""
    @State private var country = "JP"
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var suggestions: [KeywordSuggestion] = []
    @State private var isFetchingSuggestions = false

    private let countries = ["JP", "US", "GB", "AU", "CA", "KR", "DE", "FR", "CN"]

    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("キーワード追加").font(.title2.bold())

            VStack(alignment: .leading, spacing: 12) {
                VStack(alignment: .leading, spacing: 4) {
                    Text("キーワード").font(.caption).foregroundStyle(.secondary)
                    TextField("例: アプリ 無料", text: $keywordText)
                        .textFieldStyle(.roundedBorder)
                        .onSubmit { add() }
                }

                VStack(alignment: .leading, spacing: 4) {
                    Text("国").font(.caption).foregroundStyle(.secondary)
                    Picker("国", selection: $country) {
                        ForEach(countries, id: \.self) { Text($0).tag($0) }
                    }
                    .labelsHidden()
                }
            }

            // Suggestions
            VStack(alignment: .leading, spacing: 8) {
                HStack {
                    Text("検索候補").font(.caption).foregroundStyle(.secondary)
                    if isFetchingSuggestions {
                        ProgressView().controlSize(.mini)
                    }
                }
                if suggestions.isEmpty && !isFetchingSuggestions {
                    Text("候補なし（Search Ads 連携が必要です）")
                        .font(.caption)
                        .foregroundStyle(.tertiary)
                } else {
                    ScrollView(.horizontal, showsIndicators: false) {
                        HStack(spacing: 6) {
                            ForEach(suggestions) { s in
                                Button {
                                    keywordText = s.text
                                } label: {
                                    HStack(spacing: 4) {
                                        Text(s.text).font(.caption)
                                        Text("\(s.popularityScore)")
                                            .font(.system(size: 9, weight: .bold))
                                            .foregroundStyle(.white)
                                            .padding(.horizontal, 4)
                                            .padding(.vertical, 2)
                                            .background(scoreColor(s.popularityScore))
                                            .clipShape(Capsule())
                                    }
                                    .padding(.horizontal, 8)
                                    .padding(.vertical, 4)
                                    .background(Color.primary.opacity(0.07))
                                    .clipShape(RoundedRectangle(cornerRadius: 6))
                                }
                                .buttonStyle(.plain)
                            }
                        }
                    }
                    .frame(height: 30)
                }
            }

            if let err = errorMessage {
                Text(err).font(.caption).foregroundStyle(.red)
            }

            HStack(spacing: 12) {
                Button("キャンセル") { dismiss() }.buttonStyle(.bordered)
                    .keyboardShortcut(.cancelAction)
                Spacer()
                Button("追加") { add() }
                    .buttonStyle(.borderedProminent)
                    .disabled(keywordText.trimmingCharacters(in: .whitespaces).isEmpty || isLoading)
                    .keyboardShortcut(.defaultAction)
            }
        }
        .padding(32)
        .frame(width: 420)
        .task { await fetchSuggestions() }
    }

    private func scoreColor(_ score: Int) -> Color {
        switch score {
        case 4...5: return .green
        case 2...3: return .orange
        default:    return .secondary
        }
    }

    private func fetchSuggestions() async {
        isFetchingSuggestions = true
        if let result = try? await APIClient.shared.getKeywordSuggestions(token: appState.token, appID: app.id) {
            await MainActor.run { suggestions = result }
        }
        await MainActor.run { isFetchingSuggestions = false }
    }

    private func add() {
        isLoading = true
        errorMessage = nil
        Task {
            do {
                let kw = try await APIClient.shared.createKeyword(
                    token: appState.token,
                    appID: app.id,
                    keyword: keywordText.trimmingCharacters(in: .whitespaces),
                    country: country
                )
                await MainActor.run { onAdded(kw); dismiss() }
            } catch {
                await MainActor.run { errorMessage = error.localizedDescription; isLoading = false }
            }
        }
    }
}
