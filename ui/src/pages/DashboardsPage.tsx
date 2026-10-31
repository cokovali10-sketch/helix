import React, { useState, useCallback, useEffect } from 'react'
import { useQuery } from '../../context/QueryContext'
import { useToast } from '../../components/Common/ToastContainer'
import { PlusIcon, BarChartIcon, GridIcon, TrashIcon, EditIcon, EyeIcon, CopyIcon, ChevronDownIcon } from '../Layout/Sidebar'

interface DashboardsPageProps {}

export function DashboardsPage() {
  const { addToast } = useToast()
  const [dashboards, setDashboards] = useState([
    { id: '1', title: 'SOC Overview', description: 'Overall security operations center metrics', panels: 12, updatedAt: Date.now() - 86400000, owner: 'Alex Chen', isPublic: true },
    { id: '2', title: 'Network Traffic Analysis', description: 'Network flow and bandwidth monitoring', panels: 8, updatedAt: Date.now() - 172800000, owner: 'Sam Wilson', isPublic: true },
    { id: '3', title: 'Host Performance', description: 'CPU, memory, disk, and network per host', panels: 15, updatedAt: Date.now() - 3600000, owner: 'Alex Chen', isPublic: false },
    { id: '4', title: 'Application Latency', description: 'Service latency percentiles and error rates', panels: 10, updatedAt: Date.now() - 7200000, owner: 'Sam Wilson', isPublic: true },
  ])
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingDashboard, setEditingDashboard] = useState<any>(null)
  const [newDashboard, setNewDashboard] = useState({ title: '', description: '', isPublic: true })

  const handleCreate = () => {
    const dashboard = {
      id: crypto.randomUUID(),
      title: newDashboard.title,
      description: newDashboard.description,
      panels: 0,
      updatedAt: Date.now(),
      owner: 'Current User',
      isPublic: newDashboard.isPublic,
    }
    setDashboards([dashboard, ...dashboards])
    setNewDashboard({ title: '', description: '', isPublic: true })
    setShowCreateModal(false)
  }

  const handleDelete = (id: string) => {
    if (confirm('Delete this dashboard?')) {
      setDashboards(prev => prev.filter(d => d.id !== id))
    }
  }

  const handleEdit = (dashboard: any) => {
    setEditingDashboard(dashboard)
    setNewDashboard({ title: dashboard.title, description: dashboard.description, isPublic: dashboard.isPublic })
    setShowCreateModal(true)
  }

  const handleSave = () => {
    if (editingDashboard) {
      setDashboards(prev => prev.map(d => d.id === editingDashboard.id ? { ...d, ...newDashboard, updatedAt: Date.now() } : d))
    } else {
      handleCreate()
    }
  }

  const DashboardCard = ({ dashboard }: { dashboard: any }) => (
    <div className="card group hover:shadow-md transition-shadow">
      <div className="card__header flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-lg truncate">{dashboard.title}</h3>
          <p className="text-sm text-muted mt-1 truncate">{dashboard.description}</p>
        </div>
        <div className="flex items-center gap-1">
          <span className={`badge ${dashboard.isPublic ? 'badge--success' : 'badge--neutral'}`}>
            {dashboard.isPublic ? 'Public' : 'Private'}
          </span>
        </div>
      </div>
      <div className="card__body">
        <div className="flex items-center justify-between text-sm text-muted">
          <span>{dashboard.panels} panels</span>
          <span>Updated {new Date(dashboard.updatedAt).toLocaleDateString()}</span>
        </div>
        <div className="mt-3 flex items-center gap-2 text-xs text-muted">
          <span className="flex items-center gap-1">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
            {dashboard.owner}
          </span>
        </div>
      </div>
      <div className="card__footer">
        <button className="btn btn--ghost btn--sm" title="View"><EyeIcon className="w-4 h-4" /></button>
        <button className="btn btn--ghost btn--sm" onClick={() => { /* edit */ }} title="Edit"><EditIcon className="w-4 h-4" /></button>
        <button className="btn btn--ghost btn--sm" title="Duplicate"><CopyIcon className="w-4 h-4" /></button>
        <button className="btn btn--ghost btn--sm text-error hover:text-error" onClick={() => { /* delete */ }} title="Delete"><TrashIcon className="w-4 h-4" /></button>
      </div>
    </div>
  )

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Dashboards</h1>
          <p className="text-sm text-muted">Create and manage observability dashboards</p>
        </div>
        <div className="flex items-center gap-3">
          <button className="btn btn--ghost btn--sm" onClick={() => setViewMode('grid')} aria-pressed={viewMode === 'grid'} title="Grid view">
            <GridIcon className="w-4 h-4" />
          </button>
          <button className="btn btn--ghost btn--sm" onClick={() => setViewMode('list')} aria-pressed={viewMode === 'list'} title="List view">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
          </button>
          <button className="btn btn--primary" onClick={() => { setNewDashboard({ title: '', description: '', isPublic: true }); setEditingDashboard(null); setShowCreateModal(true); }}>
            <PlusIcon className="w-4 h-4 mr-2" />
            New Dashboard
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4">
        {viewMode === 'grid' ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {dashboards.map(d => <DashboardCard key={d.id} dashboard={d} />)}
          </div>
        ) : (
          <div className="space-y-3">
            {dashboards.map(d => (
              <div key={d.id} className="card group hover:shadow-md transition-shadow">
                <div className="card__body p-4 flex items-center justify-between">
                  <div className="flex items-center gap-4 flex-1 min-w-0">
                    <div className="w-12 h-12 rounded-lg bg-primary-light flex items-center justify-center">
                      <BarChartIcon className="w-6 h-6 text-primary" />
                    </div>
                    <div className="min-w-0">
                      <h3 className="font-semibold truncate">{d.title}</h3>
                      <p className="text-sm text-muted truncate">{d.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4 text-sm text-muted shrink-0">
                    <span>{d.panels} panels</span>
                    <span className={`badge ${d.isPublic ? 'badge--success' : 'badge--neutral'}`}>
                      {d.isPublic ? 'Public' : 'Private'}
                    </span>
                    <span>Updated {new Date(d.updatedAt).toLocaleDateString()}</span>
                    <span>{d.owner}</span>
                  </div>
                  <div className="flex items-center gap-1 shrink-0">
                    <button className="btn btn--ghost btn--sm" title="View"><EyeIcon className="w-4 h-4" /></button>
                    <button className="btn btn--ghost btn--sm" title="Edit"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-1 4 4-1 9.5-9.5a2.121 2.121 0 00-3-3L6 11l-4-1-1 4 1-4 1.5-1.5a2.121 2.121 0 010-3L6.5 2.5a2.121 2.121 0 00-3 0l-7 7a2.121 2.121 0 000 3l7 7"/></svg>
                    </button>
                    <button className="btn btn--ghost btn--sm text-error hover:text-error" title="Delete"><TrashIcon className="w-4 h-4" /></button>
                  </div>
                </div>
              ))}
            )}
          </div>
        )}
        {dashboards.length === 0 && (
          <div className="empty-state">
            <svg className="empty-state__icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
            <h3 className="empty-state__title">No dashboards yet</h3>
            <p className="empty-state__description">Create your first dashboard to start visualizing metrics</p>
            <button className="btn btn--primary mt-4" onClick={() => setShowCreateModal(true)}>Create Dashboard</button>
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showCreateModal && (
        <div className="modal__overlay" onClick={() => { setShowCreateModal(false); setEditingDashboard(null); setNewDashboard({ title: '', description: '', isPublic: true }); }}>
          <div className="modal">
            <div className="modal__header">
              <h2 className="modal__title">{editingDashboard ? 'Edit Dashboard' : 'Create Dashboard'}</h2>
              <button className="modal__close" onClick={() => { setShowCreateModal(false); setEditingDashboard(null); setNewDashboard({ title: '', description: '', isPublic: true }); }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
            <div className="modal__body">
              <div className="space-y-4">
                <div>
                  <label className="label" htmlFor="title">Title</label>
                  <input id="title" type="text" className="input" value={newDashboard.title} onChange={e => setNewDashboard({...newDashboard, title: e.target.value})} placeholder="Dashboard title" required />
                </div>
                <div>
                  <label className="label" htmlFor="description">Description</label>
                  <textarea id="description" className="input" rows={3} value={newDashboard.description} onChange={e => setNewDashboard({...newDashboard, description: e.target.value})} placeholder="Dashboard description" />
                </div>
                <div className="flex items-center gap-2">
                  <input type="checkbox" id="isPublic" className="w-4 h-4 rounded border-border text-primary focus:ring-primary" checked={newDashboard.isPublic} onChange={e => setNewDashboard({...newDashboard, isPublic: e.target.checked})} />
                  <label className="label" htmlFor="isPublic">Make this dashboard public</label>
                </div>
              </div>
            </div>
            <div className="modal__footer">
              <button className="btn btn--secondary" onClick={() => { setShowCreateModal(false); setEditingDashboard(null); setNewDashboard({ title: '', description: '', isPublic: true }); }}>Cancel</button>
              <button className="btn btn--primary" onClick={handleSave}>{editingDashboard ? 'Save Changes' : 'Create Dashboard'}</button>
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
function BarChartIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
}
function EyeIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M1 12s4-8 11-8 11 8 11-8"/><circle cx="12" cy="12" r="3"/></svg>
}
function EditIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 9.5-9.5a2.121 2.121 0 013 3L6 17l-4 1 1-4 1-4 1.5-1.5a2.121 2.121 0 010-3L6.5 2.5a2.121 2.121 0 00-3 0l-7 7a2.121 2.121 0 000 3l7 7"/></svg>
}
function CopyIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
}
function TrashIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
}
function GridIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
}

export { DashboardsPage }