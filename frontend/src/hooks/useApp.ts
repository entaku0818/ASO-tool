'use client'

import { useState, useEffect } from 'react'
import { App, getApp } from '@/lib/api'

export function useApp(id: string) {
  const [app, setApp] = useState<App | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchApp() {
      if (!id) return

      try {
        setIsLoading(true)
        const data = await getApp(id)
        setApp(data)
      } catch (e) {
        setError(e instanceof Error ? e : new Error('Failed to fetch app'))
      } finally {
        setIsLoading(false)
      }
    }

    fetchApp()
  }, [id])

  return { app, isLoading, error }
}
