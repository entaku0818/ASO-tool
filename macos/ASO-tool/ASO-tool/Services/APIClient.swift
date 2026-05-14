//  APIClient.swift

import Foundation

enum APIClientError: LocalizedError {
    case invalidURL
    case httpError(Int, String)
    case decodingError(Error)
    case networkError(Error)

    var errorDescription: String? {
        switch self {
        case .invalidURL: return "Invalid URL"
        case .httpError(let code, let msg): return "HTTP \(code): \(msg)"
        case .decodingError(let e): return "デコードエラー: \(e.localizedDescription)"
        case .networkError(let e): return e.localizedDescription
        }
    }
}

final class APIClient {
    static let shared = APIClient()

    var baseURL: String = "https://aso-tool.vercel.app"
    var session: URLSession = .shared

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

    // MARK: - Keywords

    func getKeywords(token: String, appID: String) async throws -> [Keyword] {
        try await get("/api/apps/\(appID)/keywords", token: token)
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

    // MARK: - Helpers

    func get<T: Decodable>(_ path: String, token: String?) async throws -> T {
        try await perform(makeRequest("GET", path: path, body: nil as [String: String]?, token: token))
    }

    func post<B: Encodable, T: Decodable>(_ path: String, body: B, token: String?) async throws -> T {
        try await perform(makeRequest("POST", path: path, body: body, token: token))
    }

    func delete(_ path: String, token: String?) async throws {
        let req = try makeRequest("DELETE", path: path, body: nil as [String: String]?, token: token)
        let (_, response) = try await session.data(for: req)
        if let http = response as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
            throw APIClientError.httpError(http.statusCode, "delete failed")
        }
    }

    private func makeRequest<B: Encodable>(_ method: String, path: String, body: B?, token: String?) throws -> URLRequest {
        guard let url = URL(string: baseURL + path) else { throw APIClientError.invalidURL }
        var req = URLRequest(url: url)
        req.httpMethod = method
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let t = token { req.setValue("Bearer \(t)", forHTTPHeaderField: "Authorization") }
        if let b = body { req.httpBody = try JSONEncoder().encode(b) }
        return req
    }

    private func perform<T: Decodable>(_ req: URLRequest) async throws -> T {
        let data: Data
        let response: URLResponse
        do {
            (data, response) = try await session.data(for: req)
        } catch {
            throw APIClientError.networkError(error)
        }
        if let http = response as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode(APIErrorResponse.self, from: data))?.error ?? "unknown error"
            throw APIClientError.httpError(http.statusCode, msg)
        }
        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIClientError.decodingError(error)
        }
    }
}
