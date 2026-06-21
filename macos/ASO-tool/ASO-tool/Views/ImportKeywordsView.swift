//  ImportKeywordsView.swift

import SwiftUI
import UniformTypeIdentifiers

struct ImportKeywordsView: View {
    @EnvironmentObject var appState: AppState
    @Environment(\.dismiss) var dismiss

    let app: ASOApp
    var onImported: () -> Void

    @State private var parsed: [(keyword: String, country: String)] = []
    @State private var isImporting = false
    @State private var result: ImportKeywordsResponse?
    @State private var errorMessage: String?
    @State private var showFilePicker = false

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Header
            VStack(alignment: .leading, spacing: 4) {
                Text("キーワードをインポート").font(.title2.bold())
                Text("Astro の CSV エクスポートファイルを選択してください")
                    .font(.subheadline).foregroundStyle(.secondary)
            }
            .padding(20)

            Divider()

            if let result {
                // 完了画面
                VStack(spacing: 16) {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.system(size: 48))
                        .foregroundStyle(.green)
                    Text("インポート完了")
                        .font(.title3.bold())
                    HStack(spacing: 24) {
                        VStack {
                            Text("\(result.imported)").font(.system(size: 28, weight: .bold))
                            Text("追加").font(.caption).foregroundStyle(.secondary)
                        }
                        VStack {
                            Text("\(result.skipped)").font(.system(size: 28, weight: .bold)).foregroundStyle(.secondary)
                            Text("スキップ（重複）").font(.caption).foregroundStyle(.secondary)
                        }
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if parsed.isEmpty {
                // ファイル未選択
                VStack(spacing: 12) {
                    Image(systemName: "doc.text")
                        .font(.system(size: 40))
                        .foregroundStyle(.secondary)
                    Text("CSV ファイルを選択").font(.subheadline).foregroundStyle(.secondary)
                    Button("ファイルを開く…") { showFilePicker = true }
                        .buttonStyle(.borderedProminent)
                    if let err = errorMessage {
                        Text(err).font(.caption).foregroundStyle(.red).multilineTextAlignment(.center)
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                // プレビュー
                VStack(alignment: .leading, spacing: 8) {
                    HStack {
                        Text("\(parsed.count) 件のキーワードが見つかりました")
                            .font(.subheadline.bold())
                        Spacer()
                        Button("別のファイル") { parsed = []; result = nil }
                            .font(.caption)
                            .buttonStyle(.borderless)
                    }
                    .padding(.horizontal, 16)
                    .padding(.top, 12)

                    Table(parsed.prefix(100).map { PreviewRow(keyword: $0.keyword, country: $0.country) }) {
                        TableColumn("キーワード", value: \.keyword)
                        TableColumn("国") { row in
                            Text(row.country).foregroundStyle(.secondary)
                        }
                        .width(60)
                    }
                    .frame(maxHeight: .infinity)

                    if parsed.count > 100 {
                        Text("… 他 \(parsed.count - 100) 件")
                            .font(.caption).foregroundStyle(.secondary)
                            .padding(.horizontal, 16)
                    }

                    if let err = errorMessage {
                        Text(err).font(.caption).foregroundStyle(.red).padding(.horizontal, 16)
                    }
                }
            }

            Divider()

            // Footer
            HStack {
                Button("閉じる") {
                    if result != nil { onImported() }
                    dismiss()
                }
                .keyboardShortcut(.cancelAction)
                Spacer()
                if result == nil && !parsed.isEmpty {
                    Button {
                        Task { await doImport() }
                    } label: {
                        if isImporting {
                            ProgressView().controlSize(.small).frame(width: 80)
                        } else {
                            Text("\(parsed.count) 件をインポート")
                        }
                    }
                    .buttonStyle(.borderedProminent)
                    .disabled(isImporting)
                    .keyboardShortcut(.defaultAction)
                }
            }
            .padding(16)
        }
        .frame(width: 500, height: 440)
        .fileImporter(
            isPresented: $showFilePicker,
            allowedContentTypes: [UTType.commaSeparatedText, UTType.plainText],
            allowsMultipleSelection: false
        ) { res in
            handleFilePick(res)
        }
    }

    private func handleFilePick(_ result: Result<[URL], Error>) {
        errorMessage = nil
        switch result {
        case .failure(let e):
            errorMessage = e.localizedDescription
        case .success(let urls):
            guard let url = urls.first else { return }
            guard url.startAccessingSecurityScopedResource() else {
                errorMessage = "ファイルへのアクセス権がありません"
                return
            }
            defer { url.stopAccessingSecurityScopedResource() }
            do {
                let text = try String(contentsOf: url, encoding: .utf8)
                parsed = parseCSV(text)
                if parsed.isEmpty {
                    errorMessage = "キーワードが見つかりませんでした。ファイル形式を確認してください。"
                }
            } catch {
                errorMessage = "ファイルを読み込めませんでした: \(error.localizedDescription)"
            }
        }
    }

    // Astro CSV: App Name, App Id, Platform, Keyword, Store Domain, ...
    private func parseCSV(_ text: String) -> [(keyword: String, country: String)] {
        var lines = text.components(separatedBy: "\n").map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
        guard let header = lines.first else { return [] }
        lines.removeFirst()

        let cols = header.components(separatedBy: ",").map { $0.lowercased().trimmingCharacters(in: .whitespaces) }
        guard let kwIdx = cols.firstIndex(where: { $0 == "keyword" }) else { return [] }
        let countryIdx = cols.firstIndex(where: { $0 == "store domain" || $0 == "country" })

        var seen = Set<String>()
        var result: [(keyword: String, country: String)] = []

        for line in lines {
            guard !line.isEmpty else { continue }
            let fields = splitCSVLine(line)
            guard fields.count > kwIdx else { continue }
            let kw = fields[kwIdx].trimmingCharacters(in: .whitespaces)
            guard !kw.isEmpty else { continue }
            var country = "JP"
            if let ci = countryIdx, fields.count > ci {
                let raw = fields[ci].trimmingCharacters(in: .whitespaces)
                if !raw.isEmpty { country = raw.uppercased() }
            }
            let key = "\(kw)|\(country)"
            if seen.insert(key).inserted {
                result.append((keyword: kw, country: country))
            }
        }
        return result
    }

    private func splitCSVLine(_ line: String) -> [String] {
        var fields: [String] = []
        var current = ""
        var inQuotes = false
        for ch in line {
            if ch == "\"" { inQuotes.toggle() }
            else if ch == "," && !inQuotes { fields.append(current); current = "" }
            else { current.append(ch) }
        }
        fields.append(current)
        return fields
    }

    private func doImport() async {
        isImporting = true
        errorMessage = nil
        let payload = parsed.map { ["keyword": $0.keyword, "country": $0.country] }
        do {
            let r = try await APIClient.shared.importKeywords(token: appState.token, appID: app.id, keywords: payload)
            await MainActor.run { result = r }
        } catch {
            await MainActor.run { errorMessage = error.localizedDescription }
        }
        await MainActor.run { isImporting = false }
    }
}

private struct PreviewRow: Identifiable {
    let id = UUID()
    let keyword: String
    let country: String
}
