'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/contexts/ThemeContext'

function ThemeToggle() {
  const { theme, toggle } = useTheme()
  return (
    <button
      onClick={toggle}
      aria-label="ダークモード切り替え"
      className="p-2 rounded-lg text-gray-500 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-white/10 transition-colors"
    >
      {theme === 'dark' ? (
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/>
        </svg>
      ) : (
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/>
        </svg>
      )}
    </button>
  )
}

export function Header() {
  const { user, isLoading, logout } = useAuth()
  const pathname = usePathname()
  const router = useRouter()

  const publicPages = ['/', '/login', '/popular-keywords', '/buy', '/buy/success', '/terms', '/privacy', '/tokushoho']
  const isPublicPage = publicPages.some(p => pathname === p || pathname.startsWith(p + '/'))

  useEffect(() => {
    if (!isLoading && !user && !isPublicPage) {
      router.push('/login')
    }
  }, [isLoading, user, pathname, router, isPublicPage])

  if (pathname === '/login') {
    return null
  }

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  if (isLoading && !isPublicPage) {
    return (
      <header className="bg-white dark:bg-[#12161e] shadow-sm border-b dark:border-white/10">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900 dark:text-gray-100">ASO Compass</span>
          <div className="text-gray-400 dark:text-gray-500">Loading...</div>
        </div>
      </header>
    )
  }

  if (isPublicPage || !user) {
    return (
      <header className="bg-white dark:bg-[#12161e] shadow-sm border-b dark:border-white/10">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <Link href="/" className="text-xl font-bold text-gray-900 dark:text-gray-100 hover:text-gray-700 dark:hover:text-white">
            ASO Compass
          </Link>
          <nav className="flex items-center gap-4">
            <Link href="/popular-keywords" className="hidden sm:block text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
              人気キーワード
            </Link>
            {user ? (
              <>
                <Link href="/dashboard" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
                  Dashboard
                </Link>
                <Link href="/keywords" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
                  Keywords
                </Link>
                <div className="flex items-center gap-4 pl-4 border-l dark:border-white/10">
                  <span className="text-sm text-gray-600 dark:text-gray-400">
                    {user.name || user.email}
                    {user.is_admin && <span className="ml-1 text-xs text-blue-600 dark:text-blue-400">(Admin)</span>}
                  </span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                  >
                    Logout
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link href="/login" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
                  Login
                </Link>
                <Link href="/popular-keywords" className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-semibold transition-colors">
                  無料で試す
                </Link>
              </>
            )}
            <ThemeToggle />
          </nav>
        </div>
      </header>
    )
  }

  return (
    <header className="bg-white dark:bg-[#12161e] shadow-sm border-b dark:border-white/10">
      <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="text-xl font-bold text-gray-900 dark:text-gray-100 hover:text-gray-700 dark:hover:text-white">
          ASO Compass
        </Link>
        <nav className="flex items-center gap-4">
          <Link href="/dashboard" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
            Dashboard
          </Link>
          <Link href="/keywords" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
            Keywords
          </Link>
          <Link href="/templates" className={`hover:text-gray-900 dark:hover:text-gray-100 ${pathname === '/templates' ? 'text-blue-600 dark:text-blue-400 font-medium' : 'text-gray-600 dark:text-gray-400'}`}>
            テンプレート
          </Link>
          <Link href="/popular-keywords" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
            人気キーワード
          </Link>
          <div className="flex items-center gap-4 pl-4 border-l dark:border-white/10">
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {user.name || user.email}
              {user.is_admin && <span className="ml-1 text-xs text-blue-600 dark:text-blue-400">(Admin)</span>}
            </span>
            <button
              onClick={handleLogout}
              className="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
            >
              Logout
            </button>
          </div>
          <ThemeToggle />
        </nav>
      </div>
    </header>
  )
}
