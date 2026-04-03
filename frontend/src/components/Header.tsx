'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '@/contexts/AuthContext'

export function Header() {
  const { user, isLoading, logout } = useAuth()
  const pathname = usePathname()
  const router = useRouter()

  // Public pages that don't require authentication
  const publicPages = ['/login', '/popular-keywords']
  const isPublicPage = publicPages.includes(pathname)

  useEffect(() => {
    if (!isLoading && !user && !isPublicPage) {
      router.push('/login')
    }
  }, [isLoading, user, pathname, router, isPublicPage])

  // Don't show header on login page
  if (pathname === '/login') {
    return null
  }

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  // Show loading state
  if (isLoading && !isPublicPage) {
    return (
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900">ASO Tool</span>
          <div className="text-gray-400">Loading...</div>
        </div>
      </header>
    )
  }

  // Public page header (no auth required)
  if (isPublicPage || !user) {
    return (
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <Link href="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
            ASO Tool
          </Link>
          <nav className="flex items-center gap-6">
            <Link href="/popular-keywords" className="text-gray-600 hover:text-gray-900">
              人気キーワード
            </Link>
            {user ? (
              <>
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
              </>
            ) : (
              <Link href="/login" className="text-gray-600 hover:text-gray-900">
                Login
              </Link>
            )}
          </nav>
        </div>
      </header>
    )
  }

  // Authenticated user header
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
          <Link href="/templates" className={`hover:text-gray-900 ${pathname === '/templates' ? 'text-blue-600 font-medium' : 'text-gray-600'}`}>
            テンプレート
          </Link>
          <Link href="/popular-keywords" className="text-gray-600 hover:text-gray-900">
            人気キーワード
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
