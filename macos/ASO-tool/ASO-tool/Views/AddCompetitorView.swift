//  AddCompetitorView.swift

import SwiftUI

struct AddCompetitorView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.dismiss) var dismiss

    let app: ASOApp
    var onAdded: (Competitor) -> Void

    @State private var searchQuery = ""
    @State private var country = "jp"
    @State private var results: [AppStoreSearchResult] = []
    @State private var isSearching = false
    @State private var isAdding = false
    @State private var errorMessage: String?

    private let countries = ["jp", "us", "gb", "au", "ca", "kr", "de", "fr"]

    var body: some View {
        VStack(spacing: 0) {
            // Header
            VStack(alignment: .leading, spacing: 4) {
                Text("競合アプリを追加").font(.title2.bold())
                Text("アプリ名や開発者名で検索してください")
                    .font(.subheadline).foregroundStyle(.secondary)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(20)

            Divider()

            // Search bar
            HStack(spacing: 8) {
                TextField("例: 録音、ボイスメモ", text: $searchQuery)
                    .textFieldStyle(.roundedBorder)
                    .onSubmit { Task { await search() } }
                Picker("", selection: $country) {
                    ForEach(countries, id: \.self) { Text($0.uppercased()).tag($0) }
                }
                .labelsHidden().frame(width: 70)
                Button("検索") { Task { await search() } }
                    .buttonStyle(.borderedProminent)
                    .disabled(searchQuery.trimmingCharacters(in: .whitespaces).isEmpty || isSearching)
            }
            .padding(.horizontal, 16).padding(.vertical, 12)

            if let err = errorMessage {
                Text(err).font(.caption).foregroundStyle(.red)
                    .padding(.horizontal, 16).padding(.bottom, 8)
            }

            Divider()

            // Results
            if isSearching {
                ProgressView("検索中…").frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if results.isEmpty {
                ContentUnavailableView("検索結果なし", systemImage: "magnifyingglass")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ScrollView {
                    LazyVStack(spacing: 1) {
                        ForEach(results) { result in
                            CompetitorSearchRow(
                                result: result,
                                isAdding: isAdding,
                                platform: app.platform
                            ) {
                                Task { await add(result) }
                            }
                        }
                    }
                    .padding(.vertical, 6)
                }
            }

            Divider()

            HStack {
                Button("閉じる") { dismiss() }.keyboardShortcut(.cancelAction)
                Spacer()
            }
            .padding(16)
        }
        .frame(width: 520, height: 500)
    }

    private func search() async {
        let q = searchQuery.trimmingCharacters(in: .whitespaces)
        guard !q.isEmpty else { return }
        isSearching = true
        errorMessage = nil
        results = []
        do {
            results = try await APIClient.shared.searchApps(
                token: appState.token, keyword: q, platform: app.platform, country: country)
        } catch {
            errorMessage = error.localizedDescription
        }
        isSearching = false
    }

    private func add(_ result: AppStoreSearchResult) async {
        isAdding = true
        do {
            let competitor = try await APIClient.shared.createCompetitor(
                token: appState.token,
                appID: app.id,
                bundleID: result.appInfo.bundleID,
                name: result.appInfo.name,
                platform: app.platform
            )
            await MainActor.run { onAdded(competitor); dismiss() }
        } catch {
            await MainActor.run { errorMessage = error.localizedDescription }
        }
        await MainActor.run { isAdding = false }
    }
}

private struct CompetitorSearchRow: View {
    let result: AppStoreSearchResult
    let isAdding: Bool
    let platform: String
    let onAdd: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            Text("\(result.rank)")
                .font(.system(size: 11, weight: .semibold, design: .monospaced))
                .foregroundStyle(.secondary)
                .frame(width: 24, alignment: .trailing)

            AsyncImage(url: URL(string: result.appInfo.iconURL)) { img in
                img.resizable().scaledToFill()
            } placeholder: {
                RoundedRectangle(cornerRadius: 8).fill(Color.secondary.opacity(0.2))
            }
            .frame(width: 40, height: 40)
            .clipShape(RoundedRectangle(cornerRadius: 8))

            VStack(alignment: .leading, spacing: 2) {
                Text(result.appInfo.name)
                    .font(.system(size: 13, weight: .medium))
                    .lineLimit(1)
                Text(result.appInfo.developer)
                    .font(.system(size: 11))
                    .foregroundStyle(.secondary)
                    .lineLimit(1)
            }

            Spacer()

            if result.appInfo.rating > 0 {
                HStack(spacing: 2) {
                    Image(systemName: "star.fill")
                        .font(.system(size: 9)).foregroundStyle(.yellow)
                    Text(String(format: "%.1f", result.appInfo.rating))
                        .font(.system(size: 11)).foregroundStyle(.secondary)
                }
            }

            Button("追加") { onAdd() }
                .buttonStyle(.borderedProminent)
                .controlSize(.small)
                .disabled(isAdding)
        }
        .padding(.horizontal, 16).padding(.vertical, 8)
    }
}
