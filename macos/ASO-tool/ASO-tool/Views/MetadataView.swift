//  MetadataView.swift

import SwiftUI

private let locales = [
    ("ja", "日本語"),
    ("en-US", "English (US)"),
    ("zh-Hans", "中文(简体)"),
    ("ko", "한국어"),
]

private let charLimits: [(String, Int)] = [
    ("title", 30),
    ("subtitle", 30),
    ("keywords", 100),
    ("promotional_text", 170),
    ("description", 4000),
]

struct MetadataView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp

    @State private var versions: [AppMetadataVersion] = []
    @State private var isLoading = false
    @State private var isSaving = false
    @State private var errorMessage: String?

    @State private var selectedLocale = "ja"
    @State private var versionTag = "draft"
    @State private var title = ""
    @State private var subtitle = ""
    @State private var keywords = ""
    @State private var promotionalText = ""
    @State private var descriptionText = ""

    private var currentVersion: AppMetadataVersion? {
        versions.first { $0.locale == selectedLocale && $0.versionTag == versionTag }
    }

    var body: some View {
        HSplitView {
            // Left: version list
            VStack(spacing: 0) {
                HStack {
                    Text("保存済み").font(.headline)
                    Spacer()
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                Divider()
                if versions.isEmpty {
                    Text("なし").foregroundStyle(.secondary).padding()
                } else {
                    List(versions, id: \.id) { v in
                        VStack(alignment: .leading, spacing: 3) {
                            Text(localeName(v.locale)).font(.body)
                            HStack(spacing: 6) {
                                Text(v.versionTag).font(.caption).foregroundStyle(.secondary)
                                if let t = v.title { Text(t).font(.caption).foregroundStyle(.secondary).lineLimit(1) }
                            }
                        }
                        .padding(.vertical, 2)
                        .contentShape(Rectangle())
                        .onTapGesture { load(v) }
                        .contextMenu {
                            Button("削除", role: .destructive) {
                                Task { await deleteVersion(v) }
                            }
                        }
                    }
                }
            }
            .frame(minWidth: 180, idealWidth: 200, maxWidth: 240)

            // Right: editor
            ScrollView {
                VStack(alignment: .leading, spacing: 16) {
                    HStack(spacing: 12) {
                        VStack(alignment: .leading) {
                            Text("言語").font(.caption).foregroundStyle(.secondary)
                            Picker("", selection: $selectedLocale) {
                                ForEach(locales, id: \.0) { code, name in
                                    Text(name).tag(code)
                                }
                            }
                            .labelsHidden()
                            .frame(width: 160)
                        }
                        VStack(alignment: .leading) {
                            Text("バージョンタグ").font(.caption).foregroundStyle(.secondary)
                            TextField("draft, v1.0 ...", text: $versionTag)
                                .textFieldStyle(.roundedBorder)
                                .frame(width: 120)
                        }
                    }

                    labeledField("タイトル（30字以内）", text: $title, limit: 30)
                    labeledField("サブタイトル（30字以内）", text: $subtitle, limit: 30)
                    labeledField("キーワード（100字・カンマ区切り）", text: $keywords, limit: 100)
                    labeledTextArea("プロモーションテキスト（170字以内）", text: $promotionalText, limit: 170, lines: 3)
                    labeledTextArea("説明文（4000字以内）", text: $descriptionText, limit: 4000, lines: 8)

                    if let err = errorMessage {
                        Text(err).foregroundStyle(.red).font(.caption)
                    }

                    HStack {
                        Spacer()
                        if currentVersion != nil {
                            Button("削除", role: .destructive) {
                                Task { if let v = currentVersion { await deleteVersion(v) } }
                            }
                        }
                        Button(isSaving ? "保存中..." : currentVersion != nil ? "更新" : "保存") {
                            Task { await save() }
                        }
                        .buttonStyle(.borderedProminent)
                        .disabled(isSaving)
                    }
                }
                .padding(20)
            }
            .frame(maxWidth: .infinity)
        }
        .task { await fetchVersions() }
        .onChange(of: selectedLocale) { _, _ in applyCurrentVersion() }
        .onChange(of: versionTag) { _, _ in applyCurrentVersion() }
        .onChange(of: versions) { _, _ in applyCurrentVersion() }
    }

    // MARK: - Helpers

    private func localeName(_ code: String) -> String {
        locales.first { $0.0 == code }?.1 ?? code
    }

    private func load(_ v: AppMetadataVersion) {
        selectedLocale = v.locale
        versionTag = v.versionTag
    }

    private func applyCurrentVersion() {
        if let v = currentVersion {
            title = v.title ?? ""
            subtitle = v.subtitle ?? ""
            keywords = v.keywords ?? ""
            promotionalText = v.promotionalText ?? ""
            descriptionText = v.description ?? ""
        } else {
            title = ""; subtitle = ""; keywords = ""; promotionalText = ""; descriptionText = ""
        }
    }

    // MARK: - API

    private func fetchVersions() async {
        isLoading = true
        do {
            versions = try await APIClient.shared.getMetadata(token: appState.token, appID: app.id)
        } catch {
            errorMessage = error.localizedDescription
        }
        isLoading = false
    }

    private func save() async {
        isSaving = true
        errorMessage = nil
        let req = UpsertMetadataRequest(
            locale: selectedLocale,
            versionTag: versionTag,
            title: title.isEmpty ? nil : title,
            subtitle: subtitle.isEmpty ? nil : subtitle,
            description: descriptionText.isEmpty ? nil : descriptionText,
            keywords: keywords.isEmpty ? nil : keywords,
            promotionalText: promotionalText.isEmpty ? nil : promotionalText
        )
        do {
            _ = try await APIClient.shared.upsertMetadata(token: appState.token, appID: app.id, request: req)
            await fetchVersions()
        } catch {
            errorMessage = error.localizedDescription
        }
        isSaving = false
    }

    private func deleteVersion(_ v: AppMetadataVersion) async {
        do {
            try await APIClient.shared.deleteMetadata(token: appState.token, appID: app.id, metadataID: v.id)
            versions.removeAll { $0.id == v.id }
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    // MARK: - Field builders

    @ViewBuilder
    private func labeledField(_ label: String, text: Binding<String>, limit: Int) -> some View {
        VStack(alignment: .leading, spacing: 4) {
            HStack {
                Text(label).font(.caption).foregroundStyle(.secondary)
                Spacer()
                charCounter(text.wrappedValue, limit: limit)
            }
            TextField("", text: text)
                .textFieldStyle(.roundedBorder)
        }
    }

    @ViewBuilder
    private func labeledTextArea(_ label: String, text: Binding<String>, limit: Int, lines: Int) -> some View {
        VStack(alignment: .leading, spacing: 4) {
            HStack {
                Text(label).font(.caption).foregroundStyle(.secondary)
                Spacer()
                charCounter(text.wrappedValue, limit: limit)
            }
            TextEditor(text: text)
                .font(.body)
                .frame(minHeight: CGFloat(lines * 22))
                .overlay(RoundedRectangle(cornerRadius: 6).stroke(Color.secondary.opacity(0.3)))
        }
    }

    @ViewBuilder
    private func charCounter(_ value: String, limit: Int) -> some View {
        let over = value.count > limit
        Text("\(value.count)/\(limit)")
            .font(.caption2)
            .foregroundStyle(over ? .red : .secondary)
            .fontWeight(over ? .bold : .regular)
    }
}
