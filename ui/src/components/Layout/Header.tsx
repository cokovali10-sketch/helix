import React, { useState, useEffect } from 'react'
import { useTheme } from '../../../context/ThemeContext'
import { useAuth } from '../../../context/AuthContext'
import { useQuery } from '../../../context/QueryContext'
import { BellIcon, SearchIcon, MenuIcon, SunIcon, MoonIcon, UserIcon, LogoutIcon, ChevronDownIcon } from '../Layout/Sidebar'

interface HeaderProps {
  onMenuClick: () => void
}

export function Header({ onMenuClick }: { onMenuClick: () => void }) {
  const { theme, toggleTheme } = useTheme()
  const { user, logout } = useAuth()
  const { executeQuery, isExecuting, currentQuery, setCurrentQuery, currentResult, error, isExecuting: queryExecuting } = useQuery()
  const [queryInput, setQueryInput] = useState('')
  const [showUserMenu, setShowUserMenu] = useState(false)
  const [showNotifications, setShowNotifications] = useState(false)
  const [showCommandPalette, setShowCommandPalette] = useState(false)
  const searchRef = React.useRef<HTMLInputElement>(null)

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        setShowCommandPalette(true)
      }
      if (e.key === 'Escape') {
        setShowCommandPalette(false)
        setShowUserMenu(false)
        setShowNotifications(false)
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (queryInput.trim()) {
      executeQuery({ expression: queryInput.trim() })
    }
  }

  const handleLogout = () => {
    // logout() // from useAuth
    setShowUserMenu(false)
  }

  return (
    <header className="layout__header">
      <div className="flex items-center gap-4 flex-1 min-w-0">
        <button
          onClick={onMenuClick}
          className="lg:hidden p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-primary"
          aria-label="Open menu"
        >
          <MenuIcon className="w-6 h-6" />
        </button>

        <form onSubmit={handleSearchSubmit} className="relative flex-1 max-w-xl">
          <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted" aria-hidden="true" />
          <input
            ref={searchRef}
            type="text"
            value={queryInput}
            onChange={(e) => setQueryInput(e.target.value)}
            placeholder="Search metrics, query, or type Ctrl+K for command palette..."
            className="input w-full pl-10 pr-4 py-2 bg-code text-sm"
            placeholder="Search metrics, query, or type Ctrl+K for command palette..."
            autoComplete="off"
            spellCheck={false}
          />
          {queryInput && (
            <button
              type="button"
              onClick={() => setQueryInput('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-tertiary transition-colors text-muted hover:text-primary"
              aria-label="Clear search"
            >
              <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" opacity="0.7">
                <path d="M4.646 4.646a.5.5 0 01.708 0L8 7.293l2.646-2.647a.5.5 0 01.708.708L9.707 8l2.647 2.646a.5.5 0 01-.708.708L8.707 8l-2.646 2.647a.5.5 0 01-.708-.708L7.293 8 4.646 5.354a.5.5 0 11.708-.708z"/>
              </svg>
            </button>
          )}
        </form>
      </div>

      <div className="flex items-center gap-2 lg:ml-auto">
        {/* Command Palette Trigger */}
        <button
          onClick={() => setShowCommandPalette(true)}
          className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-code text-sm text-muted hover:text-primary transition-colors"
          aria-label="Open command palette (Ctrl+K)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
            <kbd className="px-1.5 py-0.5 bg-tertiary rounded text-xs font-mono">⌘</kbd>
            <span>K</span>
          </svg>
          <span className="hidden md:inline">Command Palette</span>
        </button>

        {/* Theme Toggle */}
        <button
          onClick={toggleTheme}
          className="p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-primary"
          aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
          title={`${theme === 'light' ? 'Dark' : 'Light'} mode`}
        >
          {theme === 'light' ? <MoonIcon className="w-5 h-5" /> : <SunIcon className="w-5 h-5" />}
        </button>

        {/* Notifications */}
        <div className="relative">
          <button
            onClick={() => setShowNotifications(!showNotifications)}
            className="relative p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-primary"
            aria-label="Notifications"
            aria-expanded={showNotifications}
          >
            <BellIcon className="w-5 h-5" />
            <span className="absolute -top-1 -right-1 w-4 h-4 bg-error text-white text-xs font-medium rounded-full flex items-center justify-center">
              3
            </span>
          </button>

          {showNotifications && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setShowNotifications(false)} />
              <div className="absolute right-0 top-full mt-2 w-80 bg-secondary border border-border rounded-lg shadow-lg z-50 animate-slide-in">
                <div className="p-3 border-b border-border flex items-center justify-between">
                  <h3 className="font-semibold text-sm">Notifications</h3>
                  <button className="text-muted hover:text-primary text-sm">Mark all read</button>
                </div>
                <div className="max-h-96 overflow-y-auto">
                  {[
                    { id: 1, type: 'warning', title: 'High CPU Usage', message: 'Server web-01 exceeded 90% CPU', time: '5 min ago' },
                    { id: 2, type: 'error', title: 'Query Timeout', message: 'Query exceeded 30s timeout', time: '12 min ago' },
                    { id: 3, type: 'info', title: 'New Dashboard', message: 'SOC Overview dashboard created', time: '1 hour ago' },
                  ].map((n) => (
                    <button key={n.id} className="w-full p-3 hover:bg-tertiary transition-colors text-left border-b border-border last:border-0">
                      <div className="flex items-start gap-2">
                        <span className={`w-2 h-2 rounded-full mt-1.5 flex-shrink-0 ${n.type === 'error' ? 'bg-error' : n.type === 'warning' ? 'bg-warning' : 'bg-info'}`} />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium truncate">{n.title}</p>
                          <p className="text-xs text-muted truncate">{n.message}</p>
                          <p className="text-xs text-muted mt-1">{n.time}</p>
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
                <div className="p-2 border-t border-border">
                  <button className="w-full text-center text-sm text-primary hover:text-primary-hover">View all notifications</button>
                </div>
              </div>
            </>
          )}
        </div>

        {/* User Menu */}
        <div className="relative">
          <button
            onClick={() => setShowUserMenu(!showUserMenu)}
            className="flex items-center gap-2 p-1.5 rounded-lg hover:bg-tertiary transition-colors"
            aria-expanded={showUserMenu}
            aria-haspopup="true"
          >
            <div className="w-8 h-8 rounded-full bg-primary-light text-primary flex items-center justify-center font-medium text-sm">
              {user?.name?.charAt(0) || 'U'}
            </div>
            <ChevronDownIcon className="w-4 h-4 text-muted" />
          </button>

          {showUserMenu && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setShowUserMenu(false)} />
              <div className="absolute right-0 top-full mt-2 w-56 bg-secondary border border-border rounded-lg shadow-lg z-50 animate-slide-in">
                <div className="p-3 border-b border-border">
                  <p className="font-medium text-sm truncate">{user?.name || 'User'}</p>
                  <p className="text-xs text-muted truncate">{user?.email || 'user@helix.io'}</p>
                  <p className="text-xs text-muted mt-0.5">{user?.tenantName || 'SOC Team Alpha'}</p>
                </div>
                <div className="py-1">
                  <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-secondary hover:bg-tertiary hover:text-primary transition-colors">
                    <UserIcon className="w-4 h-4" />
                    Profile
                  </button>
                  <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-secondary hover:bg-tertiary hover:text-primary transition-colors">
                    <SettingsIcon className="w-4 h-4" />
                    Settings
                  </button>
                  <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-secondary hover:bg-tertiary hover:text-primary transition-colors">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3 7 7" />
                      <path d="M9 12l2 2 4-4" />
                    </svg>
                    Keyboard Shortcuts
                  </button>
                </div>
                <div className="border-t border-border">
                  <button
                    onClick={handleLogout}
                    className="w-full flex items-center gap-2 px-3 py-2 text-sm text-error hover:bg-error-light hover:text-error transition-colors"
                  >
                    <LogoutIcon className="w-4 h-4" />
                    Logout
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  )
}

// Additional icons needed
function SearchIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="11" cy="11" r="8" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>
  )
}

function ChevronDownIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
    </svg>
  )
}

function SettingsIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 019 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 014.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 011 1.51" />
    </svg>
  )
}

export { Header }