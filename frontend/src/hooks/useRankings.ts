'use client'

import { useState, useEffect } from 'react'
import { Ranking, getRankings } from '@/lib/api'

export function useRankings(appId: string, keywordId: string, days: number = 30) {
  const [rankings, setRankings] = useState<Ranking[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    if (!appId || !keywordId) return

    let cancelled = false
    async function fetchRankings() {
      try {
        setIsLoading(true)
        const data = await getRankings(appId, keywordId, days)
        if (!cancelled) setRankings(data)
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e : new Error('Failed to fetch rankings'))
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    }

    fetchRankings()
    return () => { cancelled = true }
  }, [appId, keywordId, days])

  return { rankings, isLoading, error }
}
