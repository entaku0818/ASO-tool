//  CredentialStore.swift

import Foundation

/// UserDefaults に保存する認証情報のアクセス経路を一元化する。
/// キー名は既存インストールとの互換のため変更しないこと。
struct CredentialStore: @unchecked Sendable {
    static let tokenKey = "jwt_token"
    static let emailKey = "user_email"
    static let licenseKeyKey = "license_key"

    /// アプリ本体が使う共有ストア。テストは別の UserDefaults を差した個別インスタンスを使う。
    static let standard = CredentialStore(defaults: .standard)

    let defaults: UserDefaults

    var token: String { defaults.string(forKey: Self.tokenKey) ?? "" }
    var email: String { defaults.string(forKey: Self.emailKey) ?? "" }
    var licenseKey: String { defaults.string(forKey: Self.licenseKeyKey) ?? "" }

    /// ライセンスキーとメールが揃っていれば、トークン失効時に自動で再アクティベートできる。
    var canReactivate: Bool { !licenseKey.isEmpty && !email.isEmpty }

    func save(token: String, email: String, licenseKey: String) {
        defaults.set(token, forKey: Self.tokenKey)
        defaults.set(email, forKey: Self.emailKey)
        defaults.set(licenseKey, forKey: Self.licenseKeyKey)
    }

    func update(token: String) {
        defaults.set(token, forKey: Self.tokenKey)
    }

    func clear() {
        defaults.removeObject(forKey: Self.tokenKey)
        defaults.removeObject(forKey: Self.emailKey)
        defaults.removeObject(forKey: Self.licenseKeyKey)
    }
}

extension Notification.Name {
    /// 自動再アクティベートでトークンが更新されたとき。userInfo["token"] に新トークン。
    static let asoTokenRefreshed = Notification.Name("asoTokenRefreshed")
    /// 自動再アクティベートに失敗し、ユーザーの手動再入力が必要なとき。
    static let asoReactivationRequired = Notification.Name("asoReactivationRequired")
}
