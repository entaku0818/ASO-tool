//  RankingChartView.swift

import SwiftUI
import Charts

struct RankingChartView: View {
    @EnvironmentObject var appState: AppState
    let app: ASOApp
    let keyword: Keyword
    @Environment(\.colorScheme) var colorScheme

    @State private var rankings: [RankingHistory] = []
    @State private var isLoading = false
    @State private var limit = 30

    private var sorted: [RankingHistory] {
        rankings.sorted { $0.recordedAt < $1.recordedAt }
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack {
                VStack(alignment: .leading, spacing: 2) {
                    Text(keyword.keyword)
                        .font(.title2.bold())
                        .foregroundStyle(colorScheme == .dark ? Aurora.text : .primary)
                    if let latest = sorted.last, let rank = latest.rank {
                        Text("最新: \(rank)位")
                            .font(.caption)
                            .foregroundStyle(colorScheme == .dark ? Aurora.textMuted : .secondary)
                    }
                }
                Spacer()
                Picker("期間", selection: $limit) {
                    Text("7日").tag(7)
                    Text("30日").tag(30)
                    Text("90日").tag(90)
                }
                .pickerStyle(.segmented)
                .tint(Color.asoAccent(colorScheme))
                .frame(width: 180)
            }
            .padding(.horizontal)
            .padding(.top)

            if isLoading {
                ProgressView().frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if sorted.filter({ $0.rank != nil }).isEmpty {
                ContentUnavailableView("データなし", systemImage: "chart.line.downtrend.xyaxis")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                Chart(sorted) { r in
                    if let rank = r.rank {
                        LineMark(
                            x: .value("日付", r.recordedAt),
                            y: .value("順位", rank)
                        )
                        .foregroundStyle(Color.asoAccent(colorScheme))
                        .interpolationMethod(.catmullRom)

                        PointMark(
                            x: .value("日付", r.recordedAt),
                            y: .value("順位", rank)
                        )
                        .foregroundStyle(Color.asoAccent(colorScheme))
                        .symbolSize(30)
                    }
                }
                .chartYScale(domain: .automatic(includesZero: false, reversed: true))
                .chartYAxis {
                    AxisMarks(values: .automatic(desiredCount: 6)) { value in
                        AxisGridLine()
                            .foregroundStyle(colorScheme == .dark ? Aurora.border : Color.primary.opacity(0.08))
                        AxisValueLabel {
                            if let v = value.as(Int.self) {
                                Text("\(v)位")
                                    .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                            }
                        }
                    }
                }
                .chartXAxis {
                    AxisMarks(values: .automatic(desiredCount: 6)) { _ in
                        AxisGridLine()
                            .foregroundStyle(colorScheme == .dark ? Aurora.border : Color.primary.opacity(0.08))
                        AxisValueLabel(format: .dateTime.month().day())
                            .foregroundStyle(colorScheme == .dark ? Aurora.textDim : .secondary)
                    }
                }
                .padding()
            }
        }
        .background(colorScheme == .dark ? Aurora.windowBg : Color(nsColor: .windowBackgroundColor))
        .overlay(colorScheme == .dark ? AuroraGradient() : nil)
        .task(id: "\(keyword.id)-\(limit)") { await fetch() }
        .onChange(of: limit) { _, _ in Task { await fetch() } }
    }

    private func fetch() async {
        isLoading = true
        rankings = (try? await APIClient.shared.getRankings(
            token: appState.token, appID: app.id,
            keywordID: keyword.id, limit: limit
        )) ?? []
        isLoading = false
    }
}
