import SwiftUI

struct ASOAdviceView: View {
    let app: ASOApp
    let token: String

    @State private var response: ASOAdviceResponse?
    @State private var isLoading = false
    @State private var error: String?

    private var highAdvice: [ASOAdvice] { response?.advice.filter { $0.priority == "high" } ?? [] }
    private var mediumAdvice: [ASOAdvice] { response?.advice.filter { $0.priority == "medium" } ?? [] }
    private var lowAdvice: [ASOAdvice] { response?.advice.filter { $0.priority == "low" } ?? [] }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Header bar
            HStack {
                VStack(alignment: .leading, spacing: 2) {
                    Text("ASOアドバイス")
                        .font(.headline)
                    if let summary = response?.summary {
                        Text(summary)
                            .font(.caption)
                            .foregroundColor(.secondary)
                            .lineLimit(2)
                    } else {
                        Text("AIがキーワード・競合・メタデータを分析してアドバイスを提供します")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
                Spacer()
                Button(action: generate) {
                    if isLoading {
                        ProgressView().controlSize(.small)
                    } else {
                        Label("アドバイスを生成", systemImage: "sparkles")
                    }
                }
                .disabled(isLoading)
                .buttonStyle(.borderedProminent)
            }
            .padding()

            Divider()

            if let error {
                VStack(spacing: 8) {
                    Image(systemName: "exclamationmark.triangle")
                        .font(.largeTitle)
                        .foregroundColor(.orange)
                    Text(error)
                        .multilineTextAlignment(.center)
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .padding()
            } else if response == nil && !isLoading {
                VStack(spacing: 12) {
                    Image(systemName: "sparkles")
                        .font(.system(size: 48))
                        .foregroundColor(.accentColor.opacity(0.6))
                    Text("アドバイスを生成してください")
                        .font(.title3)
                        .foregroundColor(.secondary)
                    Text("キーワード順位・競合ギャップ・メタデータを基に\nClaudeがASO改善ポイントを提案します")
                        .font(.caption)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .padding()
            } else if isLoading {
                VStack(spacing: 16) {
                    ProgressView()
                    Text("Claudeが分析中...")
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ScrollView {
                    VStack(alignment: .leading, spacing: 20) {
                        if !highAdvice.isEmpty {
                            adviceSection(title: "優先度: 高", items: highAdvice, color: .red)
                        }
                        if !mediumAdvice.isEmpty {
                            adviceSection(title: "優先度: 中", items: mediumAdvice, color: .orange)
                        }
                        if !lowAdvice.isEmpty {
                            adviceSection(title: "優先度: 低", items: lowAdvice, color: .blue)
                        }
                    }
                    .padding()
                }
            }
        }
    }

    @ViewBuilder
    private func adviceSection(title: String, items: [ASOAdvice], color: Color) -> some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack {
                Circle().fill(color).frame(width: 10, height: 10)
                Text(title)
                    .font(.subheadline)
                    .fontWeight(.semibold)
                    .foregroundColor(color)
            }
            ForEach(items) { item in
                AdviceCard(item: item, accentColor: color)
            }
        }
    }

    private func generate() {
        isLoading = true
        error = nil
        Task {
            do {
                let result = try await APIClient.shared.getASOAdvice(token: token, appID: app.id)
                await MainActor.run {
                    response = result
                    isLoading = false
                }
            } catch {
                await MainActor.run {
                    self.error = error.localizedDescription
                    isLoading = false
                }
            }
        }
    }
}

private struct AdviceCard: View {
    let item: ASOAdvice
    let accentColor: Color

    private var categoryIcon: String {
        switch item.category {
        case "keyword":    return "magnifyingglass"
        case "competitor": return "person.2"
        case "metadata":   return "doc.text"
        case "ranking":    return "chart.line.uptrend.xyaxis"
        default:           return "lightbulb"
        }
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack(spacing: 6) {
                Image(systemName: categoryIcon)
                    .foregroundColor(accentColor)
                    .font(.caption)
                Text(item.title)
                    .font(.subheadline)
                    .fontWeight(.medium)
                Spacer()
                Text(item.category)
                    .font(.caption2)
                    .padding(.horizontal, 6)
                    .padding(.vertical, 2)
                    .background(accentColor.opacity(0.1))
                    .foregroundColor(accentColor)
                    .cornerRadius(4)
            }
            Text(item.description)
                .font(.caption)
                .foregroundColor(.secondary)
            HStack(spacing: 4) {
                Image(systemName: "arrow.right.circle.fill")
                    .font(.caption2)
                    .foregroundColor(accentColor)
                Text(item.action)
                    .font(.caption)
                    .fontWeight(.medium)
            }
        }
        .padding(12)
        .background(Color(NSColor.controlBackgroundColor))
        .cornerRadius(8)
        .overlay(
            RoundedRectangle(cornerRadius: 8)
                .stroke(accentColor.opacity(0.2), lineWidth: 1)
        )
    }
}
