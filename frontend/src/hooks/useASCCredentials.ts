'use client'

import { useState, useEffect, useCallback } from 'react'
import { ASCCredentials, getASCCredentials } from '@/lib/api'

export function useASCCredentials(appId: string) {
  const [credentials, setCredentials] = useState<ASCCredentials | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const refetch = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const data = await getASCCredentials(appId)
      setCredentials(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch credentials'))
    } finally {
      setIsLoading(false)
    }
  }, [appId])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { credentials, isLoading, error, refetch }
}
