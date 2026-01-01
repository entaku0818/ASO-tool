'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '@/contexts/AuthContext'

export function Header() {
  const { user, isLoading, logout } = useAuth()
  const pathname = usePathname()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && !user && pathname !== '/login') {
      router.push('/login')
    }
  }, [isLoading, user, pathname, router])

  // Don't show header on login page
  if (pathname === '/login') {
    return null
  }

  // Show loading state
  if (isLoading) {
    return (
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900">ASO Tool</span>
          <div className="text-gray-400">Loading...</div>
        </div>
      </header>
    )
  }

  // Don't render if not authenticated
  if (!user) {
    return null
  }

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
          ASO Tool
        </Link>
        <nav className="flex items-center gap-6">
          <Link href="/" className="text-gray-600 hover:text-gray-900">
            Dashboard
          </Link>
          <Link href="/keywords" className="text-gray-600 hover:text-gray-900">
            Keywords
          </Link>
          <div className="flex items-center gap-4 pl-4 border-l">
            <span className="text-sm text-gray-600">
              {user.name || user.email}
              {user.is_admin && <span className="ml-1 text-xs text-blue-600">(Admin)</span>}
            </span>
            <button
              onClick={handleLogout}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              Logout
            </button>
          </div>
        </nav>
      </div>
    </header>
  )
}
