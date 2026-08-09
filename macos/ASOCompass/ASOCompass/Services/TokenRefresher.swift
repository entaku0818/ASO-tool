//  TokenRefresher.swift

import Foundation

/// トークン失効時の再アクティベートを直列化する。
/// アプリ起動直後のように複数リクエストが同時に401を受けても、
/// /api/licenses/activate を叩くのは1回だけにする。
actor TokenRefresher {
    private var inFlight: Task<String, Error>?

    /// - Parameters:
    ///   - staleToken: 401 になったリクエストが使っていたトークン。
    ///   - currentToken: 現在保存されているトークン。staleToken と違えば
    ///                   別のリクエストが既に更新済みなので、それをそのまま使う。
    ///   - operation: 実際に再アクティベートしてトークン文字列を返す処理。
    func refresh(
        staleToken: String?,
        currentToken: String,
        operation: @escaping @Sendable () async throws -> String
    ) async throws -> String {
        if let staleToken, !currentToken.isEmpty, currentToken != staleToken {
            return currentToken
        }

        if let existing = inFlight {
            return try await existing.value
        }

        let task = Task<String, Error> { try await operation() }
        inFlight = task
        defer { inFlight = nil }
        return try await task.value
    }
}
