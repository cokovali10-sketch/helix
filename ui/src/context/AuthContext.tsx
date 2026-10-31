import React, { createContext, useContext, useState, useEffect, useCallback } from 'react'

interface User {
  id: string
  name: string
  email: string
  avatar?: string
  role: 'admin' | 'user' | 'viewer'
  tenantId: string
  tenantName: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

// Mock user data for development
const mockUser: User = {
  id: 'user-1',
  name: 'Alex Chen',
  email: 'alex@helix.io',
  role: 'admin',
  tenantId: 'tenant-1',
  tenantName: 'SOC Team Alpha',
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const refreshUser = useCallback(async () => {
    setIsLoading(true)
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    setUser(mockUser)
    setIsLoading(false)
  }, [])

  useEffect(() => {
    refreshUser()
  }, [refreshUser])

  const login = useCallback(async (email: string, password: string) => {
    setIsLoading(true)
    await new Promise(resolve => setTimeout(resolve, 1000))
    // Simulate successful login
    setUser(mockUser)
    setIsLoading(false)
  }, [])

  const logout = useCallback(() => {
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{
      user,
      isAuthenticated: !!user,
      isLoading,
      login,
      logout,
      refreshUser,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}