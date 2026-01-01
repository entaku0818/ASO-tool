'use client'

import { useState, useEffect, useCallback, useMemo } from 'react'
import { Keyword, getKeywords, getLatestRanking } from '@/lib/api'

export type KeywordWithRanking = Keyword & {
  latestRank: number | null
}

export function useKeywords(appId: string) {
  const [keywords, setKeywords] = useState<KeywordWithRanking[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const fetchKeywords = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const keywordsData = await getKeywords(appId)

      const keywordsWithRankings = await Promise.all(
        keywordsData.map(async (keyword) => {
          const ranking = await getLatestRanking(appId, keyword.id)
          return {
            ...keyword,
            latestRank: ranking?.rank ?? null,
          }
        })
      )

      setKeywords(keywordsWithRankings)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch keywords'))
    } finally {
      setIsLoading(false)
    }
  }, [appId])

  useEffect(() => {
    fetchKeywords()
  }, [fetchKeywords])

  // Sort keywords by rank (null ranks at the end)
  const sortedKeywords = useMemo(() => {
    return [...keywords].sort((a, b) => {
      if (a.latestRank === null) return 1
      if (b.latestRank === null) return -1
      return a.latestRank - b.latestRank
    })
  }, [keywords])

  return { keywords: sortedKeywords, isLoading, error, refetch: fetchKeywords }
}
