//  ASO_toolTests.swift

import Testing
import Foundation
@testable import ASO_tool

// MARK: - AppState Tests

@Suite("AppState")
struct AppStateTests {

    @Test("初期状態: token がなければ非アクティブ")
    func initialState() {
        UserDefaults.standard.removeObject(forKey: "jwt_token")
        let state = AppState()
        #expect(state.isActivated == false)
        #expect(state.token == "")
    }

    @Test("activate: 全フィールドが設定される")
    func activate() {
        let state = AppState()
        state.activate(token: "tok123", email: "user@example.com", key: "ASOT-AAAA-BBBB-CCCC")

        #expect(state.isActivated == true)
        #expect(state.token == "tok123")
        #expect(state.email == "user@example.com")
        #expect(state.licenseKey == "ASOT-AAAA-BBBB-CCCC")

        state.deactivate()
    }

    @Test("activate: UserDefaults に永続化される")
    func activatePersists() {
        let state = AppState()
        state.activate(token: "tok123", email: "user@example.com", key: "ASOT-AAAA-BBBB-CCCC")

        #expect(UserDefaults.standard.string(forKey: "jwt_token") == "tok123")
        #expect(UserDefaults.standard.string(forKey: "user_email") == "user@example.com")
        #expect(UserDefaults.standard.string(forKey: "license_key") == "ASOT-AAAA-BBBB-CCCC")

        state.deactivate()
    }

    @Test("deactivate: 全フィールドがクリアされる")
    func deactivate() {
        let state = AppState()
        state.activate(token: "tok123", email: "user@example.com", key: "ASOT-AAAA-BBBB-CCCC")
        state.deactivate()

        #expect(state.isActivated == false)
        #expect(state.token == "")
        #expect(state.email == "")
        #expect(state.licenseKey == "")
        #expect(UserDefaults.standard.string(forKey: "jwt_token") == nil)
    }
}

// MARK: - ASOModels Decoding Tests

@Suite("ASOModels デコード")
struct ASOModelsTests {

    @Test("ASOApp のデコード")
    func decodeASOApp() throws {
        let json = """
        {"id":"app-1","name":"MyApp","bundle_id":"com.example.app","platform":"ios",
         "store_url":"https://apps.apple.com/app/id123","created_at":"2024-01-01T00:00:00Z"}
        """.data(using: .utf8)!

        let app = try JSONDecoder().decode(ASOApp.self, from: json)
        #expect(app.id == "app-1")
        #expect(app.name == "MyApp")
        #expect(app.bundleID == "com.example.app")
        #expect(app.platform == "ios")
    }

    @Test("Keyword: popularity_score あり")
    func decodeKeywordWithScore() throws {
        let json = """
        {"id":"kw-1","app_id":"app-1","keyword":"アプリ 無料","country":"JP",
         "popularity_score":4,"created_at":"2024-01-01T00:00:00Z"}
        """.data(using: .utf8)!

        let kw = try JSONDecoder().decode(Keyword.self, from: json)
        #expect(kw.keyword == "アプリ 無料")
        #expect(kw.popularityScore == 4)
    }

    @Test("Keyword: popularity_score なし")
    func decodeKeywordWithoutScore() throws {
        let json = """
        {"id":"kw-2","app_id":"app-1","keyword":"test","country":"US",
         "created_at":"2024-01-01T00:00:00Z"}
        """.data(using: .utf8)!

        let kw = try JSONDecoder().decode(Keyword.self, from: json)
        #expect(kw.popularityScore == nil)
    }

    @Test("RankingHistory: rank あり")
    func decodeRankingWithRank() throws {
        let json = """
        {"id":"r-1","keyword_id":"kw-1","rank":5,"recorded_at":"2024-06-01T12:00:00Z"}
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let r = try decoder.decode(RankingHistory.self, from: json)
        #expect(r.rank == 5)
        #expect(r.keywordID == "kw-1")
    }

    @Test("RankingHistory: rank null")
    func decodeRankingNullRank() throws {
        let json = """
        {"id":"r-2","keyword_id":"kw-1","rank":null,"recorded_at":"2024-06-01T12:00:00Z"}
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let r = try decoder.decode(RankingHistory.self, from: json)
        #expect(r.rank == nil)
    }

    @Test("ActivateLicenseResponse のデコード")
    func decodeActivateLicenseResponse() throws {
        let json = """
        {"token":"eyJ...","key":"ASOT-AAAA-BBBB-CCCC",
         "user":{"id":"u-1","email":"user@example.com","name":"Test User","plan":"free"}}
        """.data(using: .utf8)!

        let resp = try JSONDecoder().decode(ActivateLicenseResponse.self, from: json)
        #expect(resp.token == "eyJ...")
        #expect(resp.key == "ASOT-AAAA-BBBB-CCCC")
        #expect(resp.user.email == "user@example.com")
    }
}
