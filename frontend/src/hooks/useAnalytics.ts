'use client'

import { useState, useEffect, useCallback } from 'react'
import {
  Analytics,
  AnalyticsWithCorrelation,
  AnalyticsSummary,
  getAnalytics,
  getAnalyticsWithCorrelation,
  getAnalyticsSummary,
} from '@/lib/api'

export function useAnalytics(appId: string, days: number = 30) {
  const [analytics, setAnalytics] = useState<Analytics[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const refetch = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const data = await getAnalytics(appId, days)
      setAnalytics(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch analytics'))
    } finally {
      setIsLoading(false)
    }
  }, [appId, days])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { analytics, isLoading, error, refetch }
}

export function useAnalyticsWithCorrelation(appId: string, days: number = 30) {
  const [analytics, setAnalytics] = useState<AnalyticsWithCorrelation[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const refetch = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const data = await getAnalyticsWithCorrelation(appId, days)
      setAnalytics(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch analytics'))
    } finally {
      setIsLoading(false)
    }
  }, [appId, days])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { analytics, isLoading, error, refetch }
}

export function useAnalyticsSummary(appId: string, days: number = 30) {
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const refetch = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const data = await getAnalyticsSummary(appId, days)
      setSummary(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch summary'))
    } finally {
      setIsLoading(false)
    }
  }, [appId, days])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { summary, isLoading, error, refetch }
}
