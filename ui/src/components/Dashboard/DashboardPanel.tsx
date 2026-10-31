import React, { useState, useCallback, useRef, useEffect } from 'react'
import { ChartPanel } from '../Common/ChartPanel'
import { MetricCard } from '../Common/MetricCard'

interface DashboardPanelProps {
  id: string
  title: string
  type: 'timeseries' | 'stat' | 'heatmap' | 'table' | 'single' | 'html'
  query: string
  width: number
  height: number
  onEdit: (id: string, updates: Partial<DashboardPanelProps>) => void
  onDelete: (id: string) => void
  onToggleMinimize?: (id: string) => void
  minimized?: boolean
  collapsed?: boolean
  timezone?: string
  legend?: 'always' | 'hidden' | 'hover'
  yaxis?: 'left' | 'right' | 'center'
  format?: 'short' | 'ms' | 'byte' | 'duration' | 'msprecision'
  onQueryChange?: (query: string) => void
  onExecute?: () => void
  loading?: boolean
  data?: any
  error?: string
}

const PANEL_TYPES = [
  { value: 'timeseries', label: 'Time Series', icon: '📈' },
  { value: 'stat', label: 'Stat', icon: '📊' },
  { value: 'heatmap', label: 'Heatmap', icon: '🔥' },
  { value: 'table', label: 'Table', icon: '📋' },
  { value: 'single', label: 'Single Stat', icon: '🔢' },
  { value: 'html', label: 'HTML', icon: '🌐' },
]

export function DashboardPanel({
  id,
  title,
  type,
  query,
  width,
  height,
  onEdit,
  onDelete,
  onToggleMinimize,
  minimized = false,
  collapsed = false,
  timezone = 'UTC',
  legend = 'always',
  yaxis = 'left',
  format = 'short',
  onQueryChange,
  onExecute,
  loading = false,
  data,
  error,
}: DashboardPanelProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [editTitle, setEditTitle] = useState(title)
  const [editQuery, setEditQuery] = useState(query)
  const [editType, setEditType] = useState(type)
  const [showMenu, setShowMenu] = useState(false)
  const [isCollapsed, setIsCollapsed] = useState(collapsed)
  const [isMinimized, setIsMinimized] = useState(minimized)

  const handleSave = useCallback(() => {
    onEdit(id, {
      title: editTitle,
      query: editQuery,
      type: editType,
    })
    setIsEditing(false)
  }, [id, editTitle, editQuery, editType, onEdit])

  const handleCancel = useCallback(() => {
    setEditTitle(title)
    setEditQuery(query)
    setEditType(type)
    setIsEditing(false)
  }, [title, query, type])

  const handleDelete = useCallback(() => {
    if (confirm('Delete this panel?')) {
      onDelete(id)
    }
  }, [id, onDelete])

  const handleToggleCollapse = useCallback(() => {
    setIsCollapsed(!isCollapsed)
    if (onToggleMinimize) {
      onToggleMinimize(id)
    }
  }, [isCollapsed, id, onToggleMinimize])

  const handleToggleMinimize = useCallback(() => {
    setIsMinimized(!isMinimized)
    if (onToggleMinimize) {
      onToggleMinimize(id)
    }
  }, [isMinimized, id, onToggleMinimize])

  if (isEditing) {
    return (
      <div className="card border-2 border-primary">
        <div className="card__header flex items-center justify-between">
          <input
            type="text"
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            className="input input--sm font-semibold"
            placeholder="Panel title"
          />
          <div className="flex items-center gap-1">
            <button className="btn btn--ghost btn--sm" onClick={handleCancel}>Cancel</button>
            <button className="btn btn--primary btn--sm" onClick={handleSave}>Save</button>
          </div>
        </div>
        <div className="card__body space-y-3">
          <div>
            <label className="label">Panel Type</label>
            <select
              value={editType}
              onChange={(e) => setEditType(e.target.value as any)}
              className="input input--sm"
            >
              {PANEL_TYPES.map(t => (
                <option key={t.value} value={t.value}>{t.icon} {t.label}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="label">Query</label>
            <textarea
              value={editQuery}
              onChange={(e) => setEditQuery(e.target.value)}
              className="input font-mono text-sm"
              rows={4}
              placeholder="Enter PromQL/HelixQL query..."
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="label">Width (columns)</label>
              <input
                type="number"
                className="input input--sm"
                value={width}
                onChange={(e) => onEdit(id, { width: Number(e.target.value) })}
                min={1}
                max={24}
              />
            </div>
            <div>
              <label className="label">Height (rows)</label>
              <input
                type="number"
                className="input input--sm"
                value={height}
                onChange={(e) => onEdit(id, { height: Number(e.target.value) })}
                min={1}
                max={24}
              />
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`card group ${isCollapsed ? 'opacity-75' : ''}`}>
      <div className="card__header flex items-center justify-between">
        <div className="flex items-center gap-2 flex-1 min-w-0">
          <button
            onClick={handleToggleCollapse}
            className="p-1 rounded hover:bg-tertiary transition-colors text-muted hover:text-primary"
            aria-label={isCollapsed ? 'Expand' : 'Collapse'}
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={`transition-transform ${isCollapsed ? '-rotate-90' : ''}`}>
              <polyline points="6 9 12 15 18 9" />
            </svg>
          </button>
          <h3 className="font-semibold truncate">{title}</h3>
          <span className="badge badge--neutral text-xs">{type}</span>
        </div>
        <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={() => setIsEditing(true)}
            className="p-1.5 rounded hover:bg-tertiary text-muted hover:text-primary transition-colors"
            title="Edit panel"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
              <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
            </svg>
          </button>
          <button
            onClick={handleToggleMinimize}
            className="p-1.5 rounded hover:bg-tertiary text-muted hover:text-primary transition-colors"
            title={isMinimized ? 'Maximize' : 'Minimize'}
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              {isMinimized ? (
                <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 1 2-2h3M3 16v3a2 2 0 0 0 2 2h3" />
              ) : (
                <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3" />
              )}
            </svg>
          </button>
          <button
            onClick={handleDelete}
            className="p-1.5 rounded hover:bg-error-light text-muted hover:text-error transition-colors"
            title="Delete panel"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="3 6 5 6 21 6" />
              <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
            </svg>
          </button>
        </div>
      </div>

      {!isCollapsed && !isMinimized && (
        <div className="card__body" style={{ height: `${height * 40}px` }}>
          {loading ? (
            <div className="flex items-center justify-center h-full">
              <div className="loading-spinner loading-spinner--lg" />
            </div>
          ) : error ? (
            <div className="flex items-center justify-center h-full text-error">
              <div className="text-center">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mx-auto mb-2">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="8" x2="12" y2="12" />
                  <line x1="12" y1="16" x2="12.01" y2="16" />
                </svg>
                <p className="text-sm">{error}</p>
              </div>
            </div>
          ) : type === 'timeseries' ? (
            <ChartPanel
              title=""
              type="line"
              data={data || { labels: [], datasets: [] }}
              height={height * 40}
              showLegend={legend !== 'hidden'}
              showGrid={true}
            />
          ) : type === 'stat' ? (
            <div className="flex items-center justify-center h-full">
              <div className="text-center">
                <p className="text-4xl font-bold text-primary">{data?.value || 0}</p>
                <p className="text-sm text-muted mt-1">{data?.label || title}</p>
              </div>
            </div>
          ) : type === 'table' ? (
            <div className="overflow-auto h-full">
              <table className="table">
                <thead>
                  <tr>
                    {data?.columns?.map((col: string, i: number) => (
                      <th key={i}>{col}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {data?.rows?.map((row: any[], i: number) => (
                    <tr key={i}>
                      {row.map((cell: any, j: number) => (
                        <td key={j}>{cell}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : type === 'single' ? (
            <div className="flex items-center justify-center h-full">
              <MetricCard
                title={title}
                value={data?.value || 0}
                unit={data?.unit || ''}
                trend={data?.trend}
                trendValue={data?.trendValue}
                color={data?.color || 'primary'}
              />
            </div>
          ) : (
            <div className="flex items-center justify-center h-full text-muted">
              <p>Panel type "{type}" not yet implemented</p>
            </div>
          )}
        </div>
      )}

      {!isCollapsed && isMinimized && (
        <div className="card__body p-2 text-center text-sm text-muted">
          Panel minimized
        </div>
      )}
    </div>
  )
}

export { DashboardPanel }