import { useState } from 'react'

export function useUpgradeModal() {
  const [isOpen, setIsOpen] = useState(false)
  const [triggerFeature, setTriggerFeature] = useState<string | undefined>()

  const open = (featureName?: string) => {
    setTriggerFeature(featureName)
    setIsOpen(true)
  }

  const close = () => setIsOpen(false)

  return { isOpen, triggerFeature, open, close }
}
