//  AppState.swift

import Combine
import Foundation
import UserNotifications

final class AppState: ObservableObject {
    @Published var isActivated: Bool = false
    @Published var token: String = ""
    @Published var email: String = ""
    @Published var licenseKey: String = ""

    private let credentials: CredentialStore
    private var observers: [NSObjectProtocol] = []

    init(credentials: CredentialStore = .standard) {
        self.credentials = credentials
        token = credentials.token
        email = credentials.email
        licenseKey = credentials.licenseKey
        isActivated = !token.isEmpty
        requestNotificationPermission()
        observeTokenRefresh()
    }

    deinit {
        observers.forEach { NotificationCenter.default.removeObserver($0) }
    }

    /// APIClient がトークン失効を検知して自動再アクティベートした結果を反映する。
    private func observeTokenRefresh() {
        let center = NotificationCenter.default
        observers.append(center.addObserver(
            forName: .asoTokenRefreshed, object: nil, queue: .main
        ) { [weak self] note in
            guard let newToken = note.userInfo?["token"] as? String else { return }
            self?.token = newToken
            self?.isActivated = true
        })

        observers.append(center.addObserver(
            forName: .asoReactivationRequired, object: nil, queue: .main
        ) { [weak self] _ in
            self?.deactivate()
        })
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
        credentials.save(token: token, email: email, licenseKey: key)
    }

    func deactivate() {
        token = ""
        email = ""
        licenseKey = ""
        isActivated = false
        credentials.clear()
    }
}
