'use client'

import { useState, useEffect } from 'react'
import { App, getApps } from '@/lib/api'

export function useApps() {
  const [apps, setApps] = useState<App[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchApps() {
      try {
        setIsLoading(true)
        const data = await getApps()
        setApps(data)
      } catch (e) {
        setError(e instanceof Error ? e : new Error('Failed to fetch apps'))
      } finally {
        setIsLoading(false)
      }
    }

    fetchApps()
  }, [])

  const refetch = async () => {
    try {
      setIsLoading(true)
      const data = await getApps()
      setApps(data)
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch apps'))
    } finally {
      setIsLoading(false)
    }
  }

  return { apps, isLoading, error, refetch }
}
