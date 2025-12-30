'use client'

import { useState, useEffect } from 'react'
import { Review, ReviewStats, getReviews, getReviewStats } from '@/lib/api'

export function useReviews(appId: string) {
  const [reviews, setReviews] = useState<Review[]>([])
  const [stats, setStats] = useState<ReviewStats | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchReviews() {
      if (!appId) return

      try {
        setIsLoading(true)
        const [reviewsData, statsData] = await Promise.all([
          getReviews(appId),
          getReviewStats(appId).catch(() => null),
        ])
        setReviews(reviewsData)
        setStats(statsData)
      } catch (e) {
        setError(e instanceof Error ? e : new Error('Failed to fetch reviews'))
      } finally {
        setIsLoading(false)
      }
    }

    fetchReviews()
  }, [appId])

  return { reviews, stats, isLoading, error }
}
