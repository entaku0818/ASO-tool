//  LicenseActivationView.swift

import SwiftUI

struct LicenseActivationView: View {
    @EnvironmentObject var appState: AppState
    @State private var licenseKey = ""
    @State private var email = ""
    @State private var isLoading = false
    @State private var errorMessage: String?

    var body: some View {
        VStack(spacing: 32) {
            VStack(spacing: 8) {
                Image(systemName: "chart.bar.xaxis")
                    .font(.system(size: 56))
                    .foregroundStyle(Color.accentColor)
                Text("ASO Tool")
                    .font(.largeTitle.bold())
                Text("ライセンスキーを入力してアクティベートしてください")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            VStack(alignment: .leading, spacing: 12) {
                VStack(alignment: .leading, spacing: 4) {
                    Text("メールアドレス").font(.caption).foregroundStyle(.secondary)
                    TextField("email@example.com", text: $email)
                        .textFieldStyle(.roundedBorder)
                }

                VStack(alignment: .leading, spacing: 4) {
                    Text("ライセンスキー").font(.caption).foregroundStyle(.secondary)
                    TextField("ASOT-XXXX-XXXX-XXXX", text: $licenseKey)
                        .textFieldStyle(.roundedBorder)
                        .font(.system(.body, design: .monospaced))
                }
            }
            .frame(width: 340)

            if let err = errorMessage {
                Text(err)
                    .font(.caption)
                    .foregroundStyle(.red)
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: 340)
            }

            VStack(spacing: 10) {
                Button(action: activate) {
                    if isLoading {
                        ProgressView().controlSize(.small)
                            .frame(width: 340, height: 28)
                    } else {
                        Text("アクティベート")
                            .frame(width: 340, height: 28)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(licenseKey.isEmpty || email.isEmpty || isLoading)

                Button("ライセンスを購入する →") {
                    NSWorkspace.shared.open(URL(string: "https://aso-tool.vercel.app/#pricing")!)
                }
                .buttonStyle(.borderless)
                .font(.footnote)
            }
        }
        .padding(56)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    private func activate() {
        isLoading = true
        errorMessage = nil
        Task {
            do {
                let resp = try await APIClient.shared.activateLicense(
                    key: licenseKey.trimmingCharacters(in: .whitespaces),
                    email: email.trimmingCharacters(in: .whitespaces)
                )
                await MainActor.run {
                    appState.activate(token: resp.token, email: resp.user.email, key: resp.key)
                }
            } catch {
                await MainActor.run {
                    let detail: String
                    switch error {
                    case APIClientError.networkError(let e):
                        detail = "ネットワーク: \(e)"
                    case APIClientError.httpError(let code, let msg):
                        detail = "HTTP\(code): \(msg)"
                    case APIClientError.decodingError(let e):
                        detail = "デコード: \(e)"
                    default:
                        detail = error.localizedDescription
                    }
                    errorMessage = detail
                    isLoading = false
                }
            }
        }
    }
}
