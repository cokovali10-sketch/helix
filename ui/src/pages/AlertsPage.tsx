import React, { useState, useCallback } from 'react'
import { useToast } from '../../components/Common/ToastContainer'
import { PlusIcon, TrashIcon, EditIcon, EyeIcon, BellIcon, FilterIcon, ChevronDownIcon, ChevronRightIcon, SaveIcon, XIcon } from '../Layout/Sidebar'

interface AlertRule {
  id: string
  name: string
  expression: string
  severity: 'critical' | 'warning' | 'info'
  status: 'firing' | 'resolved' | 'acknowledged' | 'suppressed'
  duration: string
  labels: Record<string, string>
  annotations: Record<string, string>
  enabled: boolean
  createdAt: number
  updatedAt: number
  firingCount: number
  lastFired?: number
}

interface AlertsPageProps {}

export function AlertsPage() {
  const { addToast } = useToast()
  const [alerts, setAlerts] = useState<AlertRule[]>([
    { id: '1', name: 'HighCPUUsage', expression: 'cpu_usage > 90', severity: 'warning', status: 'firing', duration: '5m', labels: { severity: 'warning' }, annotations: { summary: 'High CPU usage detected' }, enabled: true, createdAt: Date.now() - 86400000, updatedAt: Date.now(), firingCount: 3, lastFired: Date.now() - 300000 },
    { id: '2', name: 'HighMemoryUsage', expression: 'memory_usage > 90', severity: 'warning', status: 'resolved', duration: '5m', labels: { severity: 'warning' }, annotations: { summary: 'High memory usage detected' }, enabled: true, createdAt: Date.now() - 172800000, updatedAt: Date.now() - 3600000, firingCount: 1, lastFired: Date.now() - 7200000 },
    { id: '3', name: 'DiskSpaceLow', expression: 'disk_usage > 95', severity: 'critical', status: 'firing', duration: '10m', labels: { severity: 'critical' }, annotations: { summary: 'Disk space critically low' }, enabled: true, createdAt: Date.now() - 259200000, updatedAt: Date.now(), firingCount: 5, lastFired: Date.now() - 600000 },
    { id: '4', name: 'ServiceDown', expression: 'up == 0', severity: 'critical', status: 'resolved', duration: '1m', labels: { severity: 'critical' }, annotations: { summary: 'Service is down' }, enabled: true, createdAt: Date.now() - 345600000, updatedAt: Date.now() - 86400000, firingCount: 2, lastFired: Date.now() - 172800000 },
    { id: '5', name: 'HighErrorRate', expression: 'error_rate > 5', severity: 'warning', status: 'acknowledged', duration: '5m', labels: { severity: 'warning' }, annotations: { summary: 'High error rate detected' }, enabled: true, createdAt: Date.now() - 432000000, updatedAt: Date.now() - 1800000, firingCount: 8, lastFired: Date.now() - 900000 },
  ])
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingAlert, setEditingAlert] = useState<AlertRule | null>(null)
  const [newAlert, setNewAlert] = useState<Partial<AlertRule>>({
    name: '', expression: '', severity: 'warning', duration: '5m', labels: {}, annotations: {}, enabled: true,
  })
  const [filter, setFilter] = useState<'all' | 'firing' | 'resolved' | 'acknowledged'>('all')
  const [searchQuery, setSearchQuery] = useState('')

  const filteredAlerts = alerts.filter(a => {
    if (filter !== 'all' && a.status !== filter) return false
    if (searchQuery && !a.name.toLowerCase().includes(searchQuery.toLowerCase()) && !a.expression.toLowerCase().includes(searchQuery.toLowerCase())) return false
    return true
  })

  const handleCreate = () => {
    const alert: AlertRule = {
      id: crypto.randomUUID(),
      name: newAlert.name || '',
      expression: newAlert.expression || '',
      severity: newAlert.severity || 'warning',
      status: 'resolved',
      duration: newAlert.duration || '5m',
      labels: newAlert.labels || {},
      annotations: newAlert.annotations || {},
      enabled: newAlert.enabled ?? true,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      firingCount: 0,
    }
    setAlerts([alert, ...alerts])
    setNewAlert({ name: '', expression: '', severity: 'warning', duration: '5m', labels: {}, annotations: {}, enabled: true })
    setShowCreateModal(false)
    addToast({ type: 'success', title: 'Alert rule created' })
  }

  const handleDelete = (id: string) => {
    if (confirm('Delete this alert rule?')) {
      setAlerts(prev => prev.filter(a => a.id !== id))
      addToast({ type: 'info', title: 'Alert rule deleted' })
    }
  }

  const handleToggle = (id: string) => {
    setAlerts(prev => prev.map(a => a.id === id ? { ...a, enabled: !a.enabled, updatedAt: Date.now() } : a))
  }

  const handleEdit = (alert: AlertRule) => {
    setEditingAlert(alert)
    setNewAlert({ name: alert.name, expression: alert.expression, severity: alert.severity, duration: alert.duration, labels: alert.labels, annotations: alert.annotations, enabled: alert.enabled })
    setShowCreateModal(true)
  }

  const handleSave = () => {
    if (editingAlert) {
      setAlerts(prev => prev.map(a => a.id === editingAlert.id ? { ...a, ...newAlert, updatedAt: Date.now() } as AlertRule : a))
      addToast({ type: 'success', title: 'Alert rule updated' })
    } else {
      handleCreate()
    }
    setShowCreateModal(false)
    setEditingAlert(null)
  }

  const stats = {
    total: alerts.length,
    firing: alerts.filter(a => a.status === 'firing').length,
    resolved: alerts.filter(a => a.status === 'resolved').length,
    acknowledged: alerts.filter(a => a.status === 'acknowledged').length,
    enabled: alerts.filter(a => a.enabled).length,
  }

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Alerts</h1>
          <p className="text-sm text-muted">Manage alert rules and notifications</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <input
              type="text"
              placeholder="Search alerts..."
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              className="input input--sm w-48"
            />
          </div>
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            {(['all', 'firing', 'resolved', 'acknowledged'] as const).map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${filter === f ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            ))}
          </div>
          <button className="btn btn--primary" onClick={() => { setNewAlert({ name: '', expression: '', severity: 'warning', duration: '5m', labels: {}, annotations: {}, enabled: true }); setEditingAlert(null); setShowCreateModal(true); }}>
            <PlusIcon className="w-4 h-4 mr-2" />
            New Alert
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 sm:grid-cols-5 gap-3 p-4 border-b border-border bg-secondary">
        {[
          { label: 'Total', value: stats.total, color: 'text-primary' },
          { label: 'Firing', value: stats.firing, color: 'text-error' },
          { label: 'Resolved', value: stats.resolved, color: 'text-success' },
          { label: 'Acknowledged', value: stats.acknowledged, color: 'text-warning' },
          { label: 'Enabled', value: stats.enabled, color: 'text-info' },
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
            <BellIcon className="empty-state__icon" />
            <h3 className="empty-state__title">No alerts found</h3>
            <p className="empty-state__description">Create your first alert rule to start monitoring</p>
            <button className="btn btn--primary mt-4" onClick={() => setShowCreateModal(true)}>Create Alert Rule</button>
          </div>
        ) : (
          <div className="space-y-3">
            {filteredAlerts.map(alert => (
              <div key={alert.id} className="card group hover:shadow-md transition-shadow">
                <div className="card__body p-4">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="font-semibold text-lg truncate">{alert.name}</h3>
                        <span className={`badge badge--${alert.severity === 'critical' ? 'error' : alert.severity === 'warning' ? 'warning' : 'info'}`}>
                          {alert.severity}
                        </span>
                        <span className={`badge badge--${alert.status === 'firing' ? 'error' : alert.status === 'resolved' ? 'success' : alert.status === 'acknowledged' ? 'warning' : 'neutral'}`}>
                          {alert.status}
                        </span>
                        {!alert.enabled && <span className="badge badge--neutral">Disabled</span>}
                      </div>
                      <code className="text-sm font-mono text-primary bg-code px-2 py-1 rounded block truncate">{alert.expression}</code>
                      <div className="flex items-center gap-4 mt-2 text-xs text-muted">
                        <span>Duration: {alert.duration}</span>
                        <span>Fired: {alert.firingCount} times</span>
                        {alert.lastFired && <span>Last fired: {new Date(alert.lastFired).toLocaleString()}</span>}
                        <span>Updated: {new Date(alert.updatedAt).toLocaleString()}</span>
                      </div>
                      {Object.keys(alert.labels).length > 0 && (
                        <div className="flex gap-1 mt-2 flex-wrap">
                          {Object.entries(alert.labels).map(([k, v]) => (
                            <span key={k} className="badge badge--neutral text-xs">{k}={v}</span>
                          ))}
                        </div>
                      )}
                    </div>
                    <div className="flex items-center gap-1 shrink-0">
                      <button className="btn btn--ghost btn--sm" onClick={() => handleToggle(alert.id)} title={alert.enabled ? 'Disable' : 'Enable'}>
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M18.36 6.64a9 9 0 1 1-12.73 0"/><line x1="12" y1="2" x2="12" y2="12"/></svg>
                      </button>
                      <button className="btn btn--ghost btn--sm" onClick={() => handleEdit(alert)} title="Edit">
                        <EditIcon className="w-4 h-4" />
                      </button>
                      <button className="btn btn--ghost btn--sm text-error hover:text-error" onClick={() => handleDelete(alert.id)} title="Delete">
                        <TrashIcon className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showCreateModal && (
        <div className="modal__overlay" onClick={() => { setShowCreateModal(false); setEditingAlert(null); }}>
          <div className="modal">
            <div className="modal__header">
              <h2 className="modal__title">{editingAlert ? 'Edit Alert Rule' : 'Create Alert Rule'}</h2>
              <button className="modal__close" onClick={() => { setShowCreateModal(false); setEditingAlert(null); }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
            <div className="modal__body">
              <div className="space-y-4">
                <div>
                  <label className="label" htmlFor="name">Name</label>
                  <input id="name" type="text" className="input" value={newAlert.name || ''} onChange={e => setNewAlert({...newAlert, name: e.target.value})} placeholder="Alert name" required />
                </div>
                <div>
                  <label className="label" htmlFor="expression">Expression</label>
                  <input id="expression" type="text" className="input font-mono" value={newAlert.expression || ''} onChange={e => setNewAlert({...newAlert, expression: e.target.value})} placeholder="cpu_usage > 90" required />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="label" htmlFor="severity">Severity</label>
                    <select id="severity" className="input" value={newAlert.severity || 'warning'} onChange={e => setNewAlert({...newAlert, severity: e.target.value as any})}>
                      <option value="info">Info</option>
                      <option value="warning">Warning</option>
                      <option value="critical">Critical</option>
                    </select>
                  </div>
                  <div>
                    <label className="label" htmlFor="duration">Duration</label>
                    <select id="duration" className="input" value={newAlert.duration || '5m'} onChange={e => setNewAlert({...newAlert, duration: e.target.value})}>
                      <option value="1m">1 minute</option>
                      <option value="5m">5 minutes</option>
                      <option value="10m">10 minutes</option>
                      <option value="15m">15 minutes</option>
                      <option value="30m">30 minutes</option>
                      <option value="1h">1 hour</option>
                    </select>
                  </div>
                </div>
                <div>
                  <label className="label">Labels (key=value, comma-separated)</label>
                  <input type="text" className="input" placeholder="severity=warning, team=soc" onChange={e => {
                    const labels: Record<string, string> = {}
                    e.target.value.split(',').forEach(pair => {
                      const [k, v] = pair.split('=')
                      if (k && v) labels[k.trim()] = v.trim()
                    })
                    setNewAlert({...newAlert, labels})
                  }} />
                </div>
                <div>
                  <label className="label">Annotations (key=value, comma-separated)</label>
                  <input type="text" className="input" placeholder="summary=High CPU usage" onChange={e => {
                    const annotations: Record<string, string> = {}
                    e.target.value.split(',').forEach(pair => {
                      const [k, v] = pair.split('=')
                      if (k && v) annotations[k.trim()] = v.trim()
                    })
                    setNewAlert({...newAlert, annotations})
                  }} />
                </div>
                <div className="flex items-center gap-2">
                  <input type="checkbox" id="enabled" className="w-4 h-4 rounded border-border text-primary focus:ring-primary" checked={newAlert.enabled ?? true} onChange={e => setNewAlert({...newAlert, enabled: e.target.checked})} />
                  <label className="label" htmlFor="enabled">Enabled</label>
                </div>
              </div>
            </div>
            <div className="modal__footer">
              <button className="btn btn--secondary" onClick={() => { setShowCreateModal(false); setEditingAlert(null); }}>Cancel</button>
              <button className="btn btn--primary" onClick={handleSave}>{editingAlert ? 'Save Changes' : 'Create Alert Rule'}</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function PlusIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
}
function TrashIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
}
function EditIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 1-4 1.5-1.5a2.121 2.121 0 010-3L6.5 2.5a2.121 2.121 0 00-3 0l-7 7a2.121 2.121 0 000 3l7 7"/></svg>
}
function BellIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-3-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>
}
function FilterIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
}
function ChevronDownIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
}
function ChevronRightIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
}
function SaveIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
}
function XIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
}

export { AlertsPage }