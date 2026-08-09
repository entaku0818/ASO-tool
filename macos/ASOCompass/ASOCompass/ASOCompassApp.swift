//  ASOCompassApp.swift

import SwiftUI

@main
struct ASOCompassApp: App {
    @StateObject private var appState = AppState()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .frame(minWidth: 900, minHeight: 600)
        }
        .windowResizability(.contentMinSize)

        Settings {
            SettingsView()
                .environmentObject(appState)
        }
    }
}
