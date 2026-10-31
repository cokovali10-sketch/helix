import React, { useState, useEffect } from 'react'
import { Outlet, useLocation } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Header } from './Header'

export function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false)
  const location = useLocation()

  useEffect(() => {
    // Close mobile sidebar on route change
    setMobileSidebarOpen(false)
    // On desktop, keep sidebar state; on mobile, close
    if (window.innerWidth >= 1024) {
      setSidebarOpen(true)
    }
  }, [location.pathname])

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 1024) {
        setMobileSidebarOpen(false)
      }
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  return (
    <div className="layout">
      <Sidebar
        isOpen={window.innerWidth >= 1024 ? sidebarOpen : mobileSidebarOpen}
        onClose={() => setMobileSidebarOpen(false)}
      />
      <div className="layout__main">
        <Header onMenuClick={() => setMobileSidebarOpen(true)} />
        <main className="layout__content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export { Layout }