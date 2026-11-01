import React, { useState, useCallback, useEffect, useRef, useMemo } from 'react'
import { DashboardPanel } from '../Dashboard/DashboardPanel'

interface DashboardGridProps {
  panels: Array<{
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
    data?: any
    loading?: boolean
    error?: string
  }>
  columns?: number
  rowHeight?: number
  gap?: number
  editable?: boolean
  onPanelEdit?: (id: string, updates: any) => void
  onPanelDelete?: (id: string) => void
  onPanelMove?: (id: string, x: number, y: number) => void
  onPanelResize?: (id: string, width: number, height: number) => void
  onAddPanel?: () => void
  className?: string
}

export function DashboardGrid({
  panels,
  columns = 12,
  rowHeight = 80,
  gap = 16,
  editable = false,
  onPanelEdit,
  onPanelDelete,
  onPanelMove,
  onPanelResize,
  onAddPanel,
  className = '',
}: DashboardGridProps) {
  const [draggingPanel, setDraggingPanel] = useState<string | null>(null)
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 })
  const [resizingPanel, setResizingPanel] = useState<string | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  const gridStyle = useMemo(() => ({
    display: 'grid',
    gridTemplateColumns: `repeat(${columns}, 1fr)`,
    gridAutoRows: `${rowHeight}px`,
    gap: `${gap}px`,
    position: 'relative' as const,
  }), [columns, rowHeight, gap])

  const handleDragStart = useCallback((e: React.DragEvent, panelId: string) => {
    if (!editable) return
    setDraggingPanel(panelId)
    const rect = (e.target as HTMLElement).getBoundingClientRect()
    setDragOffset({
      x: e.clientX - rect.left,
      y: e.clientY - rect.top,
    })
  }, [editable])

  const handleDragEnd = useCallback(() => {
    setDraggingPanel(null)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent, targetX: number, targetY: number) => {
    if (!editable || !draggingPanel) return
    e.preventDefault()
    onPanelMove?.(draggingPanel, targetX, targetY)
    setDraggingPanel(null)
  }, [editable, draggingPanel, onPanelMove])

  const handleResizeStart = useCallback((e: React.MouseEvent, panelId: string) => {
    if (!editable) return
    e.preventDefault()
    e.stopPropagation()
    setResizingPanel(panelId)
  }, [editable])

  const handleResizeEnd = useCallback(() => {
    setResizingPanel(null)
  }, [])

  useEffect(() => {
    if (resizingPanel) {
      const handleMouseMove = (e: MouseEvent) => {
        // Calculate new size based on mouse position
        const panel = panels.find(p => p.id === resizingPanel)
        if (!panel || !containerRef.current) return
        
        const containerRect = containerRef.current.getBoundingClientRect()
        const newWidth = Math.max(1, Math.round((e.clientX - containerRect.left) / (containerRect.width / columns)))
        const newHeight = Math.max(1, Math.round((e.clientY - containerRect.top) / rowHeight))
        
        onPanelResize?.(resizingPanel, Math.min(newWidth, columns), Math.min(newHeight, 12))
      }
      
      const handleMouseUp = () => {
        setResizingPanel(null)
      }
      
      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
      
      return () => {
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
      }
    }
  }, [resizingPanel, panels, columns, rowHeight, onPanelResize])

  const renderPanel = useCallback((panel: any) => {
    const gridArea = {
      gridColumn: `${panel.x + 1} / span ${panel.width}`,
      gridRow: `${panel.y + 1} / span ${panel.height}`,
    }

    return (
      <div
        key={panel.id}
        style={gridArea}
        className={`relative ${draggingPanel === panel.id ? 'opacity-50' : ''} ${resizingPanel === panel.id ? 'z-10' : ''}`}
        draggable={editable}
        onDragStart={(e) => handleDragStart(e, panel.id)}
        onDragEnd={handleDragEnd}
        onDragOver={(e) => e.preventDefault()}
        onDrop={(e) => handleDrop(e, panel.x, panel.y)}
      >
        <DashboardPanel
          id={panel.id}
          title={panel.title}
          type={panel.type}
          query={panel.query}
          width={panel.width}
          height={panel.height}
          onEdit={onPanelEdit || (() => {})}
          onDelete={onPanelDelete || (() => {})}
          onToggleMinimize={() => {}}
          minimized={panel.minimized}
          collapsed={panel.collapsed}
          loading={panel.loading}
          data={panel.data}
          error={panel.error}
        />
        {editable && (
          <div
            className="absolute bottom-0 right-0 w-4 h-4 cursor-se-resize opacity-0 group-hover:opacity-100 transition-opacity"
            onMouseDown={(e) => handleResizeStart(e, panel.id)}
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-muted">
              <polyline points="15 3 21 3 21 9" />
              <polyline points="9 21 3 21 3 15" />
              <line x1="21" y1="3" x2="14" y2="10" />
              <line x1="3" y1="21" x2="10" y2="14" />
            </svg>
          </div>
        )}
      </div>
    )
  }, [editable, draggingPanel, resizingPanel, handleDragStart, handleDragEnd, handleDrop, handleResizeStart, onPanelEdit, onPanelDelete])

  const totalRows = useMemo(() => {
    return Math.max(...panels.map(p => p.y + p.height), 1)
  }, [panels])

  return (
    <div className={`dashboard-grid ${className}`}>
      <div
        ref={containerRef}
        style={{
          ...gridStyle,
          minHeight: `${totalRows * rowHeight + (totalRows - 1) * gap}px`,
        }}
      >
        {panels.map(renderPanel)}
      </div>

      {editable && onAddPanel && (
        <div className="mt-4 flex justify-center">
          <button
            onClick={onAddPanel}
            className="btn btn--secondary"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
              <line x1="12" y1="5" x2="12" y2="19" />
              <line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Add Panel
          </button>
        </div>
      )}

      {panels.length === 0 && (
        <div className="empty-state">
          <svg className="empty-state__icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
            <rect x="3" y="3" width="7" height="7" rx="1" />
            <rect x="14" y="3" width="7" height="7" rx="1" />
            <rect x="3" y="14" width="7" height="7" rx="1" />
            <rect x="14" y="14" width="7" height="7" rx="1" />
          </svg>
          <h3 className="empty-state__title">No panels</h3>
          <p className="empty-state__description">Add panels to build your dashboard</p>
          {editable && onAddPanel && (
            <button className="btn btn--primary mt-4" onClick={onAddPanel}>
              Add First Panel
            </button>
          )}
        </div>
      )}
    </div>
  )
}

export { DashboardGrid }