'use client'

import { useState, useEffect } from 'react'
import { App } from '@/lib/api'

export function useSelectedApp(apps: App[], isLoading: boolean) {
  const [selectedApp, setSelectedApp] = useState<App | null>(null)

  // Auto-select first app when apps are loaded
  useEffect(() => {
    if (!isLoading && apps.length > 0 && !selectedApp) {
      setSelectedApp(apps[0])
    }
  }, [apps, isLoading, selectedApp])

  return { selectedApp, setSelectedApp }
}
