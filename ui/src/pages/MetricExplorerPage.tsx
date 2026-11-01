import React, { useState, useCallback, useEffect, useRef, useMemo } from 'react'

interface MetricExplorerProps {
  metrics: Array<{
    name: string
    type: 'counter' | 'gauge' | 'histogram' | 'summary'
    unit: string
    help: string
    labels: string[]
    cardinality: number
    lastSeen: number
    sampleCount: number
  }>
  onMetricSelect?: (metric: any) => void
  onQueryCreate?: (query: string) => void
  loading?: boolean
}

const TYPE_COLORS: Record<string, string> = {
  counter: 'badge--info',
  gauge: 'badge--success',
  histogram: 'badge--warning',
  summary: 'badge--neutral',
}

export function MetricExplorer({ metrics, onMetricSelect, onQueryCreate, loading = false }: MetricExplorerProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState<string[]>([])
  const [sortBy, setSortBy] = useState<'name' | 'cardinality' | 'sampleCount' | 'lastSeen'>('name')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('asc')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [selectedMetric, setSelectedMetric] = useState<any>(null)
  const [showFilters, setShowFilters] = useState(false)

  const filteredMetrics = useMemo(() => {
    let result = [...metrics]

    if (searchQuery) {
      const lower = searchQuery.toLowerCase()
      result = result.filter(m =>
        m.name.toLowerCase().includes(lower) ||
        m.help.toLowerCase().includes(lower) ||
        m.labels.some(l => l.toLowerCase().includes(lower))
      )
    }

    if (typeFilter.length > 0) {
      result = result.filter(m => typeFilter.includes(m.type))
    }

    result.sort((a, b) => {
      let comparison = 0
      switch (sortBy) {
        case 'name':
          comparison = a.name.localeCompare(b.name)
          break
        case 'cardinality':
          comparison = a.cardinality - b.cardinality
          break
        case 'sampleCount':
          comparison = a.sampleCount - b.sampleCount
          break
        case 'lastSeen':
          comparison = a.lastSeen - b.lastSeen
          break
      }
      return sortDir === 'asc' ? comparison : -comparison
    })

    return result
  }, [metrics, searchQuery, typeFilter, sortBy, sortDir])

  const handleMetricClick = useCallback((metric: any) => {
    setSelectedMetric(metric)
    onMetricSelect?.(metric)
  }, [onMetricSelect])

  const handleCreateQuery = useCallback((metric: any) => {
    const query = `${metric.name}{job="node"}`
    onQueryCreate?.(query)
  }, [onQueryCreate])

  const handleBack = useCallback(() => {
    setSelectedMetric(null)
  }, [])

  if (selectedMetric) {
    return (
      <div className="h-full flex flex-col bg-primary">
        <div className="flex items-center gap-3 p-4 border-b border-border bg-secondary">
          <button className="btn btn--ghost btn--sm" onClick={handleBack}>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="15 18 9 12 15 6" />
            </svg>
            Back
          </button>
          <div>
            <h1 className="text-2xl font-bold">{selectedMetric.name}</h1>
            <p className="text-sm text-muted">{selectedMetric.help}</p>
          </div>
        </div>

        <div className="flex-1 overflow-auto p-4">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6">
            <div className="card p-4">
              <p className="text-xs font-medium text-muted uppercase tracking-wide">Type</p>
              <p className="text-2xl font-bold text-primary mt-1">{selectedMetric.type}</p>
            </div>
            <div className="card p-4">
              <p className="text-xs font-medium text-muted uppercase tracking-wide">Unit</p>
              <p className="text-2xl font-bold text-primary mt-1">{selectedMetric.unit}</p>
            </div>
            <div className="card p-4">
              <p className="text-xs font-medium text-muted uppercase tracking-wide">Cardinality</p>
              <p className="text-2xl font-bold text-primary mt-1">{selectedMetric.cardinality}</p>
            </div>
            <div className="card p-4">
              <p className="text-xs font-medium text-muted uppercase tracking-wide">Samples</p>
              <p className="text-2xl font-bold text-primary mt-1">{selectedMetric.sampleCount.toLocaleString()}</p>
            </div>
          </div>

          <div className="card mb-6">
            <div className="card__header">
              <h3 className="font-semibold">Labels</h3>
            </div>
            <div className="card__body">
              <div className="flex flex-wrap gap-2">
                {selectedMetric.labels.map((label: string) => (
                  <span key={label} className="badge badge--neutral">{label}</span>
                ))}
              </div>
            </div>
          </div>

          <div className="card mb-6">
            <div className="card__header">
              <h3 className="font-semibold">Query Examples</h3>
            </div>
            <div className="card__body space-y-2">
              {[
                `${selectedMetric.name}`,
                `rate(${selectedMetric.name}[5m])`,
                `sum by (host) (rate(${selectedMetric.name}[5m]))`,
                `topk(10, rate(${selectedMetric.name}[5m]))`,
                `histogram_quantile(0.95, rate(${selectedMetric.name}_bucket[5m]))`,
              ].map((q, i) => (
                <div key={i} className="flex items-center justify-between p-2 bg-code rounded">
                  <code className="text-sm font-mono text-primary">{q}</code>
                  <button
                    className="btn btn--ghost btn--sm"
                    onClick={() => onQueryCreate?.(q)}
                  >
                    Use
                  </button>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="loading-spinner loading-spinner--lg" />
      </div>
    )
  }

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Metric Explorer</h1>
          <p className="text-sm text-muted">Discover and explore available metrics</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative">
            <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              placeholder="Search metrics..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-9 w-64"
            />
          </div>
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            {(['counter', 'gauge', 'histogram', 'summary'] as const).map(t => (
              <button
                key={t}
                onClick={() => setTypeFilter(prev => prev.includes(t) ? prev.filter(x => x !== t) : [...prev, t])}
                className={`px-2 py-1 rounded text-xs font-medium transition-colors ${typeFilter.includes(t) ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              >
                {t}
              </button>
            ))}
          </div>
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-1.5 rounded ${viewMode === 'grid' ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="Grid view"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="3" width="7" height="7" />
                <rect x="14" y="3" width="7" height="7" />
                <rect x="3" y="14" width="7" height="7" />
                <rect x="14" y="14" width="7" height="7" />
              </svg>
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-1.5 rounded ${viewMode === 'list' ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="List view"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="8" y1="6" x2="21" y2="6" />
                <line x1="8" y1="12" x2="21" y2="12" />
                <line x1="8" y1="18" x2="21" y2="18" />
                <line x1="3" y1="6" x2="3.01" y2="6" />
                <line x1="3" y1="12" x2="3.01" y2="12" />
                <line x1="3" y1="18" x2="3.01" y2="18" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4">
        {viewMode === 'grid' ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {filteredMetrics.map(metric => (
              <div
                key={metric.name}
                onClick={() => handleMetricClick(metric)}
                className="card p-4 cursor-pointer hover:shadow-md transition-shadow group"
              >
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-semibold text-sm truncate flex-1">{metric.name}</h3>
                  <span className={`badge ${TYPE_COLORS[metric.type]}`}>{metric.type}</span>
                </div>
                <p className="text-xs text-muted mb-3 line-clamp-2">{metric.help}</p>
                <div className="flex items-center justify-between text-xs text-muted">
                  <span>{metric.cardinality} series</span>
                  <span>{metric.sampleCount.toLocaleString()} samples</span>
                </div>
                <div className="mt-2 flex flex-wrap gap-1">
                  {metric.labels.slice(0, 3).map((label: string) => (
                    <span key={label} className="badge badge--neutral text-xs">{label}</span>
                  ))}
                  {metric.labels.length > 3 && (
                    <span className="badge badge--neutral text-xs">+{metric.labels.length - 3}</span>
                  )}
                </div>
                <div className="mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={(e) => { e.stopPropagation(); handleCreateQuery(metric); }}
                    className="btn btn--primary btn--sm w-full"
                  >
                    Create Query
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Metric Name</th>
                  <th>Type</th>
                  <th>Unit</th>
                  <th>Cardinality</th>
                  <th>Samples</th>
                  <th>Labels</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredMetrics.map(metric => (
                  <tr key={metric.name}>
                    <td className="font-mono text-sm">{metric.name}</td>
                    <td><span className={`badge ${TYPE_COLORS[metric.type]}`}>{metric.type}</span></td>
                    <td className="text-sm text-muted">{metric.unit}</td>
                    <td className="text-right">{metric.cardinality}</td>
                    <td className="text-right">{metric.sampleCount.toLocaleString()}</td>
                    <td>
                      <div className="flex flex-wrap gap-1">
                        {metric.labels.slice(0, 3).map((l: string) => (
                          <span key={l} className="badge badge--neutral text-xs">{l}</span>
                        ))}
                      </div>
                    </td>
                    <td>
                      <div className="flex gap-1">
                        <button className="btn btn--ghost btn--sm" onClick={() => handleMetricClick(metric)}>View</button>
                        <button className="btn btn--ghost btn--sm" onClick={() => handleCreateQuery(metric)}>Query</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {filteredMetrics.length === 0 && (
          <div className="empty-state">
            <svg className="empty-state__icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M3 3v18h18" />
              <path d="M18 17V9" />
              <path d="M13 17V5" />
              <path d="M8 17v-3" />
            </svg>
            <h3 className="empty-state__title">No metrics found</h3>
            <p className="empty-state__description">Try adjusting your search or filters</p>
          </div>
        )}
      </div>
    </div>
  )
}

export { MetricExplorer }