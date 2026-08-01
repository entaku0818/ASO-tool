//  APIClientAuthTests.swift
//  トークン失効(401)時の自動再アクティベートを検証する

import Testing
import Foundation
@testable import ASO_tool

// MARK: - URLProtocol スタブ

final class StubURLProtocol: URLProtocol {
    /// (受信リクエスト) -> (ステータス, ボディ) を返すハンドラ
    nonisolated(unsafe) static var handler: ((URLRequest) -> (Int, Data))?
    /// スタブが受け取ったリクエストの記録 (メソッド, パス, Authorizationヘッダ)
    nonisolated(unsafe) static var received: [(path: String, authorization: String?)] = []

    static func reset() {
        handler = nil
        received = []
    }

    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }
    override func stopLoading() {}

    override func startLoading() {
        Self.received.append((
            path: request.url?.path ?? "",
            authorization: request.value(forHTTPHeaderField: "Authorization")
        ))

        let (status, body) = Self.handler?(request) ?? (500, Data())
        let response = HTTPURLResponse(
            url: request.url!, statusCode: status, httpVersion: "HTTP/1.1", headerFields: nil)!
        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: body)
        client?.urlProtocolDidFinishLoading(self)
    }
}

/// UserDefaults.standard / NotificationCenter.default を共有すると
/// 並行実行される他スイート (AppState など) と干渉するため、専用のものを使う。
private func makeStubbedClient(store: CredentialStore) -> APIClient {
    let config = URLSessionConfiguration.ephemeral
    config.protocolClasses = [StubURLProtocol.self]
    let client = APIClient()
    client.baseURL = "https://api.test.local"
    client.session = URLSession(configuration: config)
    client.credentials = store
    client.notificationCenter = NotificationCenter()
    return client
}

private let activateResponseBody = """
{"token":"NEW_TOKEN","key":"ASOT-AAAA-BBBB-CCCC",
 "user":{"id":"u-1","email":"user@example.com","name":"Test User","plan":"pro"}}
""".data(using: .utf8)!

private let appsResponseBody = """
[{"id":"app-1","name":"MyApp","bundle_id":"com.example.app","platform":"ios",
  "store_url":null,"created_at":"2024-01-01T00:00:00Z"}]
""".data(using: .utf8)!

// StubURLProtocol の static を共有するため直列実行
@Suite("APIClient トークン失効時の自動再アクティベート", .serialized)
struct APIClientAuthTests {

    /// テストごとに独立した UserDefaults を作る (standard は他スイートと共有されるため使わない)
    private func makeStore(licenseKey: String?, email: String?) -> CredentialStore {
        let suiteName = "APIClientAuthTests.\(UUID().uuidString)"
        let defaults = UserDefaults(suiteName: suiteName)!
        let store = CredentialStore(defaults: defaults)
        store.update(token: "EXPIRED_TOKEN")
        if let licenseKey { defaults.set(licenseKey, forKey: CredentialStore.licenseKeyKey) }
        if let email { defaults.set(email, forKey: CredentialStore.emailKey) }
        return store
    }

    @Test("401 を受けたらライセンスで再アクティベートし、新トークンでリトライして成功する")
    func retriesAfterReactivation() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: "ASOT-AAAA-BBBB-CCCC", email: "user@example.com")
        defer { StubURLProtocol.reset() }

        StubURLProtocol.handler = { req in
            let path = req.url?.path ?? ""
            if path == "/api/licenses/activate" {
                return (200, activateResponseBody)
            }
            // 期限切れトークンなら401、新トークンなら200
            let auth = req.value(forHTTPHeaderField: "Authorization")
            return auth == "Bearer NEW_TOKEN"
                ? (200, appsResponseBody)
                : (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!)
        }

        let client = makeStubbedClient(store: store)
        let apps = try await client.getApps(token: "EXPIRED_TOKEN")

        #expect(apps.count == 1)
        #expect(apps.first?.id == "app-1")

        // 401 → activate → リトライ の3リクエストが飛ぶ
        #expect(StubURLProtocol.received.count == 3)
        #expect(StubURLProtocol.received[0].authorization == "Bearer EXPIRED_TOKEN")
        #expect(StubURLProtocol.received[1].path == "/api/licenses/activate")
        #expect(StubURLProtocol.received[2].authorization == "Bearer NEW_TOKEN")

        // 新トークンが保存される
        #expect(store.token == "NEW_TOKEN")
    }

    @Test("ライセンス情報がなければ再アクティベートせず 401 をそのまま返す")
    func surfacesUnauthorizedWithoutLicense() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: nil, email: nil)
        defer { StubURLProtocol.reset() }

        StubURLProtocol.handler = { _ in (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!) }

        let client = makeStubbedClient(store: store)
        await #expect(throws: APIClientError.self) {
            _ = try await client.getApps(token: "EXPIRED_TOKEN")
        }

        // activate を叩かず、リトライもしない
        #expect(StubURLProtocol.received.count == 1)
        #expect(StubURLProtocol.received.allSatisfy { $0.path != "/api/licenses/activate" })
    }

    @Test("再アクティベートも401なら、無限リトライせず失敗する")
    func doesNotLoopWhenReactivationFails() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: "ASOT-DEAD-BEEF-0000", email: "user@example.com")
        defer { StubURLProtocol.reset() }

        // activate も含めて全部 401 を返す
        StubURLProtocol.handler = { _ in (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!) }

        let client = makeStubbedClient(store: store)
        await #expect(throws: APIClientError.self) {
            _ = try await client.getApps(token: "EXPIRED_TOKEN")
        }

        // 元リクエスト + activate の2回で打ち止め
        #expect(StubURLProtocol.received.count == 2)
    }

    @Test("サーバー障害(5xx)で再アクティベートに失敗しても、保存済みライセンスは消さない")
    func keepsLicenseOnTransientFailure() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: "ASOT-AAAA-BBBB-CCCC", email: "user@example.com")
        defer { StubURLProtocol.reset() }

        StubURLProtocol.handler = { req in
            let path = req.url?.path ?? ""
            // activate は 503、業務APIは 401
            return path == "/api/licenses/activate"
                ? (503, #"{"error":"service unavailable"}"#.data(using: .utf8)!)
                : (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!)
        }

        let client = makeStubbedClient(store: store)
        let center = client.notificationCenter
        nonisolated(unsafe) var reactivationRequired = false
        let observer = center.addObserver(
            forName: .asoReactivationRequired, object: nil, queue: nil
        ) { _ in reactivationRequired = true }
        defer { center.removeObserver(observer) }

        await #expect(throws: APIClientError.self) {
            _ = try await client.getApps(token: "EXPIRED_TOKEN")
        }

        // 一時障害では手動再認証へ誘導せず、ライセンス情報も保持する
        #expect(reactivationRequired == false)
        #expect(store.licenseKey == "ASOT-AAAA-BBBB-CCCC")
        #expect(store.email == "user@example.com")
    }

    @Test("ライセンス無効(4xx)なら手動再認証を要求する")
    func requiresManualReactivationWhenLicenseRejected() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: "ASOT-DEAD-BEEF-0000", email: "user@example.com")
        defer { StubURLProtocol.reset() }

        StubURLProtocol.handler = { req in
            let path = req.url?.path ?? ""
            return path == "/api/licenses/activate"
                ? (404, #"{"error":"license not found"}"#.data(using: .utf8)!)
                : (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!)
        }

        let client = makeStubbedClient(store: store)
        let center = client.notificationCenter
        nonisolated(unsafe) var reactivationRequired = false
        let observer = center.addObserver(
            forName: .asoReactivationRequired, object: nil, queue: nil
        ) { _ in reactivationRequired = true }
        defer { center.removeObserver(observer) }

        await #expect(throws: APIClientError.self) {
            _ = try await client.getApps(token: "EXPIRED_TOKEN")
        }

        #expect(reactivationRequired == true)
    }

    @Test("同時に複数リクエストが401でも activate は1回しか呼ばれない")
    func coalescesConcurrentRefreshes() async throws {
        StubURLProtocol.reset()
        let store = makeStore(licenseKey: "ASOT-AAAA-BBBB-CCCC", email: "user@example.com")
        defer { StubURLProtocol.reset() }

        StubURLProtocol.handler = { req in
            let path = req.url?.path ?? ""
            if path == "/api/licenses/activate" { return (200, activateResponseBody) }
            let auth = req.value(forHTTPHeaderField: "Authorization")
            return auth == "Bearer NEW_TOKEN"
                ? (200, appsResponseBody)
                : (401, #"{"error":"unauthorized"}"#.data(using: .utf8)!)
        }

        let client = makeStubbedClient(store: store)
        try await withThrowingTaskGroup(of: [ASOApp].self) { group in
            for _ in 0..<5 {
                group.addTask { try await client.getApps(token: "EXPIRED_TOKEN") }
            }
            for try await apps in group {
                #expect(apps.count == 1)
            }
        }

        let activateCalls = StubURLProtocol.received.filter { $0.path == "/api/licenses/activate" }
        #expect(activateCalls.count == 1)
    }
}
