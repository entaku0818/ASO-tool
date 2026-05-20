'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import {
  getAppStoreSuggestions,
  AppStoreKeywordSuggestion,
  getAppRankings,
  AppRankingEntry,
  getAppRankingTrend,
  RankingTrendPoint,
  getAppRankingCountries,
  CountryRankPoint,
  getKeywordCache,
  CachedKeyword,
} from '@/lib/api'

const COUNTRIES = [
  { code: 'jp', name: '日本', flag: '🇯🇵' },
  { code: 'us', name: 'アメリカ', flag: '🇺🇸' },
  { code: 'gb', name: 'イギリス', flag: '🇬🇧' },
  { code: 'de', name: 'ドイツ', flag: '🇩🇪' },
  { code: 'fr', name: 'フランス', flag: '🇫🇷' },
  { code: 'cn', name: '中国', flag: '🇨🇳' },
  { code: 'kr', name: '韓国', flag: '🇰🇷' },
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

const TREND_RANKING_TYPES = [
  { value: 'topfreeapplications', label: '無料トップ' },
  { value: 'toppaidapplications', label: '有料トップ' },
  { value: 'topgrossingapplications', label: '売上トップ' },
]

const COUNTRY_FLAGS: Record<string, string> = {
  jp: '🇯🇵',
  us: '🇺🇸',
  kr: '🇰🇷',
  gb: '🇬🇧',
}

const COUNTRY_NAMES: Record<string, string> = {
  jp: '日本',
  us: 'アメリカ',
  kr: '韓国',
  gb: 'イギリス',
}

const SEED_KEYWORDS = [
  'ゲーム', '写真', '音楽', '動画', '家計簿', 'ダイエット', '勉強', '天気',
  '地図', 'ニュース', 'カメラ', 'マンガ', '英語', 'レシピ', '節約', '睡眠',
  'game', 'photo', 'music', 'chat', 'camera', 'map', 'weather',
  'fitness', 'manga', 'video', 'shopping', 'finance', 'travel',
]

const ASA_GENRES = [
  { value: '', label: '全ジャンル' },
  { value: 'social', label: 'SNS・コミュニケーション' },
  { value: 'video', label: '動画' },
  { value: 'navigation', label: 'ナビ・地図' },
  { value: 'photo', label: '写真' },
  { value: 'shopping', label: 'ショッピング' },
]

function PopularityBar({ value }: { value: number }) {
  const pct = Math.min(100, Math.max(0, value))
  const color = pct >= 80 ? 'bg-green-500' : pct >= 50 ? 'bg-blue-500' : pct >= 20 ? 'bg-yellow-500' : 'bg-gray-400'
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
        <div className={`h-full rounded-full ${color}`} style={{ width: `${pct}%` }} />
      </div>
      <span className="text-xs font-mono text-gray-600 dark:text-gray-400 w-8 text-right">{pct}</span>
    </div>
  )
}

const PLATFORMS = [
  { value: 'ios', label: '🍎 App Store' },
  { value: 'android', label: '🤖 Google Play' },
]

function KeywordsSection() {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCountry, setSelectedCountry] = useState('jp')
  const [selectedPlatform, setSelectedPlatform] = useState('ios')
  const [suggestions, setSuggestions] = useState<AppStoreKeywordSuggestion[]>([])
  const [isSearching, setIsSearching] = useState(false)
  const [hasSearched, setHasSearched] = useState(false)

  // Cached popularity state (iOS only)
  const [cachedKeywords, setCachedKeywords] = useState<CachedKeyword[]>([])
  const [cacheCountry, setCacheCountry] = useState('jp')
  const [cacheGenre, setCacheGenre] = useState('')
  const [isCacheLoading, setIsCacheLoading] = useState(false)

  useEffect(() => {
    if (!searchTerm.trim()) {
      setSuggestions([])
      setHasSearched(false)
      return
    }
    const timer = setTimeout(() => {
      fetchSuggestions(searchTerm.trim())
    }, 300)
    return () => clearTimeout(timer)
  }, [searchTerm, selectedCountry, selectedPlatform])

  useEffect(() => {
    if (selectedPlatform === 'ios') fetchCachedKeywords()
  }, [cacheCountry, cacheGenre, selectedPlatform])

  const fetchSuggestions = async (term: string) => {
    setIsSearching(true)
    try {
      const data = await getAppStoreSuggestions(term, selectedCountry, selectedPlatform)
      setSuggestions(data)
      setHasSearched(true)
    } catch (err) {
      console.error('Error fetching suggestions:', err)
    } finally {
      setIsSearching(false)
    }
  }

  const fetchCachedKeywords = async () => {
    setIsCacheLoading(true)
    try {
      const data = await getKeywordCache(cacheCountry, cacheGenre, 100)
      setCachedKeywords(data)
    } catch (err) {
      console.error('Error fetching keyword cache:', err)
    } finally {
      setIsCacheLoading(false)
    }
  }

  const selectedCountryData = COUNTRIES.find(c => c.code === selectedCountry)

  return (
    <div className="space-y-4">
      {/* Platform toggle */}
      <div className="flex gap-2">
        {PLATFORMS.map(p => (
          <button
            key={p.value}
            onClick={() => { setSelectedPlatform(p.value); setSearchTerm(''); setSuggestions([]); setHasSearched(false) }}
            className={`px-5 py-2 rounded-full text-sm font-semibold transition-colors border ${
              selectedPlatform === p.value
                ? 'bg-blue-600 text-white border-blue-600'
                : 'bg-white dark:bg-[#12161e] text-gray-600 dark:text-gray-400 border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'
            }`}
          >
            {p.label}
          </button>
        ))}
      </div>

      {/* Search bar */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6">
        <div className="flex gap-3">
          <div className="flex-1">
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder={selectedPlatform === 'android' ? 'キーワードを入力（例: ゲーム, 写真, 天気...）' : 'キーワードを入力（例: photo, game, music...）'}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 text-base"
              autoFocus
            />
            <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
              {selectedPlatform === 'android' ? 'Google Play の検索サジェストをリアルタイムで取得します' : 'Apple App Store の検索サジェストをリアルタイムで取得します'}
            </p>
          </div>
          <div>
            <select
              value={selectedCountry}
              onChange={(e) => setSelectedCountry(e.target.value)}
              className="h-[50px] px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {COUNTRIES.map(c => (
                <option key={c.code} value={c.code}>{c.flag} {c.name}</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Suggestion Results */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
        <div className="px-6 py-4 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14] flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
            {searchTerm ? `「${searchTerm}」のサジェスト — ${selectedCountryData?.flag} ${selectedCountryData?.name}` : 'キーワードを入力して検索'}
          </h2>
          {hasSearched && <span className="text-sm text-gray-500 dark:text-gray-400">{suggestions.length}件</span>}
        </div>

        {isSearching ? (
          <div className="flex items-center justify-center py-12 text-gray-500 dark:text-gray-400">検索中...</div>
        ) : !searchTerm.trim() ? (
          <div className="p-6">
            <p className="text-xs text-gray-400 dark:text-gray-500 mb-3">よく検索されるワード</p>
            <div className="flex flex-wrap gap-2">
              {SEED_KEYWORDS.map(kw => (
                <button
                  key={kw}
                  onClick={() => setSearchTerm(kw)}
                  className="px-3 py-1.5 bg-gray-100 dark:bg-gray-800 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-700 dark:hover:text-blue-400 text-gray-700 dark:text-gray-300 rounded-full text-sm transition-colors"
                >
                  {kw}
                </button>
              ))}
            </div>
          </div>
        ) : suggestions.length === 0 && hasSearched ? (
          <div className="flex items-center justify-center py-12 text-gray-500 dark:text-gray-400">サジェストが見つかりませんでした</div>
        ) : (
          <>
            <ul className="divide-y divide-gray-100 dark:divide-gray-800">
              {suggestions.map((s, i) => (
                <li
                  key={s.term}
                  className="flex items-center gap-4 px-6 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer"
                  onClick={() => setSearchTerm(s.term)}
                >
                  <span className="text-gray-400 dark:text-gray-500 text-sm w-6 text-right">{i + 1}</span>
                  <span className="text-gray-800 dark:text-gray-200">{s.term}</span>
                </li>
              ))}
            </ul>
            <div className="border-t dark:border-gray-700 px-6 py-4 bg-blue-50 dark:bg-blue-900/10 flex items-center justify-between gap-4">
              <p className="text-sm text-blue-800 dark:text-blue-300">
                <span className="font-semibold">「{searchTerm}」</span>の順位を毎日自動追跡できます
              </p>
              <Link
                href="/signup"
                className="flex-shrink-0 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-colors"
              >
                無料で追跡する →
              </Link>
            </div>
          </>
        )}
      </div>

      {/* ASA Popularity Cache Table — iOS only */}
      {selectedPlatform === 'android' && (
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 text-center text-gray-500 dark:text-gray-400 text-sm">
          Google Play のキーワード人気スコアは現在準備中です。上の検索ボックスでサジェストを確認できます。
        </div>
      )}
      {selectedPlatform === 'ios' && <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
        <div className="px-6 py-4 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14]">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">ASA 人気キーワード</h2>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Apple Search Ads の人気スコア（0–100）</p>
            </div>
            <div className="flex items-center gap-2">
              <select
                value={cacheCountry}
                onChange={(e) => setCacheCountry(e.target.value)}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {COUNTRIES.map(c => (
                  <option key={c.code} value={c.code}>{c.flag} {c.name}</option>
                ))}
              </select>
              <select
                value={cacheGenre}
                onChange={(e) => setCacheGenre(e.target.value)}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {ASA_GENRES.map(g => (
                  <option key={g.value} value={g.value}>{g.label}</option>
                ))}
              </select>
              {isCacheLoading && <div className="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin" />}
            </div>
          </div>
        </div>

        {isCacheLoading && cachedKeywords.length === 0 ? (
          <div className="flex items-center justify-center py-12 text-gray-500 dark:text-gray-400">読み込み中...</div>
        ) : cachedKeywords.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 gap-3 text-center px-6">
            <div className="text-4xl">🔍</div>
            <p className="text-gray-700 dark:text-gray-300 font-medium">まだデータが準備中です</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">上の検索ボックスでキーワードを入力すると、リアルタイムでサジェストを確認できます</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50 dark:bg-[#0d0f14] border-b dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  <th className="px-4 py-3 w-8">#</th>
                  <th className="px-4 py-3">キーワード</th>
                  <th className="px-4 py-3">ジャンル</th>
                  <th className="px-4 py-3 min-w-[180px]">人気スコア</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-gray-800">
                {cachedKeywords.map((kw, i) => (
                  <tr key={`${kw.keyword}-${kw.genre}`} className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                    <td className="px-4 py-2 text-gray-400 dark:text-gray-500 text-xs">{i + 1}</td>
                    <td className="px-4 py-2">
                      <button
                        onClick={() => setSearchTerm(kw.keyword)}
                        className="font-medium text-gray-800 dark:text-gray-200 hover:text-blue-600 dark:hover:text-blue-400"
                      >
                        {kw.keyword}
                      </button>
                    </td>
                    <td className="px-4 py-2 text-gray-500 dark:text-gray-400">
                      {ASA_GENRES.find(g => g.value === kw.genre)?.label ?? kw.genre}
                    </td>
                    <td className="px-4 py-2">
                      <PopularityBar value={kw.popularity} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>}
    </div>
  )
}

const ANDROID_RANKING_TYPES = [
  { value: 'topfreeapplications', label: '無料トップ' },
  { value: 'toppaidapplications', label: '有料トップ' },
  { value: 'topgrossingapplications', label: '収益トップ' },
]

function AppRankingSection() {
  const [rankings, setRankings] = useState<AppRankingEntry[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [selectedPlatform, setSelectedPlatform] = useState('ios')
  const [selectedCountry, setSelectedCountry] = useState('jp')
  const [selectedRankingType, setSelectedRankingType] = useState('topfreeapplications')
  const [selectedGenre, setSelectedGenre] = useState('')

  const rankingTypes = selectedPlatform === 'android' ? ANDROID_RANKING_TYPES : RANKING_TYPES

  useEffect(() => {
    fetchRankings()
  }, [selectedCountry, selectedRankingType, selectedGenre, selectedPlatform])

  const fetchRankings = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await getAppRankings(selectedCountry, selectedRankingType, selectedPlatform === 'android' ? '' : selectedGenre, 100, selectedPlatform)
      setRankings(data)
    } catch (err) {
      setError('ランキングの取得に失敗しました')
      console.error('Error fetching app rankings:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const selectedCountryData = COUNTRIES.find(c => c.code === selectedCountry)
  const selectedRankingTypeData = rankingTypes.find(r => r.value === selectedRankingType)
  const selectedGenreData = GENRES.find(g => g.id === selectedGenre)

  return (
    <>
      {/* Platform toggle */}
      <div className="flex gap-2 mb-4">
        {PLATFORMS.map(p => (
          <button
            key={p.value}
            onClick={() => { setSelectedPlatform(p.value); setSelectedRankingType('topfreeapplications') }}
            className={`px-5 py-2 rounded-full text-sm font-semibold transition-colors border ${
              selectedPlatform === p.value
                ? 'bg-blue-600 text-white border-blue-600'
                : 'bg-white dark:bg-[#12161e] text-gray-600 dark:text-gray-400 border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'
            }`}
          >
            {p.label}
          </button>
        ))}
      </div>

      {/* Ranking type sub-tabs */}
      <div className="flex gap-2 mb-4 overflow-x-auto pb-1">
        {rankingTypes.map(type => (
          <button
            key={type.value}
            onClick={() => setSelectedRankingType(type.value)}
            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-colors ${
              selectedRankingType === type.value
                ? 'bg-blue-600 dark:bg-blue-500 text-white'
                : 'bg-white dark:bg-[#12161e] text-gray-600 dark:text-gray-400 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'
            }`}
          >
            {type.label}
          </button>
        ))}
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">国</label>
            <select
              value={selectedCountry}
              onChange={(e) => setSelectedCountry(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {COUNTRIES.map(country => (
                <option key={country.code} value={country.code}>
                  {country.flag} {country.name}
                </option>
              ))}
            </select>
          </div>

          {selectedPlatform === 'ios' && (
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">ジャンル</label>
            <select
              value={selectedGenre}
              onChange={(e) => setSelectedGenre(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {GENRES.map(genre => (
                <option key={genre.id} value={genre.id}>{genre.name}</option>
              ))}
            </select>
          </div>
          )}
        </div>
      </div>

      {/* Results */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
        <div className="px-6 py-4 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14]">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {selectedPlatform === 'android' ? '🤖 Google Play' : '🍎 App Store'} — {selectedCountryData?.flag} {selectedCountryData?.name} — {selectedRankingTypeData?.label}
              {selectedPlatform === 'ios' && selectedGenreData?.id && ` / ${selectedGenreData.name}`}
            </h2>
            <span className="text-sm text-gray-500 dark:text-gray-400">{rankings.length}件</span>
          </div>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500 dark:text-gray-400">読み込み中...</div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-red-500 dark:text-red-400">{error}</div>
          </div>
        ) : rankings.length === 0 ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-gray-500 dark:text-gray-400">ランキングが見つかりませんでした</div>
          </div>
        ) : (
          <div className="divide-y dark:divide-gray-700">
            {rankings.map((app) => (
              <div key={app.app_id} className="flex items-center gap-4 px-6 py-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                <span className="w-8 text-center text-sm font-bold text-gray-400 dark:text-gray-500">{app.rank}</span>
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
                  <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-xl flex-shrink-0" />
                )}
                <div className="flex-1 min-w-0">
                  <a
                    href={app.store_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-medium text-gray-900 dark:text-gray-100 hover:text-blue-600 dark:hover:text-blue-400 truncate block"
                  >
                    {app.name}
                  </a>
                  <p className="text-sm text-gray-500 dark:text-gray-400 truncate">{app.developer}</p>
                </div>
                <div className="text-right flex-shrink-0">
                  <p className="text-sm text-gray-600 dark:text-gray-400">{app.category}</p>
                  <p className="text-xs text-gray-400 dark:text-gray-500">{app.price || '無料'}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  )
}

function RankingTrendSection() {
  const [appId, setAppId] = useState('')
  const [country, setCountry] = useState('jp')
  const [rankingType, setRankingType] = useState('topfreeapplications')
  const [days, setDays] = useState(30)
  const [trendData, setTrendData] = useState<RankingTrendPoint[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searched, setSearched] = useState(false)

  const handleSearch = async () => {
    if (!appId.trim()) return
    setIsLoading(true)
    setError(null)
    setSearched(true)
    try {
      const data = await getAppRankingTrend(appId.trim(), country, rankingType, days)
      setTrendData(data)
    } catch (err) {
      setError('トレンドデータの取得に失敗しました')
      console.error('Error fetching ranking trend:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const countryData = COUNTRIES.find(c => c.code === country)
  const rankingTypeData = TREND_RANKING_TYPES.find(r => r.value === rankingType)

  return (
    <>
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
          <div className="md:col-span-1">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">App ID</label>
            <input
              type="text"
              value={appId}
              onChange={(e) => setAppId(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              placeholder="例: 389801252"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">国</label>
            <select
              value={country}
              onChange={(e) => setCountry(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {[
                { code: 'jp', name: '日本', flag: '🇯🇵' },
                { code: 'us', name: 'アメリカ', flag: '🇺🇸' },
                { code: 'kr', name: '韓国', flag: '🇰🇷' },
                { code: 'gb', name: 'イギリス', flag: '🇬🇧' },
              ].map(c => (
                <option key={c.code} value={c.code}>{c.flag} {c.name}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">ランキングタイプ</label>
            <select
              value={rankingType}
              onChange={(e) => setRankingType(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {TREND_RANKING_TYPES.map(r => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">期間</label>
            <select
              value={days}
              onChange={(e) => setDays(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value={7}>7日間</option>
              <option value={14}>14日間</option>
              <option value={30}>30日間</option>
            </select>
          </div>
        </div>
        <button
          onClick={handleSearch}
          disabled={!appId.trim() || isLoading}
          className="px-6 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {isLoading ? '検索中...' : '検索'}
        </button>
      </div>

      {searched && (
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
          <div className="px-6 py-4 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14]">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {countryData?.flag} {countryData?.name} — {rankingTypeData?.label} — 過去{days}日間
            </h2>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500 dark:text-gray-400">読み込み中...</div>
            </div>
          ) : error ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-red-500 dark:text-red-400">{error}</div>
            </div>
          ) : trendData.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500 dark:text-gray-400">データが見つかりませんでした</div>
            </div>
          ) : (
            <div className="p-6">
              <div className="h-64 mb-6">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={trendData} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="date" tick={{ fontSize: 12 }} />
                    <YAxis reversed={true} domain={['auto', 'auto']} tick={{ fontSize: 12 }} label={{ value: '順位', angle: -90, position: 'insideLeft' }} />
                    <Tooltip formatter={(value: number) => [`${value}位`, '順位']} />
                    <Line type="monotone" dataKey="rank" stroke="#2563eb" strokeWidth={2} dot={{ r: 3 }} activeDot={{ r: 5 }} />
                  </LineChart>
                </ResponsiveContainer>
              </div>

              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-gray-50 dark:bg-[#0d0f14] border-b dark:border-gray-700 text-left font-medium text-gray-500 dark:text-gray-400">
                      <th className="px-4 py-2">日付</th>
                      <th className="px-4 py-2 text-right">順位</th>
                    </tr>
                  </thead>
                  <tbody>
                    {[...trendData].reverse().map((point) => (
                      <tr key={point.date} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
                        <td className="px-4 py-2 text-gray-700 dark:text-gray-300">{point.date}</td>
                        <td className="px-4 py-2 text-right font-bold text-blue-600 dark:text-blue-400">{point.rank}位</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </>
  )
}

function CountryComparisonSection() {
  const [appId, setAppId] = useState('')
  const [rankingType, setRankingType] = useState('topfreeapplications')
  const [countryData, setCountryData] = useState<CountryRankPoint[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searched, setSearched] = useState(false)

  const handleSearch = async () => {
    if (!appId.trim()) return
    setIsLoading(true)
    setError(null)
    setSearched(true)
    try {
      const data = await getAppRankingCountries(appId.trim(), rankingType)
      setCountryData(data)
    } catch (err) {
      setError('国別ランキングの取得に失敗しました')
      console.error('Error fetching country rankings:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const rankingTypeData = TREND_RANKING_TYPES.find(r => r.value === rankingType)

  return (
    <>
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">App ID</label>
            <input
              type="text"
              value={appId}
              onChange={(e) => setAppId(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              placeholder="例: 389801252"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">ランキングタイプ</label>
            <select
              value={rankingType}
              onChange={(e) => setRankingType(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {TREND_RANKING_TYPES.map(r => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>
        </div>
        <button
          onClick={handleSearch}
          disabled={!appId.trim() || isLoading}
          className="px-6 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {isLoading ? '検索中...' : '検索'}
        </button>
      </div>

      {searched && (
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
          <div className="px-6 py-4 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14]">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              国別ランキング — {rankingTypeData?.label}
            </h2>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500 dark:text-gray-400">読み込み中...</div>
            </div>
          ) : error ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-red-500 dark:text-red-400">{error}</div>
            </div>
          ) : countryData.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500 dark:text-gray-400">データが見つかりませんでした</div>
            </div>
          ) : (
            <div className="p-6 grid grid-cols-2 md:grid-cols-4 gap-4">
              {countryData.map((point) => (
                <div
                  key={point.country}
                  className="bg-gray-50 dark:bg-[#0d0f14] rounded-lg p-4 text-center border dark:border-gray-700 hover:shadow-md transition-shadow"
                >
                  <div className="text-3xl mb-1">{COUNTRY_FLAGS[point.country] ?? '🌐'}</div>
                  <div className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
                    {COUNTRY_NAMES[point.country] ?? point.country.toUpperCase()}
                  </div>
                  <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">{point.rank}位</div>
                  <div className="text-xs text-gray-400 dark:text-gray-500 mt-1">
                    {new Date(point.fetched_at).toLocaleDateString('ja-JP')}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </>
  )
}

type Tab = 'keywords' | 'rankings' | 'trend' | 'countries'

const TABS: { id: Tab; label: string }[] = [
  { id: 'keywords', label: '人気キーワード' },
  { id: 'rankings', label: 'App Store ランキング' },
  { id: 'trend', label: 'ランキングトレンド' },
  { id: 'countries', label: '国別比較' },
]

function StickySignupBar() {
  const [visible, setVisible] = useState(false)
  const [dismissed, setDismissed] = useState(false)

  useEffect(() => {
    if (dismissed) return
    // Show after 20 seconds or after scroll
    const timer = setTimeout(() => setVisible(true), 20000)
    const onScroll = () => {
      if (window.scrollY > 400) { setVisible(true); window.removeEventListener('scroll', onScroll) }
    }
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => { clearTimeout(timer); window.removeEventListener('scroll', onScroll) }
  }, [dismissed])

  if (!visible || dismissed) return null

  return (
    <div className="fixed bottom-0 left-0 right-0 z-40 bg-white dark:bg-[#12161e] border-t border-gray-200 dark:border-gray-700 shadow-2xl dark:shadow-black/40 px-4 py-3">
      <div className="max-w-4xl mx-auto flex items-center justify-between gap-4">
        <div className="min-w-0">
          <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">自分のアプリのキーワード順位を毎日追跡しませんか？</p>
          <p className="text-xs text-gray-500 dark:text-gray-400 hidden sm:block">無料プランで1アプリ・10キーワードまでずっと無料</p>
        </div>
        <div className="flex items-center gap-2 flex-shrink-0">
          <Link
            href="/signup"
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-colors"
          >
            無料で始める
          </Link>
          <button
            onClick={() => setDismissed(true)}
            className="p-1.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
            aria-label="閉じる"
          >
            ✕
          </button>
        </div>
      </div>
    </div>
  )
}

export default function PopularKeywordsPage() {
  const [activeTab, setActiveTab] = useState<Tab>('keywords')

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-[#0d0f14]">
      {/* Header */}
      <div className="bg-white dark:bg-[#12161e] border-b dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">App Store データ</h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            人気キーワードとApp Storeランキングを確認する
          </p>
        </div>
      </div>

      {/* Main tabs */}
      <div className="bg-white dark:bg-[#12161e] border-b dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-0 overflow-x-auto">
            {TABS.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`px-6 py-4 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${
                  activeTab === tab.id
                    ? 'border-blue-600 text-blue-600 dark:text-blue-400'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 pb-12">
        {activeTab === 'keywords' && <KeywordsSection />}
        {activeTab === 'rankings' && <AppRankingSection />}
        {activeTab === 'trend' && <RankingTrendSection />}
        {activeTab === 'countries' && <CountryComparisonSection />}
      </div>

      {/* Signup CTA */}
      <div className="bg-gradient-to-br from-blue-600 to-blue-700 dark:from-blue-700 dark:to-blue-800">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16 text-center">
          <h2 className="text-2xl sm:text-3xl font-bold text-white mb-3">
            自分のアプリのASO管理もはじめませんか？
          </h2>
          <p className="text-blue-100 mb-8 text-base sm:text-lg">
            キーワード順位追跡・競合分析・レビュー管理をすべて一箇所で。無料で始められます。
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/signup"
              className="inline-block px-8 py-3 bg-white text-blue-700 font-semibold rounded-lg hover:bg-blue-50 transition-colors text-base"
            >
              無料アカウントを作成
            </Link>
            <Link
              href="/login"
              className="inline-block px-8 py-3 border border-white text-white font-semibold rounded-lg hover:bg-white/10 transition-colors text-base"
            >
              ログイン
            </Link>
          </div>
        </div>
      </div>

      <StickySignupBar />
    </div>
  )
}
