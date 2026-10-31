import React, { useState, useCallback, useEffect, useRef } from 'react'
import { DashboardPanel } from '../Dashboard/DashboardPanel'
import { PlusIcon, TrashIcon, SaveIcon, DownloadIcon, UploadIcon, GridIcon, LayoutIcon, SettingsIcon } from '../Layout/Sidebar'

interface DashboardBuilderProps {
  initialPanels?: any[]
  onSave?: (panels: any[]) => void
  onExport?: (panels: any[]) => void
  onImport?: (panels: any[]) => void
}

interface Panel {
  id: string
  title: string
  type: 'timeseries' | 'stat' | 'heatmap' | 'table' | 'single' | 'html'
  query: string
  width: number
  height: number
  x: number
  y: number
  minimized?: boolean
  collapsed?: boolean
}

const PANEL_TEMPLATES = [
  { type: 'timeseries', title: 'Time Series', query: 'rate(cpu_usage[5m])', width: 6, height: 3 },
  { type: 'stat', title: 'Stat', query: 'up', width: 3, height: 2 },
  { type: 'table', title: 'Table', query: 'topk(10, rate(cpu_usage[5m]))', width: 6, height: 3 },
  { type: 'single', title: 'Single Stat', query: 'avg(cpu_usage)', width: 3, height: 2 },
]

const GRID_COLUMNS = 12
const GRID_ROW_HEIGHT = 40

export function DashboardBuilder({ initialPanels = [], onSave, onExport, onImport }: DashboardBuilderProps) {
  const [panels, setPanels] = useState<Panel[]>(initialPanels)
  const [selectedPanel, setSelectedPanel] = useState<string | null>(null)
  const [isDragging, setIsDragging] = useState(false)
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 })
  const [showTemplateModal, setShowTemplateModal] = useState(false)
  const [showGrid, setShowGrid] = useState(true)
  const [zoom, setZoom] = useState(1)
  const containerRef = useRef<HTMLDivElement>(null)

  const handleAddPanel = useCallback((template: any) => {
    const newPanel: Panel = {
      id: crypto.randomUUID(),
      title: template.title,
      type: template.type,
      query: template.query,
      width: template.width,
      height: template.height,
      x: 0,
      y: 0,
      minimized: false,
      collapsed: false,
    }
    setPanels(prev => [...prev, newPanel])
    setShowTemplateModal(false)
  }, [])

  const handleDeletePanel = useCallback((id: string) => {
    if (confirm('Delete this panel?')) {
      setPanels(prev => prev.filter(p => p.id !== id))
      if (selectedPanel === id) {
        setSelectedPanel(null)
      }
    }
  }, [selectedPanel])

  const handleEditPanel = useCallback((id: string, updates: Partial<Panel>) => {
    setPanels(prev => prev.map(p => p.id === id ? { ...p, ...updates } : p))
  }, [])

  const handleToggleMinimize = useCallback((id: string) => {
    setPanels(prev => prev.map(p => p.id === id ? { ...p, minimized: !p.minimized } : p))
  }, [])

  const handleSave = useCallback(() => {
    onSave?.(panels)
  }, [panels, onSave])

  const handleExport = useCallback(() => {
    const data = JSON.stringify(panels, null, 2)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'dashboard.json'
    a.click()
    URL.revokeObjectURL(url)
    onExport?.(panels)
  }, [panels, onExport])

  const handleImport = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (event) => {
      try {
        const data = JSON.parse(event.target?.result as string)
        setPanels(data)
        onImport?.(data)
      } catch (err) {
        console.error('Failed to import dashboard:', err)
      }
    }
    reader.readAsText(file)
  }, [onImport])

  const handleClear = useCallback(() => {
    if (confirm('Clear all panels?')) {
      setPanels([])
    }
  }, [])

  const totalWidth = GRID_COLUMNS * 100
  const totalHeight = Math.max(400, panels.reduce((max, p) => Math.max(max, (p.y + p.height) * GRID_ROW_HEIGHT), 0))

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Dashboard Builder</h1>
          <p className="text-sm text-muted">Create and arrange dashboard panels</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 bg-tertiary rounded-lg p-1">
            <button
              onClick={() => setShowGrid(!showGrid)}
              className={`p-1.5 rounded ${showGrid ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="Toggle grid"
            >
              <GridIcon className="w-4 h-4" />
            </button>
            <button
              onClick={() => setZoom(z => Math.max(0.5, z - 0.1))}
              className="p-1.5 rounded text-muted hover:text-primary"
              title="Zoom out"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
                <line x1="8" y1="11" x2="14" y2="11" />
              </svg>
            </button>
            <span className="text-xs text-muted px-2">{Math.round(zoom * 100)}%</span>
            <button
              onClick={() => setZoom(z => Math.min(2, z + 0.1))}
              className="p-1.5 rounded text-muted hover:text-primary"
              title="Zoom in"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
                <line x1="11" y1="8" x2="11" y2="14" />
                <line x1="8" y1="11" x2="14" y2="11" />
              </svg>
            </button>
          </div>
          <button className="btn btn--ghost btn--sm" onClick={() => setShowTemplateModal(true)}>
            <PlusIcon className="w-4 h-4 mr-2" />
            Add Panel
          </button>
          <button className="btn btn--ghost btn--sm" onClick={handleExport}>
            <DownloadIcon className="w-4 h-4 mr-2" />
            Export
          </button>
          <label className="btn btn--ghost btn--sm cursor-pointer">
            <UploadIcon className="w-4 h-4 mr-2" />
            Import
            <input type="file" accept=".json" onChange={handleImport} className="hidden" />
          </label>
          <button className="btn btn--ghost btn--sm text-error" onClick={handleClear}>
            <TrashIcon className="w-4 h-4 mr-2" />
            Clear
          </button>
          <button className="btn btn--primary" onClick={handleSave}>
            <SaveIcon className="w-4 h-4 mr-2" />
            Save
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4">
        <div
          ref={containerRef}
          className="relative bg-secondary rounded-lg border border-border"
          style={{
            width: `${totalWidth * zoom}px`,
            height: `${totalHeight * zoom}px`,
            backgroundImage: showGrid ? 'linear-gradient(rgba(0,0,0,0.05) 1px, transparent 1px), linear-gradient(90deg, rgba(0,0,0,0.05) 1px, transparent 1px)' : 'none',
            backgroundSize: `${100 * zoom}px ${GRID_ROW_HEIGHT * zoom}px`,
          }}
        >
          {panels.map(panel => (
            <div
              key={panel.id}
              className={`absolute bg-primary border border-border rounded-lg shadow-sm hover:shadow-md transition-shadow ${selectedPanel === panel.id ? 'ring-2 ring-primary' : ''}`}
              style={{
                left: `${panel.x * 100 * zoom}px`,
                top: `${panel.y * GRID_ROW_HEIGHT * zoom}px`,
                width: `${panel.width * 100 * zoom}px`,
                height: `${panel.height * GRID_ROW_HEIGHT * zoom}px`,
              }}
              onClick={() => setSelectedPanel(panel.id)}
            >
              <DashboardPanel
                id={panel.id}
                title={panel.title}
                type={panel.type}
                query={panel.query}
                width={panel.width}
                height={panel.height}
                onEdit={handleEditPanel}
                onDelete={handleDeletePanel}
                onToggleMinimize={handleToggleMinimize}
                minimized={panel.minimized}
                collapsed={panel.collapsed}
                loading={false}
              />
            </div>
          ))}

          {panels.length === 0 && (
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="text-center">
                <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="mx-auto mb-4 text-muted">
                  <rect x="3" y="3" width="7" height="7" rx="1" />
                  <rect x="14" y="3" width="7" height="7" rx="1" />
                  <rect x="3" y="14" width="7" height="7" rx="1" />
                  <rect x="14" y="14" width="7" height="7" rx="1" />
                </svg>
                <h3 className="text-lg font-semibold mb-2">No panels yet</h3>
                <p className="text-sm text-muted mb-4">Add your first panel to start building your dashboard</p>
                <button className="btn btn--primary" onClick={() => setShowTemplateModal(true)}>
                  <PlusIcon className="w-4 h-4 mr-2" />
                  Add Panel
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Template Modal */}
      {showTemplateModal && (
        <div className="modal__overlay" onClick={() => setShowTemplateModal(false)}>
          <div className="modal">
            <div className="modal__header">
              <h2 className="modal__title">Add Panel</h2>
              <button className="modal__close" onClick={() => setShowTemplateModal(false)}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <div className="modal__body">
              <div className="grid gap-3 sm:grid-cols-2">
                {PANEL_TEMPLATES.map((template, i) => (
                  <button
                    key={i}
                    onClick={() => handleAddPanel(template)}
                    className="p-4 bg-secondary border border-border rounded-lg hover:bg-tertiary transition-colors text-left"
                  >
                    <div className="flex items-center gap-2 mb-2">
                      <span className="text-2xl">{template.type === 'timeseries' ? '📈' : template.type === 'stat' ? '📊' : template.type === 'table' ? '📋' : '🔢'}</span>
                      <span className="font-semibold">{template.title}</span>
                    </div>
                    <code className="text-xs font-mono text-primary bg-code px-2 py-1 rounded block truncate">{template.query}</code>
                    <p className="text-xs text-muted mt-2">{template.width}×{template.height} grid</p>
                  </button>
                ))}
              </div>
              <div className="mt-4 pt-4 border-t border-border">
                <h3 className="font-semibold mb-2">Custom Panel</h3>
                <button
                  onClick={() => handleAddPanel({ type: 'timeseries', title: 'Custom Panel', query: '', width: 6, height: 3 })}
                  className="btn btn--secondary w-full"
                >
                  <PlusIcon className="w-4 h-4 mr-2" />
                  Create Custom Panel
                </button>
              </div>
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
function SaveIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
}
function DownloadIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
}
function UploadIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
}
function GridIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
}
function LayoutIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
}
function SettingsIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 019 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 014.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 011 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001 1.51" /></svg>
}

export { DashboardBuilder }