'use client'

import { useState, useEffect } from 'react'
import { Ranking, getRankings } from '@/lib/api'

export function useRankings(keywordId: string) {
  const [rankings, setRankings] = useState<Ranking[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchRankings() {
      if (!keywordId) return

      try {
        setIsLoading(true)
        const data = await getRankings(keywordId)
        setRankings(data)
      } catch (e) {
        setError(e instanceof Error ? e : new Error('Failed to fetch rankings'))
      } finally {
        setIsLoading(false)
      }
    }

    fetchRankings()
  }, [keywordId])

  return { rankings, isLoading, error }
}
