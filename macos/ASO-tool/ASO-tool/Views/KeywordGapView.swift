//  KeywordGapView.swift

import SwiftUI

struct KeywordGapView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp
    @Environment(\.colorScheme) var colorScheme

    @State private var gaps: [KeywordGap] = []
    @State private var isLoading = false
    @State private var isUpdating = false
    @State private var errorMessage: String?

    var body: some View {
        VStack(spacing: 0) {
            // Header
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    VStack(alignment: .leading, spacing: 2) {
                        Text("キーワードギャップ")
                            .font(.system(size: 11, weight: .semibold))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .tertiary)
                            .tracking(0.4)
                            .textCase(.uppercase)
                        Text("競合が上位、自社が圏外")
                            .font(.system(size: 18, weight: .semibold))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    }
                    Spacer()
                    if isLoading || isUpdating {
                        ProgressView().scaleEffect(0.7)
                    }
                    Button {
                        Task { await updateRankings() }
                    } label: {
                        Label("順位更新", systemImage: "arrow.clockwise")
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                    .buttonStyle(.borderless)
                    .disabled(isUpdating)

                    Button {
                        Task { await fetchGaps() }
                    } label: {
                        Label("再取得", systemImage: "arrow.counterclockwise")
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                    .buttonStyle(.borderless)
                    .disabled(isLoading)
                }

                Text("競合が20位以内・自社が30位以下またはランク外のキーワード")
                    .font(.system(size: 12))
                    .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)

                if let err = errorMessage {
                    Text(err)
                        .foregroundStyle(colorScheme == .dark ? Aurora.negText : .red)
                        .font(.caption)
                }
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 14)

            Rectangle()
                .fill(colorScheme == .dark ? Aurora.divider : Color.primary.opacity(0.08))
                .frame(height: 1)

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
                    TableColumn("キーワード") { g in
                        Text(g.keyword)
                            .font(.system(size: 13, weight: .medium))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    }
                    TableColumn("国") { g in
                        Text(g.country)
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                    TableColumn("競合アプリ") { g in
                        Text(g.competitorName)
                            .font(.system(size: 12))
                            .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    }
                    TableColumn("競合順位") { g in
                        Text("\(g.competitorRank)位")
                            .font(.system(size: 12, weight: .semibold))
                            .monospacedDigit()
                            .foregroundStyle(Color.asoPosText(colorScheme))
                    }
                    TableColumn("自社順位") { g in
                        if let r = g.ourRank {
                            Text("\(r)位")
                                .font(.system(size: 12, weight: .semibold))
                                .monospacedDigit()
                                .foregroundStyle(Color.asoNegText(colorScheme))
                        } else {
                            Text("圏外")
                                .font(.system(size: 10, weight: .bold))
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(
                                    RoundedRectangle(cornerRadius: 4)
                                        .fill(colorScheme == .dark ? Color(hex: "FF8A8A").opacity(0.14) : Color(hex: "C24545").opacity(0.10))
                                )
                                .foregroundStyle(Color.asoNegText(colorScheme))
                        }
                    }
                }
            }
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
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
