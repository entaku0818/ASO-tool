//  AddAppView.swift

import SwiftUI

struct AddAppView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.dismiss) var dismiss

    var onCreated: (ASOApp) -> Void

    @State private var query = ""
    @State private var platform = "ios"
    @State private var country = "JP"
    @State private var results: [AppStoreSearchResult] = []
    @State private var isSearching = false
    @State private var isAdding: String? = nil  // bundleID being added
    @State private var errorMessage: String?
    @State private var searched = false

    private let countries = ["JP", "US", "GB", "AU", "CA", "KR", "DE", "FR"]

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Header
            VStack(alignment: .leading, spacing: 12) {
                Text("アプリを追加").font(.title2.bold())

                HStack(spacing: 8) {
                    // Search field
                    HStack(spacing: 6) {
                        Image(systemName: "magnifyingglass")
                            .foregroundStyle(.secondary)
                            .font(.system(size: 13))
                        TextField("アプリ名で検索...", text: $query)
                            .textFieldStyle(.plain)
                            .onSubmit { search() }
                    }
                    .padding(.horizontal, 10)
                    .padding(.vertical, 7)
                    .background(Color.primary.opacity(0.06))
                    .clipShape(RoundedRectangle(cornerRadius: 8))

                    Picker("", selection: $platform) {
                        Text("iOS").tag("ios")
                        Text("Android").tag("android")
                    }
                    .pickerStyle(.segmented)
                    .frame(width: 110)
                    .labelsHidden()

                    Picker("", selection: $country) {
                        ForEach(countries, id: \.self) { Text($0).tag($0) }
                    }
                    .frame(width: 70)
                    .labelsHidden()

                    Button(action: search) {
                        if isSearching {
                            ProgressView().controlSize(.small).frame(width: 28, height: 28)
                        } else {
                            Text("検索").frame(width: 44)
                        }
                    }
                    .buttonStyle(.borderedProminent)
                    .disabled(query.trimmingCharacters(in: .whitespaces).isEmpty || isSearching)
                    .keyboardShortcut(.defaultAction)
                }

                if let err = errorMessage {
                    Text(err).font(.caption).foregroundStyle(.red)
                }
            }
            .padding(20)

            Divider()

            // Results
            if results.isEmpty {
                Spacer()
                Group {
                    if !searched {
                        Label("アプリ名を入力して検索してください", systemImage: "magnifyingglass")
                    } else {
                        Label("見つかりませんでした", systemImage: "questionmark.circle")
                    }
                }
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .frame(maxWidth: .infinity)
                Spacer()
            } else {
                ScrollView {
                    LazyVStack(spacing: 0) {
                        ForEach(results) { result in
                            AppSearchRow(
                                result: result,
                                isAdding: isAdding == result.appInfo.bundleID
                            ) {
                                addApp(result)
                            }
                            Divider().padding(.leading, 68)
                        }
                    }
                }
            }

            Divider()
            HStack {
                Spacer()
                Button("閉じる") { dismiss() }
                    .keyboardShortcut(.cancelAction)
            }
            .padding(16)
        }
        .frame(width: 520, height: 480)
    }

    private func search() {
        let kw = query.trimmingCharacters(in: .whitespaces)
        guard !kw.isEmpty else { return }
        isSearching = true
        errorMessage = nil
        results = []
        Task {
            do {
                let r = try await APIClient.shared.searchApps(
                    token: appState.token, keyword: kw, platform: platform, country: country)
                await MainActor.run { results = r; searched = true }
            } catch {
                await MainActor.run { errorMessage = error.localizedDescription; searched = true }
            }
            await MainActor.run { isSearching = false }
        }
    }

    private func addApp(_ result: AppStoreSearchResult) {
        isAdding = result.appInfo.bundleID
        errorMessage = nil
        Task {
            do {
                let app = try await APIClient.shared.createApp(
                    token: appState.token,
                    name: result.appInfo.name,
                    bundleID: result.appInfo.bundleID,
                    platform: platform,
                    storeURL: result.appInfo.storeURL
                )
                await MainActor.run { onCreated(app); dismiss() }
            } catch {
                await MainActor.run { errorMessage = error.localizedDescription; isAdding = nil }
            }
        }
    }
}

private struct AppSearchRow: View {
    let result: AppStoreSearchResult
    let isAdding: Bool
    let onAdd: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            // Rank
            Text("\(result.rank)")
                .font(.system(size: 11, weight: .medium))
                .foregroundStyle(.tertiary)
                .frame(width: 20, alignment: .trailing)

            // Icon
            AsyncImage(url: URL(string: result.appInfo.iconURL)) { image in
                image.resizable().scaledToFill()
            } placeholder: {
                RoundedRectangle(cornerRadius: 10)
                    .fill(Color.primary.opacity(0.08))
            }
            .frame(width: 44, height: 44)
            .clipShape(RoundedRectangle(cornerRadius: 10))

            // Info
            VStack(alignment: .leading, spacing: 3) {
                Text(result.appInfo.name)
                    .font(.system(size: 13, weight: .semibold))
                    .lineLimit(1)
                Text(result.appInfo.developer)
                    .font(.system(size: 11))
                    .foregroundStyle(.secondary)
                    .lineLimit(1)
                HStack(spacing: 4) {
                    Image(systemName: "star.fill")
                        .font(.system(size: 9))
                        .foregroundStyle(.yellow)
                    Text(String(format: "%.1f", result.appInfo.rating))
                        .font(.system(size: 10))
                        .foregroundStyle(.secondary)
                    Text("(\(result.appInfo.ratingCount.formatted()))")
                        .font(.system(size: 10))
                        .foregroundStyle(.tertiary)
                }
            }

            Spacer()

            // Add button
            Button(action: onAdd) {
                if isAdding {
                    ProgressView().controlSize(.small).frame(width: 44)
                } else {
                    Text("追加")
                }
            }
            .buttonStyle(.bordered)
            .disabled(isAdding)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 10)
    }
}
