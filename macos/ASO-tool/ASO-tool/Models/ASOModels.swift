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

struct ImportKeywordsResponse: Codable {
    let imported: Int
    let skipped: Int
}

struct KeywordSuggestion: Codable, Identifiable {
    let text: String
    let popularityScore: Int
    var id: String { text }
}

struct AutocompleteSuggestion: Codable, Identifiable {
    let term: String
    var id: String { term }
}

struct KeywordRankSummary: Codable, Identifiable {
    let keywordID: String
    let keyword: String
    let country: String
    let currentRank: Int?
    let previousRank: Int?
    let change: Int?
    var id: String { keywordID }

    enum CodingKeys: String, CodingKey {
        case keyword, country
        case keywordID    = "keyword_id"
        case currentRank  = "current_rank"
        case previousRank = "previous_rank"
        case change
    }
}

struct RisingKeyword: Codable, Identifiable {
    let keywordID: String
    let keyword: String
    let country: String
    let currentRank: Int
    let previousRank: Int
    let improvement: Int
    var id: String { keywordID }

    enum CodingKeys: String, CodingKey {
        case keyword, country, improvement
        case keywordID    = "keyword_id"
        case currentRank  = "current_rank"
        case previousRank = "previous_rank"
    }
}

struct AppStoreSearchResult: Codable, Identifiable {
    let rank: Int
    let appInfo: AppStoreAppInfo
    var id: String { appInfo.bundleID }

    enum CodingKeys: String, CodingKey {
        case rank = "Rank"
        case appInfo = "AppInfo"
    }
}

struct AppStoreAppInfo: Codable {
    let bundleID: String
    let name: String
    let developer: String
    let rating: Double
    let ratingCount: Int
    let storeURL: String
    let iconURL: String
    let price: Double

    enum CodingKeys: String, CodingKey {
        case bundleID = "BundleID"
        case name = "Name"
        case developer = "Developer"
        case rating = "Rating"
        case ratingCount = "RatingCount"
        case storeURL = "StoreURL"
        case iconURL = "IconURL"
        case price = "Price"
    }
}

// MARK: - Competitor

struct Competitor: Codable, Identifiable, Hashable {
    let id: String
    let appID: String
    let competitorBundleID: String
    let competitorName: String
    let platform: String
    let notes: String?
    let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, platform, notes
        case appID = "app_id"
        case competitorBundleID = "competitor_bundle_id"
        case competitorName = "competitor_name"
        case createdAt = "created_at"
    }
}

struct KeywordGap: Codable, Identifiable {
    let keywordID: String
    let keyword: String
    let country: String
    let competitorName: String
    let competitorRank: Int
    let ourRank: Int?

    var id: String { keywordID }

    enum CodingKeys: String, CodingKey {
        case keyword, country
        case keywordID = "keyword_id"
        case competitorName = "competitor_name"
        case competitorRank = "competitor_rank"
        case ourRank = "our_rank"
    }
}

/// POST /api/apps/{appID}/competitors/update-rankings のレスポンス
struct UpdateRankingsResponse: Codable {
    let updated: Int
}

/// POST /api/apps/{appID}/scrape/rankings のレスポンス。
/// 競合側の update-rankings とはキー名が違う (`keywords_updated`) ので別型にしている。
struct ScrapeRankingsResponse: Codable {
    let message: String
    let keywordsUpdated: Int

    enum CodingKeys: String, CodingKey {
        case message
        case keywordsUpdated = "keywords_updated"
    }
}

// MARK: - Metadata

struct AppMetadataVersion: Codable, Identifiable, Hashable {
    let id: String
    let appID: String
    let locale: String
    let versionTag: String
    let title: String?
    let subtitle: String?
    let description: String?
    let keywords: String?
    let promotionalText: String?
    let createdAt: String
    let updatedAt: String

    enum CodingKeys: String, CodingKey {
        case id, locale, title, subtitle, description, keywords
        case appID = "app_id"
        case versionTag = "version_tag"
        case promotionalText = "promotional_text"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

struct UpsertMetadataRequest: Codable {
    let locale: String
    let versionTag: String
    let title: String?
    let subtitle: String?
    let description: String?
    let keywords: String?
    let promotionalText: String?

    enum CodingKeys: String, CodingKey {
        case locale, title, subtitle, description, keywords
        case versionTag = "version_tag"
        case promotionalText = "promotional_text"
    }
}

// MARK: - ASO Advice

struct ASOAdvice: Codable, Identifiable {
    let priority: String
    let category: String
    let title: String
    let description: String
    let action: String

    var id: String { "\(priority)-\(title)" }
}

struct ASOAdviceResponse: Codable {
    let summary: String
    let advice: [ASOAdvice]
}
