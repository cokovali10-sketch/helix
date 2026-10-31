import React, { createContext, useContext, useState, useCallback, useEffect } from 'react'

type ToastType = 'info' | 'success' | 'warning' | 'error'

interface Toast {
  id: string
  type: ToastType
  title: string
  message?: string
  duration?: number
  action?: {
    label: string
    onClick: () => void
  }
}

interface ToastContextType {
  toasts: Toast[]
  addToast: (toast: Omit<Toast, 'id'>) => string
  removeToast: (id: string) => void
  clearToasts: () => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])

  const addToast = useCallback((toast: Omit<Toast, 'id'>) => {
    const id = crypto.randomUUID()
    const newToast = { ...toast, id }
    setToasts(prev => [...prev, newToast])
    return id
  }, [])

  const removeToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(t => t.id !== id))
  }, [])

  const clearToasts = useCallback(() => {
    setToasts([])
  }, [])

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast, clearToasts }}>
      {children}
    </ToastContext.Provider>
  )
}

export function useToast() {
  const context = React.useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}

interface ToastProps {
  toast: Toast
  onClose: (id: string) => void
}

function Toast({ toast, onClose }: ToastProps) {
  useEffect(() => {
    if (toast.duration !== 0) {
      const timer = setTimeout(() => onClose(toast.id), toast.duration ?? 5000)
      return () => clearTimeout(timer)
    }
  }, [toast, onClose])

  const icons = {
    info: (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
      </svg>
    ),
    success: (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
      </svg>
    ),
    warning: (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM10 13a1 1 0 011 1v2a1 1 0 11-2 0v-2a1 1 0 011-1zm-1-8a1 1 0 00-1 1v3a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
    error: (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
      </svg>
    ),
  }

  const colors = {
    info: 'text-info bg-info-light border-info',
    success: 'text-success bg-success-light border-success',
    warning: 'text-warning bg-warning-light border-warning',
    error: 'text-error bg-error-light border-error',
  }

  const handleAction = () => {
    if (toast.action) {
      toast.action.onClick()
      onClose(toast.id)
    }
  }

  return (
    <div className={`flex items-start gap-3 p-4 rounded-lg border ${colors[toast.type]} animate-slide-in`} role="alert">
      <div className="flex-shrink-0">
        {icons[toast.type]}
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-medium text-sm">{toast.title}</p>
        {toast.message && (
          <p className="text-sm opacity-90 mt-1">{toast.message}</p>
        )}
        {toast.action && (
          <button
            onClick={handleAction}
            className="mt-3 text-sm font-medium underline hover:no-underline focus:outline-none focus:ring-2 focus:ring-offset-2"
            style={{ 
              '--focus-ring-color': `var(--color-${toast.type})` 
            }}
          >
            {toast.action.label}
          </button>
        )}
      </div>
      <button
        onClick={() => onClose(toast.id)}
        className="flex-shrink-0 p-1 rounded hover:bg-black/10 transition-colors"
        aria-label="Dismiss"
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" opacity="0.7">
          <path d="M4.646 4.646a.5.5 0 01.708 0L8 7.293l2.646-2.647a.5.5 0 01.708.708L9.707 8l2.647 2.646a.5.5 0 01-.708.708L8.707 8.707l-2.646 2.647a.5.5 0 01-.708-.708L7.293 8 4.646 5.354a.5.5 0 11.708-.708z"/>
        </svg>
      </button>
    </div>
  )
}

interface ToastContainerProps {
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left'
  maxToasts?: number
}

export function ToastContainer({ position = 'top-right', maxToasts = 5 }: ToastContainerProps) {
  const { toasts, removeToast } = useToast()

  const positions = {
    'top-right': 'top-4 right-4',
    'top-left': 'top-4 left-4',
    'bottom-right': 'bottom-4 right-4',
    'bottom-left': 'bottom-4 left-4',
  }

  const visibleToasts = toasts.slice(-maxToasts)

  return (
    <div
      className={`fixed ${positions[position]} z-50 flex flex-col gap-2 pointer-events-none`}
      role="region"
      aria-label="Notifications"
      aria-live="polite"
    >
      {visibleToasts.map(toast => (
        <Toast
          key={toast.id}
          toast={toast}
          onClose={removeToast}
        />
      ))}
    </div>
  )
}

export function useToast() {
  const context = React.useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}

export { ToastProvider }