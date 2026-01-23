'use client'

import { useState, useEffect } from 'react'
import { getPopularKeywords, PopularKeyword } from '@/lib/api'

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

export default function PopularKeywordsPage() {
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
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-3xl font-bold text-gray-900">人気キーワード</h1>
          <p className="mt-2 text-sm text-gray-600">
            国とプラットフォーム別の最も検索されたキーワードを探す
          </p>
        </div>
      </div>

      {/* Filters */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {/* Country Select */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                国
              </label>
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

            {/* Platform Select */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                プラットフォーム
              </label>
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

            {/* Limit Select */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                表示件数
              </label>
              <select
                value={selectedLimit}
                onChange={(e) => setSelectedLimit(Number(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {LIMITS.map(limit => (
                  <option key={limit} value={limit}>
                    {limit}件
                  </option>
                ))}
              </select>
            </div>

            {/* Search */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                検索
              </label>
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
      </div>

      {/* Results */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-12">
        <div className="bg-white rounded-lg shadow overflow-hidden">
          {/* Header */}
          <div className="px-6 py-4 border-b bg-gray-50">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">
                {selectedCountryData?.flag} {selectedCountryData?.name} - {selectedPlatform.toUpperCase()}
              </h2>
              <span className="text-sm text-gray-500">
                {filteredKeywords.length}件のキーワード
              </span>
            </div>
          </div>

          {/* Table */}
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
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {index + 1}
                      </td>
                      <td className="px-6 py-4">
                        <span className="font-medium text-gray-900">
                          {keyword.keyword}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-right text-sm text-gray-600">
                        {keyword.tracking_count}
                      </td>
                      <td className="px-6 py-4 text-right text-sm text-gray-600">
                        {keyword.results_count}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
