//  KeywordGapView.swift

import SwiftUI

struct KeywordGapView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp

    @State private var gaps: [KeywordGap] = []
    @State private var isLoading = false
    @State private var isUpdating = false
    @State private var errorMessage: String?

    var body: some View {
        VStack(spacing: 0) {
            HStack {
                Text("キーワードギャップ").font(.headline)
                Spacer()
                if isLoading || isUpdating {
                    ProgressView().scaleEffect(0.7)
                }
                Button {
                    Task { await updateRankings() }
                } label: {
                    Label("順位更新", systemImage: "arrow.clockwise")
                        .font(.caption)
                }
                .buttonStyle(.borderless)
                .disabled(isUpdating)
                Button {
                    Task { await fetchGaps() }
                } label: {
                    Label("再取得", systemImage: "arrow.counterclockwise")
                        .font(.caption)
                }
                .buttonStyle(.borderless)
                .disabled(isLoading)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)

            if let err = errorMessage {
                Text(err).foregroundStyle(.red).font(.caption).padding(.horizontal, 12)
            }

            Text("競合が20位以内・自社が30位以下またはランク外のキーワード")
                .font(.caption)
                .foregroundStyle(.secondary)
                .padding(.horizontal, 12)
                .padding(.bottom, 6)

            Divider()

            if isLoading {
                ProgressView().frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if gaps.isEmpty {
                ContentUnavailableView(
                    "ギャップなし",
                    systemImage: "checkmark.shield",
                    description: Text("競合が上位かつ自社が圏外のキーワードはありません")
                )
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                Table(gaps) {
                    TableColumn("キーワード", value: \.keyword)
                    TableColumn("国") { g in
                        Text(g.country).foregroundStyle(.secondary)
                    }
                    TableColumn("競合アプリ", value: \.competitorName)
                    TableColumn("競合順位") { g in
                        Text("\(g.competitorRank)位")
                            .foregroundStyle(.green)
                            .fontWeight(.semibold)
                    }
                    TableColumn("自社順位") { g in
                        if let r = g.ourRank {
                            Text("\(r)位").foregroundStyle(.red)
                        } else {
                            Text("圏外").foregroundStyle(.secondary)
                        }
                    }
                }
            }
        }
        .task { await fetchGaps() }
    }

    private func fetchGaps() async {
        isLoading = true
        errorMessage = nil
        do {
            gaps = try await APIClient.shared.getKeywordGap(token: appState.token, appID: app.id)
        } catch {
            errorMessage = error.localizedDescription
        }
        isLoading = false
    }

    private func updateRankings() async {
        isUpdating = true
        do {
            let result = try await APIClient.shared.updateCompetitorRankings(token: appState.token, appID: app.id)
            // Refresh gaps after update
            if result.updated > 0 {
                await fetchGaps()
            }
        } catch {
            errorMessage = error.localizedDescription
        }
        isUpdating = false
    }
}
