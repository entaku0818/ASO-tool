//  AppState.swift

import Combine
import Foundation
import UserNotifications

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
        requestNotificationPermission()
    }

    private func requestNotificationPermission() {
        UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .sound]) { _, _ in }
    }

    static func sendRankingNotification(keyword: String, appName: String, oldRank: Int?, newRank: Int?) {
        let content = UNMutableNotificationContent()
        content.title = "順位変動: \(appName)"
        let oldStr = oldRank.map { "\($0)位" } ?? "圏外"
        let newStr = newRank.map { "\($0)位" } ?? "圏外"
        let arrow = (newRank ?? 999) < (oldRank ?? 999) ? "↑" : "↓"
        content.body = "\(keyword): \(oldStr) → \(newStr) \(arrow)"
        content.sound = .default
        let req = UNNotificationRequest(identifier: UUID().uuidString, content: content, trigger: nil)
        UNUserNotificationCenter.current().add(req, withCompletionHandler: nil)
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
