'use client'

import { useState, useCallback } from 'react'
import { Ranking, getRankings } from '@/lib/api'

export function useRankingHistory() {
  const [rankings, setRankings] = useState<Ranking[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  const fetchHistory = useCallback(async (appId: string, keywordId: string) => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await getRankings(appId, keywordId)
      // Sort by date ascending for chart display
      const sorted = [...data].sort((a, b) =>
        new Date(a.recorded_at).getTime() - new Date(b.recorded_at).getTime()
      )
      setRankings(sorted)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch rankings'))
    } finally {
      setIsLoading(false)
    }
  }, [])

  const clearHistory = useCallback(() => {
    setRankings([])
  }, [])

  return { rankings, isLoading, error, fetchHistory, clearHistory }
}
