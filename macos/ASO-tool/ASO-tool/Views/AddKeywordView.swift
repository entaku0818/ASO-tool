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

    private let countries = ["JP", "US", "GB", "AU", "CA", "KR", "DE", "FR", "CN"]

    var body: some View {
        VStack(spacing: 24) {
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
            .frame(width: 300)

            if let err = errorMessage {
                Text(err).font(.caption).foregroundStyle(.red).frame(maxWidth: 300)
            }

            HStack(spacing: 12) {
                Button("キャンセル") { dismiss() }.buttonStyle(.bordered)
                Button("追加") { add() }
                    .buttonStyle(.borderedProminent)
                    .disabled(keywordText.trimmingCharacters(in: .whitespaces).isEmpty || isLoading)
            }
        }
        .padding(32)
        .frame(width: 380)
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
