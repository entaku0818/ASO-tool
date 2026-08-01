//  APIClient.swift

import Foundation
import os.log

private let logger = Logger(subsystem: "com.entaku.ASO-tool", category: "APIClient")

enum APIClientError: LocalizedError {
    case invalidURL
    case httpError(Int, String)
    case decodingError(Error)
    case networkError(Error)
    case reactivationRequired

    var errorDescription: String? {
        switch self {
        case .invalidURL: return "Invalid URL"
        case .httpError(let code, let msg): return "HTTP \(code): \(msg)"
        case .decodingError(let e): return "デコードエラー: \(e.localizedDescription)"
        case .networkError(let e): return e.localizedDescription
        case .reactivationRequired: return "ライセンスの再認証が必要です。設定画面からライセンスキーを入力し直してください。"
        }
    }
}

final class APIClient {
    static let shared = APIClient()

    var baseURL: String = "https://aso-api-671942133800.asia-northeast1.run.app"
    var session: URLSession = .shared
    /// トークン失効時の自動再アクティベートに使う認証情報 (テストでは差し替える)
    var credentials: CredentialStore = .standard
    var notificationCenter: NotificationCenter = .default

    private let refresher = TokenRefresher()

    private lazy var decoder: JSONDecoder = {
        let d = JSONDecoder()
        d.dateDecodingStrategy = .iso8601
        return d
    }()

    // MARK: - License

    func activateLicense(key: String, email: String) async throws -> ActivateLicenseResponse {
        try await post("/api/licenses/activate", body: ["key": key, "email": email], token: nil)
    }

    // MARK: - Apps

    func getApps(token: String) async throws -> [ASOApp] {
        try await get("/api/apps", token: token)
    }

    func searchApps(token: String, keyword: String, platform: String, country: String, limit: Int = 20) async throws -> [AppStoreSearchResult] {
        let q = keyword.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? keyword
        return try await get("/api/scraper/search?keyword=\(q)&platform=\(platform)&country=\(country)&limit=\(limit)", token: token)
    }

    func createApp(token: String, name: String, bundleID: String, platform: String, storeURL: String?) async throws -> ASOApp {
        struct Body: Encodable {
            let name: String
            let bundle_id: String
            let platform: String
            let store_url: String?
        }
        return try await post("/api/apps", body: Body(name: name, bundle_id: bundleID, platform: platform, store_url: storeURL.flatMap { $0.isEmpty ? nil : $0 }), token: token)
    }

    // MARK: - Keywords

    func getKeywords(token: String, appID: String) async throws -> [Keyword] {
        try await get("/api/apps/\(appID)/keywords", token: token)
    }

    func importKeywords(token: String, appID: String, keywords: [[String: String]]) async throws -> ImportKeywordsResponse {
        struct Body: Encodable { let keywords: [[String: String]] }
        return try await post("/api/apps/\(appID)/keywords/import", body: Body(keywords: keywords), token: token)
    }

    func getKeywordSuggestions(token: String, appID: String, limit: Int = 20) async throws -> [KeywordSuggestion] {
        try await get("/api/apps/\(appID)/keywords/suggestions?limit=\(limit)", token: token)
    }

    func getCompetitorSuggestions(token: String, appID: String, adamID: String, limit: Int = 50) async throws -> [KeywordSuggestion] {
        try await get("/api/apps/\(appID)/keywords/competitor-suggestions?adam_id=\(adamID)&limit=\(limit)", token: token)
    }

    func getKeywordAutocomplete(token: String, term: String, country: String = "jp") async throws -> [AutocompleteSuggestion] {
        let q = term.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? term
        return try await get("/api/scraper/keyword-suggestions?term=\(q)&country=\(country.lowercased())", token: token)
    }

    func createKeyword(token: String, appID: String, keyword: String, country: String) async throws -> Keyword {
        let body = ["app_id": appID, "keyword": keyword, "country": country]
        return try await post("/api/apps/\(appID)/keywords", body: body, token: token)
    }

    func deleteKeyword(token: String, appID: String, keywordID: String) async throws {
        try await delete("/api/apps/\(appID)/keywords/\(keywordID)", token: token)
    }

    // MARK: - Rankings

    func getRankings(token: String, appID: String, keywordID: String, limit: Int = 30) async throws -> [RankingHistory] {
        try await get("/api/apps/\(appID)/keywords/\(keywordID)/rankings?limit=\(limit)", token: token)
    }

    func scrapeRankings(token: String, appID: String) async throws {
        struct Res: Decodable { let updated: Int }
        let _: Res = try await post("/api/apps/\(appID)/scrape/rankings", body: EmptyBody(), token: token)
    }

    func getKeywordRanks(token: String, appID: String) async throws -> [KeywordRankSummary] {
        try await get("/api/apps/\(appID)/keywords/ranks", token: token)
    }

    func getRisingKeywords(token: String, appID: String, days: Int = 7) async throws -> [RisingKeyword] {
        try await get("/api/apps/\(appID)/keywords/rising?days=\(days)", token: token)
    }

    func getLatestRanking(token: String, appID: String, keywordID: String) async throws -> RankingHistory? {
        do {
            return try await get("/api/apps/\(appID)/keywords/\(keywordID)/rankings/latest", token: token)
        } catch APIClientError.httpError(404, _) {
            return nil
        }
    }

    func deleteApp(token: String, appID: String) async throws {
        try await delete("/api/apps/\(appID)", token: token)
    }

    // MARK: - Competitors

    func getCompetitors(token: String, appID: String) async throws -> [Competitor] {
        try await get("/api/apps/\(appID)/competitors", token: token)
    }

    func createCompetitor(token: String, appID: String, bundleID: String, name: String, platform: String) async throws -> Competitor {
        struct Body: Encodable {
            let app_id: String; let competitor_bundle_id: String
            let competitor_name: String; let platform: String
        }
        return try await post("/api/apps/\(appID)/competitors",
            body: Body(app_id: appID, competitor_bundle_id: bundleID, competitor_name: name, platform: platform),
            token: token)
    }

    func deleteCompetitor(token: String, appID: String, competitorID: String) async throws {
        try await delete("/api/apps/\(appID)/competitors/\(competitorID)", token: token)
    }

    func getKeywordGap(token: String, appID: String) async throws -> [KeywordGap] {
        try await get("/api/apps/\(appID)/competitors/keyword-gap", token: token)
    }

    func updateCompetitorRankings(token: String, appID: String) async throws -> UpdateRankingsResponse {
        try await post("/api/apps/\(appID)/competitors/update-rankings", body: EmptyBody(), token: token)
    }

    // MARK: - ASO Advice

    func getASOAdvice(token: String, appID: String) async throws -> ASOAdviceResponse {
        try await get("/api/apps/\(appID)/aso-advice", token: token)
    }

    // MARK: - Metadata

    func getMetadata(token: String, appID: String) async throws -> [AppMetadataVersion] {
        try await get("/api/apps/\(appID)/metadata", token: token)
    }

    func upsertMetadata(token: String, appID: String, request: UpsertMetadataRequest) async throws -> AppMetadataVersion {
        try await put("/api/apps/\(appID)/metadata", body: request, token: token)
    }

    func deleteMetadata(token: String, appID: String, metadataID: String) async throws {
        try await delete("/api/apps/\(appID)/metadata/\(metadataID)", token: token)
    }

    // MARK: - Helpers

    func get<T: Decodable>(_ path: String, token: String?) async throws -> T {
        try await perform(makeRequest("GET", path: path, body: nil as [String: String]?, token: token))
    }

    func post<B: Encodable, T: Decodable>(_ path: String, body: B, token: String?) async throws -> T {
        try await perform(makeRequest("POST", path: path, body: body, token: token))
    }

    func put<B: Encodable, T: Decodable>(_ path: String, body: B, token: String?) async throws -> T {
        try await perform(makeRequest("PUT", path: path, body: body, token: token))
    }

    func delete(_ path: String, token: String?) async throws {
        let req = try makeRequest("DELETE", path: path, body: nil as [String: String]?, token: token)
        let (_, http) = try await send(req)
        if let http, !(200..<300).contains(http.statusCode) {
            throw APIClientError.httpError(http.statusCode, "delete failed")
        }
    }

    private struct EmptyBody: Encodable {}

    private func makeRequest<B: Encodable>(_ method: String, path: String, body: B?, token: String?) throws -> URLRequest {
        guard let url = URL(string: baseURL + path) else { throw APIClientError.invalidURL }
        var req = URLRequest(url: url)
        req.httpMethod = method
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let t = token { req.setValue("Bearer \(t)", forHTTPHeaderField: "Authorization") }
        if let b = body { req.httpBody = try JSONEncoder().encode(b) }
        return req
    }

    /// リクエストを送る。401 が返り、かつ保存済みライセンスで再アクティベートできる場合は
    /// トークンを更新して1度だけリトライする。
    private func send(_ req: URLRequest) async throws -> (Data, HTTPURLResponse?) {
        let (data, http) = try await rawSend(req)

        guard http?.statusCode == 401,
              let staleToken = bearerToken(of: req) else {
            return (data, http)
        }

        logger.info("401 を受信。ライセンスの自動再アクティベートを試行")
        let newToken: String
        do {
            newToken = try await refreshToken(staleToken: staleToken)
        } catch {
            logger.error("自動再アクティベート失敗: \(error, privacy: .public)")
            return (data, http)
        }

        var retried = req
        retried.setValue("Bearer \(newToken)", forHTTPHeaderField: "Authorization")
        logger.info("新トークンでリトライ: \(req.url?.absoluteString ?? "?", privacy: .public)")
        return try await rawSend(retried)
    }

    private func rawSend(_ req: URLRequest) async throws -> (Data, HTTPURLResponse?) {
        logger.info("→ \(req.httpMethod ?? "?", privacy: .public) \(req.url?.absoluteString ?? "?", privacy: .public)")
        let data: Data
        let response: URLResponse
        do {
            (data, response) = try await session.data(for: req)
        } catch {
            logger.error("networkError: \(error, privacy: .public)")
            throw APIClientError.networkError(error)
        }
        let http = response as? HTTPURLResponse
        logger.info("← status: \(http?.statusCode ?? -1, privacy: .public)")
        logger.debug("← body: \(String(data: data, encoding: .utf8) ?? "(binary)", privacy: .public)")
        return (data, http)
    }

    private func bearerToken(of req: URLRequest) -> String? {
        guard let header = req.value(forHTTPHeaderField: "Authorization"),
              header.hasPrefix("Bearer ") else { return nil }
        return String(header.dropFirst("Bearer ".count))
    }

    /// 保存済みのライセンスキー＋メールで再アクティベートし、新しいトークンを保存して返す。
    /// 認証ヘッダを付けない activate エンドポイントを使うため、401 リトライが再帰することはない。
    private func refreshToken(staleToken: String?) async throws -> String {
        let store = credentials
        let center = notificationCenter

        return try await refresher.refresh(
            staleToken: staleToken, currentToken: store.token
        ) { [self] in
            guard store.canReactivate else {
                await post(.asoReactivationRequired, on: center)
                throw APIClientError.reactivationRequired
            }

            do {
                let resp = try await activateLicense(key: store.licenseKey, email: store.email)
                store.update(token: resp.token)
                await post(.asoTokenRefreshed, on: center, userInfo: ["token": resp.token])
                return resp.token
            } catch {
                // ライセンス失効・無効 (4xx) のときだけ手動再認証へ誘導する。
                // オフラインやサーバー障害 (5xx) で保存済みライセンスを捨てないこと。
                if Self.isPermanentRejection(error) {
                    await post(.asoReactivationRequired, on: center)
                }
                throw error
            }
        }
    }

    /// ライセンスが恒久的に使えない (無効・失効・別アカウントで使用済み) 応答かどうか。
    /// 通信断やサーバー障害と区別して、保存済みライセンスを消してよいかの判断に使う。
    private static func isPermanentRejection(_ error: Error) -> Bool {
        guard case .httpError(let code, _) = error as? APIClientError ?? .invalidURL else {
            return false
        }
        return (400..<500).contains(code)
    }

    @MainActor
    private func post(
        _ name: Notification.Name, on center: NotificationCenter, userInfo: [String: Any]? = nil
    ) {
        center.post(name: name, object: nil, userInfo: userInfo)
    }

    private func perform<T: Decodable>(_ req: URLRequest) async throws -> T {
        let (data, http) = try await send(req)
        if let http = http, !(200..<300).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode(APIErrorResponse.self, from: data))?.error ?? "unknown error"
            logger.error("httpError \(http.statusCode, privacy: .public): \(msg, privacy: .public)")
            throw APIClientError.httpError(http.statusCode, msg)
        }
        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            logger.error("decodeError: \(error, privacy: .public)")
            throw APIClientError.decodingError(error)
        }
    }
}
