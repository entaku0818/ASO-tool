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

// MARK: - 旧 Bundle ID からの移行

extension CredentialStore {
    /// ASO-tool 時代の Bundle ID。
    static let legacyBundleID = "com.entaku.ASO-tool"

    /// サンドボックス下では NSHomeDirectory() がコンテナを指すため、実ホームを取得する。
    private static var realHomeDirectory: URL {
        if let pw = getpwuid(getuid()), let dir = pw.pointee.pw_dir {
            return URL(fileURLWithPath: String(cString: dir))
        }
        return URL(fileURLWithPath: NSHomeDirectory())
    }

    private static var legacyPreferencesURL: URL {
        realHomeDirectory.appendingPathComponent(
            "Library/Containers/\(legacyBundleID)/Data/Library/Preferences/\(legacyBundleID).plist")
    }

    /// ASO-tool から ASO Compass への改名で Bundle ID が変わり、サンドボックスの
    /// コンテナが別になったため、旧コンテナの plist を直接読んで引き継ぐ。
    ///
    /// - 移行期間限定の処理。旧バージョンが十分に入れ替わったら、この extension と
    ///   entitlements の temporary-exception.files.home-relative-path.read-only を
    ///   まとめて削除すること。
    /// - テスト用の差し替えストアには作用させない（実ファイルを読ませないため）。
    static func migrateLegacyCredentialsIfNeeded(into store: CredentialStore = .standard) {
        guard store.defaults == .standard else { return }
        guard store.token.isEmpty, store.licenseKey.isEmpty else { return }
        guard let dict = NSDictionary(contentsOf: legacyPreferencesURL) as? [String: Any] else { return }

        let token = dict[tokenKey] as? String ?? ""
        let email = dict[emailKey] as? String ?? ""
        let licenseKey = dict[licenseKeyKey] as? String ?? ""

        // トークンが生きているか、再アクティベートに必要な組が揃っている場合のみ引き継ぐ。
        guard !token.isEmpty || (!licenseKey.isEmpty && !email.isEmpty) else { return }
        store.save(token: token, email: email, licenseKey: licenseKey)
    }
}

extension Notification.Name {
    /// 自動再アクティベートでトークンが更新されたとき。userInfo["token"] に新トークン。
    static let asoTokenRefreshed = Notification.Name("asoTokenRefreshed")
    /// 自動再アクティベートに失敗し、ユーザーの手動再入力が必要なとき。
    static let asoReactivationRequired = Notification.Name("asoReactivationRequired")
}
