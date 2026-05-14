'use client'

import { useState, useEffect } from 'react'
import { getRankings } from '@/lib/api'
import { KeywordRankingData } from '@/components/RankingChart'

export function useMultiRankings(
  appId: string,
  keywords: { id: string; keyword: string }[],
  days: number
) {
  const [data, setData] = useState<KeywordRankingData[]>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!appId || keywords.length === 0) {
      setData([])
      return
    }

    let cancelled = false
    async function fetch() {
      setIsLoading(true)
      const results = await Promise.allSettled(
        keywords.map((kw) => getRankings(appId, kw.id, days).then((r) => ({ ...kw, rankings: r })))
      )
      if (!cancelled) {
        setData(
          results
            .filter((r): r is PromiseFulfilledResult<KeywordRankingData> => r.status === 'fulfilled')
            .map((r) => r.value)
        )
        setIsLoading(false)
      }
    }

    fetch()
    return () => { cancelled = true }
  }, [appId, keywords.map((k) => k.id).join(','), days]) // eslint-disable-line react-hooks/exhaustive-deps

  return { data, isLoading }
}
