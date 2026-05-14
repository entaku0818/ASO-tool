//  AppState.swift

import Foundation

final class AppState: ObservableObject {
    @Published var isActivated: Bool = false
    @Published var token: String = ""
    @Published var email: String = ""
    @Published var licenseKey: String = ""

    private let defaults = UserDefaults.standard

    init() {
        token = defaults.string(forKey: "jwt_token") ?? ""
        email = defaults.string(forKey: "user_email") ?? ""
        licenseKey = defaults.string(forKey: "license_key") ?? ""
        isActivated = !token.isEmpty
    }

    func activate(token: String, email: String, key: String) {
        self.token = token
        self.email = email
        self.licenseKey = key
        self.isActivated = true
        defaults.set(token, forKey: "jwt_token")
        defaults.set(email, forKey: "user_email")
        defaults.set(key, forKey: "license_key")
    }

    func deactivate() {
        token = ""
        email = ""
        licenseKey = ""
        isActivated = false
        defaults.removeObject(forKey: "jwt_token")
        defaults.removeObject(forKey: "user_email")
        defaults.removeObject(forKey: "license_key")
    }
}
