'use client'

import { useState, useEffect, useCallback } from 'react'
import { Screenshot, DeviceType, getScreenshots } from '@/lib/api'

export function useScreenshots(appId: string, deviceType?: DeviceType) {
  const [screenshots, setScreenshots] = useState<Screenshot[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const fetchScreenshots = useCallback(async () => {
    if (!appId) return

    try {
      setIsLoading(true)
      const data = await getScreenshots(appId, deviceType)
      setScreenshots(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch screenshots'))
    } finally {
      setIsLoading(false)
    }
  }, [appId, deviceType])

  useEffect(() => {
    fetchScreenshots()
  }, [fetchScreenshots])

  const refetch = useCallback(() => {
    fetchScreenshots()
  }, [fetchScreenshots])

  return { screenshots, isLoading, error, refetch }
}
