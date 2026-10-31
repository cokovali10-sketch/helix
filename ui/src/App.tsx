import React, { useState, useEffect } from 'react'
import { Routes, Route, Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useTheme } from './context/ThemeContext'
import { useAuth } from './context/AuthContext'
import { Layout } from './components/Layout/Layout'
import { Sidebar } from './components/Layout/Sidebar'
import { Header } from './components/Layout/Header'
import { QueryPage } from './pages/QueryPage'
import { DashboardsPage } from './pages/DashboardsPage'
import { ExplorerPage } from './pages/ExplorerPage'
import { AlertsPage } from './pages/AlertsPage'
import { SettingsPage } from './pages/SettingsPage'
import { LoginPage } from './pages/LoginPage'
import { ProtectedRoute } from './components/Common/ProtectedRoute'
import { ToastContainer } from './components/Common/ToastContainer'
import './App.css'

function App() {
  const { theme } = useTheme()
  const { isAuthenticated, isLoading } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()
  const [showMobileSidebar, setShowMobileSidebar] = useState(false)

  // Redirect to login if not authenticated and not on login page
  useEffect(() => {
    if (!isLoading && !isAuthenticated && location.pathname !== '/login') {
      navigate('/login', { replace: true })
    }
  }, [isAuthenticated, isLoading, location.pathname, navigate])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-primary">
        <div className="loading-spinner loading-spinner--lg" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return <LoginPage onLogin={() => navigate('/')} />
  }

  return (
    <div className="layout" data-theme={theme}>
      <Sidebar isOpen={showMobileSidebar} onClose={() => setShowMobileSidebar(false)} />
      <div className="layout__main">
        <Header onMenuClick={() => setShowMobileSidebar(true)} />
        <main className="layout__content">
          <Routes>
            <Route path="/login" element={<LoginPage onLogin={() => navigate('/')} />} />
            <Route element={<ProtectedRoute><Layout /></ProtectedRoute>}>
              <Route path="/" element={<Navigate to="/query" replace />} />
              <Route path="/query" element={<QueryPage />} />
              <Route path="/dashboards" element={<DashboardsPage />} />
              <Route path="/dashboards/:id" element={<DashboardsPage />} />
              <Route path="/explorer" element={<ExplorerPage />} />
              <Route path="/alerts" element={<AlertsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
      <ToastContainer />
    </div>
  )
}

export default App