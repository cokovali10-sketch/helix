import React, { useState, useCallback, useEffect, useMemo } from 'react'

interface QueryHistoryProps {
  history: Array<{
    id: string
    expression: string
    timestamp: number
    tenant?: string
    resultType: string
    rowCount: number
    executionTimeMs: number
    error?: string
  }>
  onQuerySelect: (expression: string) => void
  onClear: () => void
  onDelete: (id: string) => void
  maxItems?: number
}

export function QueryHistory({ history, onQuerySelect, onClear, onDelete, maxItems = 50 }: QueryHistoryProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [filter, setFilter] = useState<'all' | 'success' | 'error'>('all')
  const [sortBy, setSortBy] = useState<'time' | 'duration' | 'rows'>('time')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')
  const [selectedItems, setSelectedItems] = useState<Set<string>>(new Set())

  const filteredHistory = useMemo(() => {
    let result = [...history]

    if (searchQuery) {
      const lower = searchQuery.toLowerCase()
      result = result.filter(h => h.expression.toLowerCase().includes(lower))
    }

    if (filter === 'success') {
      result = result.filter(h => !h.error)
    } else if (filter === 'error') {
      result = result.filter(h => h.error)
    }

    result.sort((a, b) => {
      let comparison = 0
      switch (sortBy) {
        case 'time':
          comparison = a.timestamp - b.timestamp
          break
        case 'duration':
          comparison = a.executionTimeMs - b.executionTimeMs
          break
        case 'rows':
          comparison = a.rowCount - b.rowCount
          break
      }
      return sortDir === 'asc' ? comparison : -comparison
    })

    return result.slice(0, maxItems)
  }, [history, searchQuery, filter, sortBy, sortDir, maxItems])

  const handleSelectAll = useCallback(() => {
    if (selectedItems.size === filteredHistory.length) {
      setSelectedItems(new Set())
    } else {
      setSelectedItems(new Set(filteredHistory.map(h => h.id)))
    }
  }, [selectedItems, filteredHistory])

  const handleSelectItem = useCallback((id: string) => {
    const newSelected = new Set(selectedItems)
    if (newSelected.has(id)) {
      newSelected.delete(id)
    } else {
      newSelected.add(id)
    }
    setSelectedItems(newSelected)
  }, [selectedItems])

  const handleDeleteSelected = useCallback(() => {
    if (selectedItems.size === 0) return
    if (confirm(`Delete ${selectedItems.size} selected queries?`)) {
      selectedItems.forEach(id => onDelete(id))
      setSelectedItems(new Set())
    }
  }, [selectedItems, onDelete])

  const handleExport = useCallback(() => {
    const data = JSON.stringify(filteredHistory, null, 2)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'query-history.json'
    a.click()
    URL.revokeObjectURL(url)
  }, [filteredHistory])

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms.toFixed(0)}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(2)}s`
    return `${(ms / 60000).toFixed(2)}m`
  }

  const formatTime = (ts: number) => {
    const date = new Date(ts)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    
    if (diff < 60000) return 'Just now'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
    return date.toLocaleDateString()
  }

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Query History</h1>
          <p className="text-sm text-muted">Browse and manage your query history</p>
        </div>
        <div className="flex items-center gap-3">
          <input
            type="text"
            placeholder="Search queries..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="input input--sm w-48"
          />
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            {(['all', 'success', 'error'] as const).map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${filter === f ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
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
            <option value="time">Time</option>
            <option value="duration">Duration</option>
            <option value="rows">Rows</option>
          </select>
          <button
            onClick={() => setSortDir(d => d === 'asc' ? 'desc' : 'asc')}
            className="btn btn--ghost btn--sm"
          >
            {sortDir === 'asc' ? '↑' : '↓'}
          </button>
          <button className="btn btn--ghost btn--sm" onClick={handleExport}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
              <polyline points="7 10 12 15 17 10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
            Export
          </button>
          <button className="btn btn--ghost btn--sm text-error" onClick={onClear}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <polyline points="3 6 5 6 21 6" />
              <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
            </svg>
            Clear All
          </button>
        </div>
      </div>

      {selectedItems.size > 0 && (
        <div className="flex items-center gap-3 p-3 bg-primary-light border-b border-border">
          <span className="text-sm font-medium">{selectedItems.size} selected</span>
          <button className="btn btn--ghost btn--sm text-error" onClick={handleDeleteSelected}>
            Delete Selected
          </button>
          <button className="btn btn--ghost btn--sm" onClick={() => setSelectedItems(new Set())}>
            Clear Selection
          </button>
        </div>
      )}

      <div className="flex-1 overflow-auto">
        {filteredHistory.length === 0 ? (
          <div className="empty-state">
            <svg className="empty-state__icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
            <h3 className="empty-state__title">No query history</h3>
            <p className="empty-state__description">Your executed queries will appear here</p>
          </div>
        ) : (
          <div className="divide-y divide-border">
            {filteredHistory.map(item => (
              <div
                key={item.id}
                className={`p-4 hover:bg-tertiary transition-colors ${selectedItems.has(item.id) ? 'bg-primary-light' : ''}`}
              >
                <div className="flex items-start gap-3">
                  <input
                    type="checkbox"
                    checked={selectedItems.has(item.id)}
                    onChange={() => handleSelectItem(item.id)}
                    className="mt-1 w-4 h-4 rounded border-border"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2 mb-1">
                      <button
                        onClick={() => onQuerySelect(item.expression)}
                        className="text-left flex-1 min-w-0"
                      >
                        <code className="text-sm font-mono text-primary hover:underline block truncate">
                          {item.expression}
                        </code>
                      </button>
                      <span className={`badge ${item.error ? 'badge--error' : 'badge--success'} shrink-0`}>
                        {item.error ? 'Error' : item.resultType}
                      </span>
                    </div>
                    <div className="flex items-center gap-3 text-xs text-muted flex-wrap">
                      <span className="flex items-center gap-1">
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <circle cx="12" cy="12" r="10" />
                          <polyline points="12 6 12 12 16 14" />
                        </svg>
                        {formatTime(item.timestamp)}
                      </span>
                      <span>•</span>
                      <span>{item.rowCount} rows</span>
                      <span>•</span>
                      <span>{formatDuration(item.executionTimeMs)}</span>
                      {item.tenant && (
                        <>
                          <span>•</span>
                          <span>{item.tenant}</span>
                        </>
                      )}
                    </div>
                    {item.error && (
                      <p className="text-xs text-error mt-1 truncate">{item.error}</p>
                    )}
                  </div>
                  <div className="flex items-center gap-1 shrink-0">
                    <button
                      onClick={() => onQuerySelect(item.expression)}
                      className="p-1.5 rounded hover:bg-tertiary text-muted hover:text-primary transition-colors"
                      title="Run query"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polygon points="5 3 19 12 5 21 5 3" />
                      </svg>
                    </button>
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText(item.expression)
                      }}
                      className="p-1.5 rounded hover:bg-tertiary text-muted hover:text-primary transition-colors"
                      title="Copy"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                        <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
                      </svg>
                    </button>
                    <button
                      onClick={() => onDelete(item.id)}
                      className="p-1.5 rounded hover:bg-error-light text-muted hover:text-error transition-colors"
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
            ))}
          </div>
        )}
      </div>

      <div className="p-3 border-t border-border bg-secondary flex items-center justify-between text-sm text-muted">
        <span>{filteredHistory.length} of {history.length} queries</span>
        <span>Showing up to {maxItems} items</span>
      </div>
    </div>
  )
}

export { QueryHistory }