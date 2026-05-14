//  SettingsView.swift

import SwiftUI

struct SettingsView: View {
    @EnvironmentObject var appState: AppState
    @State private var showDeactivateAlert = false

    var body: some View {
        Form {
            Section("ライセンス") {
                LabeledContent("メールアドレス", value: appState.email)
                LabeledContent("ライセンスキー", value: appState.licenseKey)
                Button("ライセンスを無効化", role: .destructive) {
                    showDeactivateAlert = true
                }
            }

            Section("API") {
                LabeledContent("エンドポイント", value: APIClient.shared.baseURL)
            }
        }
        .formStyle(.grouped)
        .frame(width: 420, height: 260)
        .alert("ライセンスを無効化", isPresented: $showDeactivateAlert) {
            Button("キャンセル", role: .cancel) {}
            Button("無効化", role: .destructive) { appState.deactivate() }
        } message: {
            Text("ライセンスを無効化するとアプリを使用できなくなります。再度ライセンスキーを入力することで再アクティベートできます。")
        }
    }
}
