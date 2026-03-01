'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { useApp } from '@/hooks/useApp'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'
import { useRankings } from '@/hooks/useRankings'
import { RankingChart } from '@/components/RankingChart'
import {
  createKeyword,
  deleteKeyword,
  getRisingKeywords,
  RisingKeyword,
  getSearchAdsCredentials,
  setSearchAdsCredentials,
  deleteSearchAdsCredentials,
  refreshKeywordPopularity,
  getKeywordSuggestions,
  SearchAdsCredentials,
  KeywordPopularitySuggestion,
} from '@/lib/api'
import { ReviewsSection } from '@/components/ReviewsSection'
import { CompetitorSection } from '@/components/CompetitorSection'
import { AnalyticsSection } from '@/components/AnalyticsSection'

function PopularityBar({ score }: { score?: number }) {
  if (score === undefined || score === null) {
    return <span className="text-gray-400">−</span>
  }
  const filled = score
  const empty = 5 - score
  return (
    <span className="inline-flex items-center gap-0.5" title={`人気スコア: ${score}/5`}>
      {Array.from({ length: filled }).map((_, i) => (
        <span key={`f${i}`} className="text-orange-400">●</span>
      ))}
      {Array.from({ length: empty }).map((_, i) => (
        <span key={`e${i}`} className="text-gray-300">●</span>
      ))}
    </span>
  )
}

function KeywordRow({
  keyword,
  isSelected,
  onSelect,
  onDelete,
}: {
  keyword: KeywordWithRanking
  isSelected: boolean
  onSelect: () => void
  onDelete: () => void
}) {
  const rankColor = keyword.latestRank === null
    ? 'text-gray-500'
    : keyword.latestRank <= 10
    ? 'text-green-600'
    : keyword.latestRank <= 50
    ? 'text-yellow-600'
    : 'text-red-600'

  return (
    <tr
      className={`border-b cursor-pointer ${isSelected ? 'bg-blue-50' : 'hover:bg-gray-50'}`}
      onClick={onSelect}
    >
      <td className="py-3 px-4">{keyword.keyword}</td>
      <td className="py-3 px-4">{keyword.country}</td>
      <td className={`py-3 px-4 font-bold ${rankColor}`}>
        {keyword.latestRank === null ? '圏外' : `${keyword.latestRank}位`}
      </td>
      <td className="py-3 px-4">
        <PopularityBar score={keyword.popularity_score} />
      </td>
      <td className="py-3 px-4">
        <button
          onClick={(e) => {
            e.stopPropagation()
            onDelete()
          }}
          className="text-red-500 hover:text-red-700 text-sm"
        >
          削除
        </button>
      </td>
    </tr>
  )
}

function RisingKeywordsSection({ appId }: { appId: string }) {
  const [keywords, setKeywords] = useState<RisingKeyword[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    getRisingKeywords(appId)
      .then(setKeywords)
      .catch(() => setKeywords([]))
      .finally(() => setIsLoading(false))
  }, [appId])

  if (isLoading || keywords.length === 0) return null

  return (
    <div className="bg-white rounded-lg shadow mt-6">
      <div className="p-4 border-b">
        <h3 className="text-lg font-semibold">急上昇キーワード</h3>
        <p className="text-sm text-gray-500">過去7日間で順位が上昇したキーワード</p>
      </div>
      <table className="w-full">
        <thead className="bg-gray-50">
          <tr>
            <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
            <th className="py-3 px-4 text-left font-medium text-gray-600">国</th>
            <th className="py-3 px-4 text-right font-medium text-gray-600">現在順位</th>
            <th className="py-3 px-4 text-right font-medium text-gray-600">以前の順位</th>
            <th className="py-3 px-4 text-right font-medium text-gray-600">上昇幅</th>
          </tr>
        </thead>
        <tbody>
          {keywords.map((kw) => (
            <tr key={kw.keyword_id} className="border-b hover:bg-gray-50">
              <td className="py-3 px-4 font-medium text-gray-900">{kw.keyword}</td>
              <td className="py-3 px-4 text-gray-500">{kw.country}</td>
              <td className="py-3 px-4 text-right font-bold text-green-600">{kw.current_rank}位</td>
              <td className="py-3 px-4 text-right text-gray-500">{kw.previous_rank}位</td>
              <td className="py-3 px-4 text-right">
                <span className="inline-flex items-center gap-1 text-green-600 font-bold">
                  ▲ {kw.improvement}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function RankingChartSection({ appId, keywordId, keywordName }: { appId: string; keywordId: string; keywordName: string }) {
  const { rankings, isLoading } = useRankings(appId, keywordId)

  return (
    <div className="bg-white rounded-lg shadow p-4 mb-6">
      <h3 className="text-lg font-semibold mb-4">「{keywordName}」の順位推移</h3>
      {isLoading ? (
        <p className="text-gray-500">読み込み中...</p>
      ) : (
        <RankingChart rankings={rankings} />
      )}
    </div>
  )
}

function SearchAdsSection({ appId, onPopularityRefreshed }: { appId: string; onPopularityRefreshed: () => void }) {
  const [isOpen, setIsOpen] = useState(false)
  const [credentials, setCredentials] = useState<SearchAdsCredentials | null>(null)
  const [isLoadingCreds, setIsLoadingCreds] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [isRefreshing, setIsRefreshing] = useState(false)

  // Form state
  const [clientId, setClientId] = useState('')
  const [teamId, setTeamId] = useState('')
  const [keyId, setKeyId] = useState('')
  const [privateKey, setPrivateKey] = useState('')
  const [orgId, setOrgId] = useState('')
  const [adamId, setAdamId] = useState('')

  useEffect(() => {
    if (!isOpen) return
    setIsLoadingCreds(true)
    getSearchAdsCredentials(appId)
      .then((creds) => {
        setCredentials(creds)
        if (creds) {
          setClientId(creds.client_id)
          setTeamId(creds.team_id)
          setKeyId(creds.key_id)
          setOrgId(creds.org_id ?? '')
          setAdamId(String(creds.adam_id))
        }
      })
      .finally(() => setIsLoadingCreds(false))
  }, [appId, isOpen])

  const handleSave = async () => {
    setIsSaving(true)
    try {
      await setSearchAdsCredentials(appId, {
        client_id: clientId,
        team_id: teamId,
        key_id: keyId,
        private_key: privateKey,
        org_id: orgId || undefined,
        adam_id: parseInt(adamId, 10),
      })
      const updated = await getSearchAdsCredentials(appId)
      setCredentials(updated)
      setPrivateKey('')
      alert('認証情報を保存しました')
    } catch {
      alert('保存に失敗しました')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('Search Ads認証情報を削除しますか？')) return
    setIsDeleting(true)
    try {
      await deleteSearchAdsCredentials(appId)
      setCredentials(null)
      setClientId('')
      setTeamId('')
      setKeyId('')
      setPrivateKey('')
      setOrgId('')
      setAdamId('')
    } catch {
      alert('削除に失敗しました')
    } finally {
      setIsDeleting(false)
    }
  }

  const handleRefresh = async () => {
    setIsRefreshing(true)
    try {
      const result = await refreshKeywordPopularity(appId)
      alert(`${result.updated}件のキーワードのスコアを更新しました`)
      onPopularityRefreshed()
    } catch {
      alert('人気スコアの更新に失敗しました')
    } finally {
      setIsRefreshing(false)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow mb-6">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full p-4 flex items-center justify-between text-left"
      >
        <div>
          <h3 className="text-lg font-semibold">Search Ads 設定</h3>
          <p className="text-sm text-gray-500">キーワード人気スコア取得のための認証情報</p>
        </div>
        <span className="text-gray-400">{isOpen ? '▲' : '▼'}</span>
      </button>

      {isOpen && (
        <div className="border-t p-4">
          {isLoadingCreds ? (
            <p className="text-gray-500">読み込み中...</p>
          ) : (
            <>
              {credentials && (
                <div className="mb-4 p-3 bg-green-50 rounded-lg flex items-center justify-between">
                  <div>
                    <p className="text-sm text-green-800 font-medium">認証情報が設定されています</p>
                    <p className="text-xs text-green-600">Client ID: {credentials.client_id}</p>
                    <p className="text-xs text-green-600">Adam ID: {credentials.adam_id}</p>
                  </div>
                  <button
                    onClick={handleRefresh}
                    disabled={isRefreshing}
                    className="px-3 py-1.5 bg-orange-500 text-white rounded text-sm hover:bg-orange-600 disabled:opacity-50"
                  >
                    {isRefreshing ? '更新中...' : '人気スコア更新'}
                  </button>
                </div>
              )}

              <div className="grid grid-cols-2 gap-3 mb-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Client ID</label>
                  <input
                    type="text"
                    value={clientId}
                    onChange={(e) => setClientId(e.target.value)}
                    placeholder="APPLE_CLIENT_ID"
                    className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Team ID</label>
                  <input
                    type="text"
                    value={teamId}
                    onChange={(e) => setTeamId(e.target.value)}
                    placeholder="APPLE_TEAM_ID"
                    className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Key ID</label>
                  <input
                    type="text"
                    value={keyId}
                    onChange={(e) => setKeyId(e.target.value)}
                    placeholder="SEARCHADS_KEY_ID"
                    className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Adam ID（iTunes数値ID）</label>
                  <input
                    type="number"
                    value={adamId}
                    onChange={(e) => setAdamId(e.target.value)}
                    placeholder="123456789"
                    className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Org ID（任意）</label>
                  <input
                    type="text"
                    value={orgId}
                    onChange={(e) => setOrgId(e.target.value)}
                    placeholder="組織ID"
                    className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
              <div className="mb-3">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  秘密鍵 (.p8 ファイルの内容をBase64エンコード)
                  {credentials && <span className="text-gray-400 ml-1">（変更する場合のみ入力）</span>}
                </label>
                <textarea
                  value={privateKey}
                  onChange={(e) => setPrivateKey(e.target.value)}
                  placeholder="Base64エンコードされた .p8 ファイルの内容"
                  rows={3}
                  className="w-full px-3 py-2 border rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleSave}
                  disabled={isSaving}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50"
                >
                  {isSaving ? '保存中...' : '保存'}
                </button>
                {credentials && (
                  <button
                    onClick={handleDelete}
                    disabled={isDeleting}
                    className="px-4 py-2 bg-red-500 text-white rounded-lg text-sm hover:bg-red-600 disabled:opacity-50"
                  >
                    {isDeleting ? '削除中...' : '削除'}
                  </button>
                )}
              </div>
            </>
          )}
        </div>
      )}
    </div>
  )
}

function KeywordSuggestionsSection({ appId, onAddKeyword }: { appId: string; onAddKeyword: (keyword: string) => void }) {
  const [suggestions, setSuggestions] = useState<KeywordPopularitySuggestion[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [hasLoaded, setHasLoaded] = useState(false)

  const handleFetch = async () => {
    setIsLoading(true)
    try {
      const result = await getKeywordSuggestions(appId, 25)
      setSuggestions(result)
      setHasLoaded(true)
    } catch {
      alert('キーワード提案の取得に失敗しました。Search Ads認証情報を確認してください。')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow mt-6">
      <div className="p-4 border-b flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">キーワード提案</h3>
          <p className="text-sm text-gray-500">Apple Search Ads によるキーワード候補（人気スコア付き）</p>
        </div>
        <button
          onClick={handleFetch}
          disabled={isLoading}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50"
        >
          {isLoading ? '取得中...' : '提案を取得'}
        </button>
      </div>

      {hasLoaded && (
        suggestions.length === 0 ? (
          <p className="p-4 text-gray-500">提案が見つかりませんでした</p>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">人気スコア</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600"></th>
              </tr>
            </thead>
            <tbody>
              {suggestions
                .sort((a, b) => b.popularityScore - a.popularityScore)
                .map((s) => (
                  <tr key={s.text} className="border-b hover:bg-gray-50">
                    <td className="py-3 px-4 font-medium text-gray-900">{s.text}</td>
                    <td className="py-3 px-4">
                      <PopularityBar score={s.popularityScore} />
                    </td>
                    <td className="py-3 px-4">
                      <button
                        onClick={() => onAddKeyword(s.text)}
                        className="text-blue-600 hover:text-blue-800 text-sm"
                      >
                        追加
                      </button>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        )
      )}
    </div>
  )
}

export default function AppDetailPage() {
  const params = useParams()
  const appId = params.id as string
  const { app, isLoading: appLoading, error: appError } = useApp(appId)
  const { keywords, isLoading: keywordsLoading, refetch } = useKeywords(appId)
  const [selectedKeyword, setSelectedKeyword] = useState<KeywordWithRanking | null>(null)
  const [newKeyword, setNewKeyword] = useState('')
  const [isAdding, setIsAdding] = useState(false)

  const handleAddKeyword = async (kw?: string) => {
    const text = kw ?? newKeyword.trim()
    if (!text) return
    setIsAdding(true)
    try {
      await createKeyword(appId, text)
      if (!kw) setNewKeyword('')
      refetch()
    } catch {
      alert('キーワードの追加に失敗しました')
    } finally {
      setIsAdding(false)
    }
  }

  const handleDeleteKeyword = async (keywordId: string) => {
    if (!confirm('このキーワードを削除しますか？')) return
    try {
      await deleteKeyword(appId, keywordId)
      if (selectedKeyword?.id === keywordId) {
        setSelectedKeyword(null)
      }
      refetch()
    } catch {
      alert('キーワードの削除に失敗しました')
    }
  }

  if (appLoading) {
    return (
      <div>
        <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; ダッシュボードに戻る
        </Link>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  if (appError || !app) {
    return (
      <div>
        <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; ダッシュボードに戻る
        </Link>
        <p className="text-red-600">アプリが見つかりません</p>
      </div>
    )
  }

  const sortedKeywords = [...keywords].sort((a, b) => {
    if (a.latestRank === null) return 1
    if (b.latestRank === null) return -1
    return a.latestRank - b.latestRank
  })

  return (
    <div>
      <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
        &larr; ダッシュボードに戻る
      </Link>

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold">{app.name}</h2>
            <p className="text-gray-500">{app.bundle_id}</p>
          </div>
          <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
            {app.platform.toUpperCase()}
          </span>
        </div>
        {app.store_url && (
          <a
            href={app.store_url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline text-sm mt-2 inline-block"
          >
            ストアページを開く
          </a>
        )}
      </div>

      {app.platform === 'ios' && (
        <div className="mb-6">
          <AnalyticsSection appId={appId} platform={app.platform} />
        </div>
      )}

      {app.platform === 'ios' && (
        <SearchAdsSection appId={appId} onPopularityRefreshed={refetch} />
      )}

      {selectedKeyword && (
        <RankingChartSection
          appId={appId}
          keywordId={selectedKeyword.id}
          keywordName={selectedKeyword.keyword}
        />
      )}

      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold">キーワード順位 ({keywords.length}件)</h3>
            <p className="text-sm text-gray-500">クリックで順位推移を表示</p>
          </div>
        </div>

        <div className="p-4 border-b bg-gray-50">
          <div className="flex gap-2">
            <input
              type="text"
              value={newKeyword}
              onChange={(e) => setNewKeyword(e.target.value)}
              placeholder="新しいキーワード"
              className="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              onKeyDown={(e) => e.key === 'Enter' && handleAddKeyword()}
            />
            <button
              onClick={() => handleAddKeyword()}
              disabled={isAdding || !newKeyword.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isAdding ? '追加中...' : '追加'}
            </button>
          </div>
        </div>

        {keywordsLoading ? (
          <p className="p-4 text-gray-500">読み込み中...</p>
        ) : keywords.length === 0 ? (
          <p className="p-4 text-gray-500">キーワードが登録されていません</p>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">国</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">順位</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">人気</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600"></th>
              </tr>
            </thead>
            <tbody>
              {sortedKeywords.map((keyword) => (
                <KeywordRow
                  key={keyword.id}
                  keyword={keyword}
                  isSelected={selectedKeyword?.id === keyword.id}
                  onSelect={() => setSelectedKeyword(keyword)}
                  onDelete={() => handleDeleteKeyword(keyword.id)}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>

      <RisingKeywordsSection appId={appId} />

      {app.platform === 'ios' && (
        <KeywordSuggestionsSection appId={appId} onAddKeyword={(kw) => handleAddKeyword(kw)} />
      )}

      <div className="mt-6">
        <CompetitorSection
          appId={appId}
          platform={app.platform}
          selectedKeywordId={selectedKeyword?.id}
        />
      </div>

      <div className="mt-6">
        <ReviewsSection appId={appId} />
      </div>
    </div>
  )
}
