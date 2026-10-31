import React, { useState, useEffect } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import { useTheme } from '../../../context/ThemeContext'
import { useAuth } from '../../../context/AuthContext'

interface SidebarProps {
  isOpen: boolean
  onClose: () => void
}

const navigation = [
  { path: '/query', label: 'Query', icon: QueryIcon, description: 'Query and analyze metrics' },
  { path: '/dashboards', label: 'Dashboards', icon: DashboardIcon, description: 'Create and manage dashboards' },
  { path: '/explorer', label: 'Explorer', icon: ExplorerIcon, description: 'Explore metrics and metadata' },
  { path: '/alerts', label: 'Alerts', icon: AlertIcon, description: 'Manage alert rules and notifications' },
  { path: '/settings', label: 'Settings', icon: SettingsIcon, description: 'Configure platform settings' },
]

function QueryIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="11" cy="11" r="8" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>
  )
}

function DashboardIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="3" width="7" height="7" rx="1" />
      <rect x="14" y="3" width="7" height="7" rx="1" />
      <rect x="3" y="14" width="7" height="7" rx="1" />
      <rect x="14" y="14" width="7" height="7" rx="1" />
    </svg>
  )
}

function ExplorerIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
  )
}

function AlertIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
      <line x1="12" y1="9" x2="12" y2="13" />
      <line x1="12" y1="17" x2="12.01" y2="17" />
    </svg>
  )
}

function SettingsIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 019 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 014.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 01-1.51 1z" />
    </svg>
  )
}

function ChevronRightIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
    </svg>
  )
}

function MenuIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  )
}

function SunIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="5" />
      <line x1="12" y1="1" x2="12" y2="3" />
      <line x1="12" y1="21" x2="12" y2="23" />
      <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
      <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
      <line x1="1" y1="12" x2="3" y2="12" />
      <line x1="21" y1="12" x2="23" y2="12" />
      <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
      <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
    </svg>
  )
}

function MoonIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
    </svg>
  )
}

function UserIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  )
}

function LogoutIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  )
}

function BellIcon({ className }: { className?: string }) {
  return (
    <svg className={className} width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-3-3-9" />
      <path d="M13.73 21a2 2 0 01-3.46 0" />
    </svg>
  )
}

export function Sidebar({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) {
  const { theme, toggleTheme } = useTheme()
  const { user, logout } = useAuth()
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  useEffect(() => {
    if (isOpen) {
      setMobileOpen(true)
      document.body.style.overflow = 'hidden'
    } else {
      setMobileOpen(false)
      document.body.style.overflow = ''
    }
    return () => {
      document.body.style.overflow = ''
    }
  }, [isOpen])

  const handleNavClick = () => {
    if (window.innerWidth < 1024) {
      onClose()
    }
  }

  if (!isOpen && window.innerWidth < 1024) {
    return null
  }

  return (
    <>
      {window.innerWidth < 1024 && isOpen && (
        <div 
          className="sidebar-overlay fixed inset-0 z-15 bg-black/50"
          onClick={onClose}
          aria-hidden="true"
        />
      )}
      <aside
        className={`layout__sidebar fixed left-0 top-0 h-full z-20 transition-all duration-300 ease-in-out ${
          isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        } ${collapsed ? 'w-18' : 'w-72'}`}
        style={{ width: collapsed ? '72px' : '288px' }}
        aria-label="Main navigation"
      >
        {/* Logo */}
        <div className="flex items-center justify-between h-16 px-4 border-b border-border">
          {!collapsed && (
            <NavLink to="/" className="flex items-center gap-2 font-bold text-lg text-primary" aria-label="Helix Home">
              <svg width="28" height="28" viewBox="0 0 32 32" fill="none">
                <rect width="32" height="32" rx="6" fill="currentColor" />
                <path d="M8 16L14 22 24 10" stroke="white" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" fill="none"/>
              </svg>
              <span className="font-semibold text-primary">Helix</span>
            </NavLink>
          )}
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-primary flex-shrink-0"
            aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
            aria-expanded={!collapsed}
          >
            <ChevronRightIcon className={`w-5 h-5 transition-transform duration-200 ${collapsed ? 'rotate-180' : ''}`} />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-3 space-y-1 overflow-y-auto" aria-label="Main navigation">
          {navigation.map((item) => {
            const isActive = location.pathname === item.path || location.pathname.startsWith(item.path + '/')
            return (
              <NavLink
                key={item.path}
                to={item.path}
                onClick={handleNavClick}
                className={({ isActive }) => `
                  flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200
                  ${isActive
                    ? 'bg-primary-light text-primary'
                    : 'text-secondary hover:bg-tertiary hover:text-primary'
                  }
                  ${collapsed ? 'justify-center px-2' : ''}
                `}
                title={collapsed ? item.label : undefined}
                aria-current={isActive ? 'page' : undefined}
                aria-label={collapsed ? item.label : undefined}
              >
                <item.icon className="w-5 h-5 flex-shrink-0 transition-transform duration-200" aria-hidden="true" />
                {!collapsed && (
                  <>
                    <span className="flex-1 truncate">{item.label}</span>
                    <span className="text-xs text-muted opacity-0 group-hover:opacity-100 transition-opacity">
                      {item.description}
                    </span>
                  </>
                )}
              </NavLink>
            )
          })}
        </nav>

        {/* Bottom section */}
        <div className="p-3 border-t border-border">
          {!collapsed && (
            <div className="mb-3 p-3 bg-tertiary rounded-lg">
              <p className="text-xs text-muted mb-1">Tenant</p>
              <p className="font-medium text-sm truncate">{user?.tenantName || 'SOC Team Alpha'}</p>
            </div>
          )}

          <div className="flex items-center gap-2">
            <button
              onClick={() => toggleTheme()}
              className="p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-primary flex-shrink-0"
              aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
              title={`${theme === 'light' ? 'Dark' : 'Light'} mode`}
            >
              {theme === 'light' ? <MoonIcon className="w-5 h-5" /> : <SunIcon className="w-5 h-5" />}
            </button>

            {!collapsed && (
              <div className="flex-1 flex items-center gap-2 ml-2">
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted truncate">{user?.name || 'User'}</p>
                  <p className="text-xs text-muted truncate">{user?.tenantName || 'SOC Team'}</p>
                </div>
              </div>
            )}

            <button
              onClick={() => logout()}
              className="p-2 rounded-lg hover:bg-tertiary transition-colors text-muted hover:text-error flex-shrink-0"
              aria-label="Logout"
              title="Logout"
            >
              <LogoutIcon className="w-5 h-5" />
            </button>
          </div>
        </div>
      </aside>

      {/* Mobile overlay */}
      {isOpen && window.innerWidth < 1024 && (
        <div
          className="fixed inset-0 bg-black/50 z-15 lg:hidden"
          onClick={() => setMobileOpen(false)}
          aria-hidden="true"
        />
      )}
    </>
  )
}