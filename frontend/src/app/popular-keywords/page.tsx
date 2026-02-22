'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import { getPopularKeywords, PopularKeyword, getAppRankings, AppRankingEntry } from '@/lib/api'

const COUNTRIES = [
  { code: 'jp', name: '日本', flag: '🇯🇵' },
  { code: 'us', name: 'アメリカ', flag: '🇺🇸' },
  { code: 'gb', name: 'イギリス', flag: '🇬🇧' },
  { code: 'de', name: 'ドイツ', flag: '🇩🇪' },
  { code: 'fr', name: 'フランス', flag: '🇫🇷' },
  { code: 'cn', name: '中国', flag: '🇨🇳' },
  { code: 'kr', name: '韓国', flag: '🇰🇷' },
]

const PLATFORMS = [
  { value: 'ios', label: 'iOS' },
  { value: 'android', label: 'Android' },
]

const LIMITS = [100, 500, 1000]

const RANKING_TYPES = [
  { value: 'topfreeapplications', label: '無料トップ' },
  { value: 'toppaidapplications', label: '有料トップ' },
  { value: 'topgrossingapplications', label: '売上トップ' },
  { value: 'newfreeapplications', label: '新着無料' },
  { value: 'newpaidapplications', label: '新着有料' },
]

const GENRES = [
  { id: '', name: '全ジャンル' },
  { id: '6014', name: 'ゲーム' },
  { id: '6015', name: 'ファイナンス' },
  { id: '6016', name: 'エンターテインメント' },
  { id: '6017', name: '教育' },
  { id: '6013', name: 'ヘルス&フィットネス' },
  { id: '6012', name: 'ライフスタイル' },
  { id: '6007', name: '仕事効率化' },
  { id: '6005', name: 'SNS' },
  { id: '6024', name: 'ショッピング' },
  { id: '6003', name: '旅行' },
  { id: '6023', name: 'フード&ドリンク' },
]

function KeywordsSection() {
  const [keywords, setKeywords] = useState<PopularKeyword[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [selectedCountry, setSelectedCountry] = useState('jp')
  const [selectedPlatform, setSelectedPlatform] = useState('ios')
  const [selectedLimit, setSelectedLimit] = useState(100)
  const [searchTerm, setSearchTerm] = useState('')

  useEffect(() => {
    fetchKeywords()
  }, [selectedCountry, selectedPlatform, selectedLimit])

  const fetchKeywords = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await getPopularKeywords(selectedCountry, selectedPlatform, selectedLimit)
      setKeywords(data)
    } catch (err) {
      setError('データの取得に失敗しました')
      console.error('Error fetching popular keywords:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const filteredKeywords = keywords.filter(kw =>
    kw.keyword.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const selectedCountryData = COUNTRIES.find(c => c.code === selectedCountry)

  return (
    <>
      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">国</label>
            <select
              value={selectedCountry}
              onChange={(e) => setSelectedCountry(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {COUNTRIES.map(country => (
                <option key={country.code} value={country.code}>
                  {country.flag} {country.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">プラットフォーム</label>
            <select
              value={selectedPlatform}
              onChange={(e) => setSelectedPlatform(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {PLATFORMS.map(platform => (
                <option key={platform.value} value={platform.value}>
                  {platform.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">表示件数</label>
            <select
              value={selectedLimit}
              onChange={(e) => setSelectedLimit(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {LIMITS.map(limit => (
                <option key={limit} value={limit}>{limit}件</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">検索</label>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="キーワードを検索..."
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Results */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-6 py-4 border-b bg-gray-50">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">
              {selectedCountryData?.flag} {selectedCountryData?.name} - {selectedPlatform.toUpperCase()}
            </h2>
            <span className="text-sm text-gray-500">{filteredKeywords.length}件のキーワード</span>
          </div>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500">読み込み中...</div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-red-500">{error}</div>
          </div>
        ) : filteredKeywords.length === 0 ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500">キーワードが見つかりませんでした</div>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="bg-gray-50 border-b text-left text-sm font-medium text-gray-500">
                  <th className="px-6 py-3">#</th>
                  <th className="px-6 py-3">キーワード</th>
                  <th className="px-6 py-3 text-right">追跡数</th>
                  <th className="px-6 py-3 text-right">検索結果数</th>
                </tr>
              </thead>
              <tbody>
                {filteredKeywords.map((keyword, index) => (
                  <tr
                    key={`${keyword.keyword}-${index}`}
                    className="border-b hover:bg-gray-50 transition-colors"
                  >
                    <td className="px-6 py-4 text-sm text-gray-500">{index + 1}</td>
                    <td className="px-6 py-4">
                      <span className="font-medium text-gray-900">{keyword.keyword}</span>
                    </td>
                    <td className="px-6 py-4 text-right text-sm text-gray-600">{keyword.tracking_count}</td>
                    <td className="px-6 py-4 text-right text-sm text-gray-600">{keyword.results_count}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  )
}

function AppRankingSection() {
  const [rankings, setRankings] = useState<AppRankingEntry[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [selectedCountry, setSelectedCountry] = useState('jp')
  const [selectedRankingType, setSelectedRankingType] = useState('topfreeapplications')
  const [selectedGenre, setSelectedGenre] = useState('')

  useEffect(() => {
    fetchRankings()
  }, [selectedCountry, selectedRankingType, selectedGenre])

  const fetchRankings = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await getAppRankings(selectedCountry, selectedRankingType, selectedGenre, 100)
      setRankings(data)
    } catch (err) {
      setError('ランキングの取得に失敗しました')
      console.error('Error fetching app rankings:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const selectedCountryData = COUNTRIES.find(c => c.code === selectedCountry)
  const selectedRankingTypeData = RANKING_TYPES.find(r => r.value === selectedRankingType)
  const selectedGenreData = GENRES.find(g => g.id === selectedGenre)

  return (
    <>
      {/* Ranking type sub-tabs */}
      <div className="flex gap-2 mb-4 overflow-x-auto pb-1">
        {RANKING_TYPES.map(type => (
          <button
            key={type.value}
            onClick={() => setSelectedRankingType(type.value)}
            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-colors ${
              selectedRankingType === type.value
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-600 border border-gray-300 hover:bg-gray-50'
            }`}
          >
            {type.label}
          </button>
        ))}
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">国</label>
            <select
              value={selectedCountry}
              onChange={(e) => setSelectedCountry(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {COUNTRIES.map(country => (
                <option key={country.code} value={country.code}>
                  {country.flag} {country.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">ジャンル</label>
            <select
              value={selectedGenre}
              onChange={(e) => setSelectedGenre(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {GENRES.map(genre => (
                <option key={genre.id} value={genre.id}>{genre.name}</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Results */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-6 py-4 border-b bg-gray-50">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">
              {selectedCountryData?.flag} {selectedCountryData?.name} — {selectedRankingTypeData?.label}
              {selectedGenreData?.id && ` / ${selectedGenreData.name}`}
            </h2>
            <span className="text-sm text-gray-500">{rankings.length}件</span>
          </div>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500">読み込み中...</div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-red-500">{error}</div>
          </div>
        ) : rankings.length === 0 ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500">ランキングが見つかりませんでした</div>
          </div>
        ) : (
          <div className="divide-y">
            {rankings.map((app) => (
              <div key={app.app_id} className="flex items-center gap-4 px-6 py-4 hover:bg-gray-50 transition-colors">
                <span className="w-8 text-center text-sm font-bold text-gray-400">{app.rank}</span>
                {app.icon_url ? (
                  <Image
                    src={app.icon_url}
                    alt={app.name}
                    width={48}
                    height={48}
                    className="rounded-xl flex-shrink-0"
                    unoptimized
                  />
                ) : (
                  <div className="w-12 h-12 bg-gray-200 rounded-xl flex-shrink-0" />
                )}
                <div className="flex-1 min-w-0">
                  <a
                    href={app.store_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-medium text-gray-900 hover:text-blue-600 truncate block"
                  >
                    {app.name}
                  </a>
                  <p className="text-sm text-gray-500 truncate">{app.developer}</p>
                </div>
                <div className="text-right flex-shrink-0">
                  <p className="text-sm text-gray-600">{app.category}</p>
                  <p className="text-xs text-gray-400">{app.price || '無料'}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  )
}

export default function PopularKeywordsPage() {
  const [activeTab, setActiveTab] = useState<'keywords' | 'rankings'>('keywords')

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-3xl font-bold text-gray-900">App Store データ</h1>
          <p className="mt-2 text-sm text-gray-600">
            人気キーワードとApp Storeランキングを確認する
          </p>
        </div>
      </div>

      {/* Main tabs */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-0">
            <button
              onClick={() => setActiveTab('keywords')}
              className={`px-6 py-4 text-sm font-medium border-b-2 transition-colors ${
                activeTab === 'keywords'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              人気キーワード
            </button>
            <button
              onClick={() => setActiveTab('rankings')}
              className={`px-6 py-4 text-sm font-medium border-b-2 transition-colors ${
                activeTab === 'rankings'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              App Store ランキング
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 pb-12">
        {activeTab === 'keywords' ? <KeywordsSection /> : <AppRankingSection />}
      </div>
    </div>
  )
}
