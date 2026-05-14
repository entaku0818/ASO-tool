//  ASOModels.swift

import Foundation

struct ASOApp: Codable, Identifiable, Hashable {
    let id: String
    let name: String
    let bundleID: String
    let platform: String
    let storeURL: String?
    let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, name, platform
        case bundleID = "bundle_id"
        case storeURL = "store_url"
        case createdAt = "created_at"
    }
}

struct Keyword: Codable, Identifiable, Hashable {
    let id: String
    let appID: String
    let keyword: String
    let country: String
    let popularityScore: Int?
    let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, keyword, country
        case appID = "app_id"
        case popularityScore = "popularity_score"
        case createdAt = "created_at"
    }
}

struct RankingHistory: Codable, Identifiable {
    let id: String
    let keywordID: String
    let rank: Int?
    let recordedAt: Date

    enum CodingKeys: String, CodingKey {
        case id, rank
        case keywordID = "keyword_id"
        case recordedAt = "recorded_at"
    }
}

struct ASOUser: Codable {
    let id: String
    let email: String
    let name: String
    let plan: String
}

struct ActivateLicenseResponse: Codable {
    let token: String
    let user: ASOUser
    let key: String
}

struct APIErrorResponse: Codable {
    let error: String
}
