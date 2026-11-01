import React, { useState, useCallback, useEffect } from 'react'
import { useToast } from '../../components/Common/ToastContainer'
import { AlertRuleEditor } from './AlertRuleEditor'

interface Alert {
  id: string
  name: string
  severity: 'critical' | 'warning' | 'info'
  status: 'firing' | 'resolved' | 'acknowledged' | 'suppressed'
  labels: Record<string, string>
  annotations: Record<string, string>
  startsAt: number
  endsAt?: number
  updatedAt: number
  firingCount: number
  lastFired?: number
  value?: number
}

interface AlertManagerProps {
  alerts: Alert[]
  onAcknowledge: (id: string) => void
  onSilence: (id: string, duration: number) => void
  onDelete: (id: string) => void
  onEdit: (alert: Alert) => void
}

export function AlertManager({ alerts, onAcknowledge, onSilence, onDelete, onEdit }: AlertManagerProps) {
  const { addToast } = useToast()
  const [selectedAlert, setSelectedAlert] = useState<Alert | null>(null)
  const [showSilenceModal, setShowSilenceModal] = useState(false)
  const [silenceDuration, setSilenceDuration] = useState(3600000)
  const [filter, setFilter] = useState<'all' | 'firing' | 'resolved' | 'acknowledged' | 'suppressed'>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [sortBy, setSortBy] = useState<'name' | 'severity' | 'status' | 'lastFired'>('lastFired')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')

  const filteredAlerts = alerts
    .filter(a => {
      if (filter !== 'all' && a.status !== filter) return false
      if (searchQuery && !a.name.toLowerCase().includes(searchQuery.toLowerCase())) return false
      return true
    })
    .sort((a, b) => {
      let comparison = 0
      switch (sortBy) {
        case 'name':
          comparison = a.name.localeCompare(b.name)
          break
        case 'severity':
          const severityOrder = { critical: 0, warning: 1, info: 2 }
          comparison = severityOrder[a.severity] - severityOrder[b.severity]
          break
        case 'status':
          comparison = a.status.localeCompare(b.status)
          break
        case 'lastFired':
          comparison = (a.lastFired || 0) - (b.lastFired || 0)
          break
      }
      return sortDir === 'asc' ? comparison : -comparison
    })

  const handleAcknowledge = useCallback((id: string) => {
    onAcknowledge(id)
    addToast({ type: 'success', title: 'Alert acknowledged' })
  }, [onAcknowledge, addToast])

  const handleSilence = useCallback((id: string) => {
    setSelectedAlert(alerts.find(a => a.id === id) || null)
    setShowSilenceModal(true)
  }, [alerts])

  const handleSilenceConfirm = useCallback(() => {
    if (selectedAlert) {
      onSilence(selectedAlert.id, silenceDuration)
      addToast({ type: 'success', title: 'Alert silenced', message: `For ${silenceDuration / 60000} minutes` })
      setShowSilenceModal(false)
      setSelectedAlert(null)
    }
  }, [selectedAlert, silenceDuration, onSilence, addToast])

  const handleDelete = useCallback((id: string) => {
    if (confirm('Delete this alert?')) {
      onDelete(id)
      addToast({ type: 'info', title: 'Alert deleted' })
    }
  }, [onDelete, addToast])

  const severityColors = {
    critical: 'text-error bg-error-light border-error',
    warning: 'text-warning bg-warning-light border-warning',
    info: 'text-info bg-info-light border-info',
  }

  const statusColors = {
    firing: 'badge--error',
    resolved: 'badge--success',
    acknowledged: 'badge--warning',
    suppressed: 'badge--neutral',
  }

  const stats = {
    total: alerts.length,
    firing: alerts.filter(a => a.status === 'firing').length,
    resolved: alerts.filter(a => a.status === 'resolved').length,
    acknowledged: alerts.filter(a => a.status === 'acknowledged').length,
    suppressed: alerts.filter(a => a.status === 'suppressed').length,
  }

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Alert Manager</h1>
          <p className="text-sm text-muted">Manage active alerts and notifications</p>
        </div>
        <div className="flex items-center gap-3">
          <input
            type="text"
            placeholder="Search alerts..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="input input--sm w-48"
          />
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            {(['all', 'firing', 'resolved', 'acknowledged', 'suppressed'] as const).map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-2 py-1 rounded text-xs font-medium transition-colors ${filter === f ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            ))}
          </div>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as any)}
            className="input input--sm w-auto"
          >
            <option value="lastFired">Last Fired</option>
            <option value="name">Name</option>
            <option value="severity">Severity</option>
            <option value="status">Status</option>
          </select>
          <button
            onClick={() => setSortDir(d => d === 'asc' ? 'desc' : 'asc')}
            className="btn btn--ghost btn--sm"
          >
            {sortDir === 'asc' ? '↑' : '↓'}
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 sm:grid-cols-5 gap-3 p-4 border-b border-border bg-secondary">
        {[
          { label: 'Total', value: stats.total, color: 'text-primary' },
          { label: 'Firing', value: stats.firing, color: 'text-error' },
          { label: 'Resolved', value: stats.resolved, color: 'text-success' },
          { label: 'Acknowledged', value: stats.acknowledged, color: 'text-warning' },
          { label: 'Suppressed', value: stats.suppressed, color: 'text-muted' },
        ].map(s => (
          <div key={s.label} className="card p-3 text-center">
            <p className="text-xs text-muted">{s.label}</p>
            <p className={`text-xl font-bold ${s.color}`}>{s.value}</p>
          </div>
        ))}
      </div>

      <div className="flex-1 overflow-auto p-4">
        {filteredAlerts.length === 0 ? (
          <div className="empty-state">
            <svg className="empty-state__icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-3-3-9" />
              <path d="M13.73 21a2 2 0 01-3.46 0" />
            </svg>
            <h3 className="empty-state__title">No alerts</h3>
            <p className="empty-state__description">No alerts match your current filters</p>
          </div>
        ) : (
          <div className="space-y-3">
            {filteredAlerts.map(alert => (
              <div key={alert.id} className={`card group hover:shadow-md transition-shadow border-l-4 ${alert.severity === 'critical' ? 'border-l-error' : alert.severity === 'warning' ? 'border-l-warning' : 'border-l-info'}`}>
                <div className="card__body p-4">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-2 flex-wrap">
                        <h3 className="font-semibold text-lg">{alert.name}</h3>
                        <span className={`badge ${severityColors[alert.severity]}`}>{alert.severity}</span>
                        <span className={`badge ${statusColors[alert.status]}`}>{alert.status}</span>
                        {alert.value !== undefined && (
                          <span className="text-sm text-muted">Value: {alert.value}</span>
                        )}
                      </div>
                      {Object.keys(alert.labels).length > 0 && (
                        <div className="flex gap-1 mb-2 flex-wrap">
                          {Object.entries(alert.labels).map(([k, v]) => (
                            <span key={k} className="badge badge--neutral text-xs">{k}={v}</span>
                          ))}
                        </div>
                      )}
                      {Object.keys(alert.annotations).length > 0 && (
                        <div className="text-sm text-muted mb-2">
                          {Object.entries(alert.annotations).map(([k, v]) => (
                            <p key={k}><span className="font-medium">{k}:</span> {v}</p>
                          ))}
                        </div>
                      )}
                      <div className="flex items-center gap-4 text-xs text-muted">
                        <span>Started: {new Date(alert.startsAt).toLocaleString()}</span>
                        {alert.lastFired && <span>Last fired: {new Date(alert.lastFired).toLocaleString()}</span>}
                        <span>Fired: {alert.firingCount} times</span>
                        {alert.endsAt && <span>Ended: {new Date(alert.endsAt).toLocaleString()}</span>}
                      </div>
                    </div>
                    <div className="flex items-center gap-1 shrink-0">
                      {alert.status === 'firing' && (
                        <>
                          <button
                            onClick={() => handleAcknowledge(alert.id)}
                            className="btn btn--ghost btn--sm"
                            title="Acknowledge"
                          >
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                              <polyline points="20 6 9 17 4 12" />
                            </svg>
                          </button>
                          <button
                            onClick={() => handleSilence(alert.id)}
                            className="btn btn--ghost btn--sm"
                            title="Silence"
                          >
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                              <path d="M11 5L6 9H2v6h4l5 4V5z" />
                              <line x1="23" y1="9" x2="17" y2="15" />
                              <line x1="17" y1="9" x2="23" y2="15" />
                            </svg>
                          </button>
                        </>
                      )}
                      <button
                        onClick={() => onEdit(alert)}
                        className="btn btn--ghost btn--sm"
                        title="Edit"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
                          <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
                        </svg>
                      </button>
                      <button
                        onClick={() => handleDelete(alert.id)}
                        className="btn btn--ghost btn--sm text-error hover:text-error"
                        title="Delete"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <polyline points="3 6 5 6 21 6" />
                          <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Silence Modal */}
      {showSilenceModal && selectedAlert && (
        <div className="modal__overlay" onClick={() => { setShowSilenceModal(false); setSelectedAlert(null); }}>
          <div className="modal">
            <div className="modal__header">
              <h2 className="modal__title">Silence Alert</h2>
              <button className="modal__close" onClick={() => { setShowSilenceModal(false); setSelectedAlert(null); }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <div className="modal__body">
              <p className="mb-4">Silence <strong>{selectedAlert.name}</strong> for:</p>
              <div className="space-y-2">
                {[
                  { label: '15 minutes', value: 900000 },
                  { label: '30 minutes', value: 1800000 },
                  { label: '1 hour', value: 3600000 },
                  { label: '2 hours', value: 7200000 },
                  { label: '6 hours', value: 21600000 },
                  { label: '12 hours', value: 43200000 },
                  { label: '24 hours', value: 86400000 },
                ].map(opt => (
                  <label key={opt.value} className={`flex items-center gap-2 p-2 rounded cursor-pointer ${silenceDuration === opt.value ? 'bg-primary-light' : 'hover:bg-tertiary'}`}>
                    <input
                      type="radio"
                      name="duration"
                      value={opt.value}
                      checked={silenceDuration === opt.value}
                      onChange={() => setSilenceDuration(opt.value)}
                      className="w-4 h-4 text-primary"
                    />
                    <span className="text-sm">{opt.label}</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="modal__footer">
              <button className="btn btn--secondary" onClick={() => { setShowSilenceModal(false); setSelectedAlert(null); }}>Cancel</button>
              <button className="btn btn--primary" onClick={handleSilenceConfirm}>Silence</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export { AlertManager }