'use client'

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react'
import { wakeDatabase } from '@/lib/api'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://aso-api-671942133800.asia-northeast1.run.app'

export type User = {
  id: string
  email: string
  name: string
  is_admin: boolean
  plan: 'free' | 'pro'
  plan_expires_at?: string
  created_at: string
}

type AuthContextType = {
  user: User | null
  token: string | null
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const logout = useCallback(() => {
    localStorage.removeItem('token')
    setToken(null)
    setUser(null)
  }, [])

  const fetchUser = useCallback(async (authToken: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      })

      if (!response.ok) {
        throw new Error('Failed to fetch user')
      }

      const userData = await response.json()
      setUser(userData)
      return true
    } catch {
      logout()
      return false
    }
  }, [logout])

  useEffect(() => {
    // Wake up database on app load
    wakeDatabase()

    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setToken(storedToken)
      fetchUser(storedToken).finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [fetchUser])

  const login = async (email: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Login failed')
    }

    const data = await response.json()
    localStorage.setItem('token', data.token)
    setToken(data.token)
    setUser(data.user)
  }

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
