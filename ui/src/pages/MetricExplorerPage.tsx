import React, { useState, useCallback, useEffect, useMemo } from 'react'
import { FilterBar } from '../Common/FilterBar'
import { DataTable } from '../Common/DataTable'
import { MetricCard } from '../Common/MetricCard'
import { ChartPanel } from '../Common/ChartPanel'
import { TimeRangePicker } from '../Common/TimeRangePicker'

interface Metric {
  name: string
  type: 'counter' | 'gauge' | 'histogram' | 'summary'
  unit: string
  help: string
  labels: string[]
  cardinality: number
  lastSeen: number
  sampleCount: number
}

interface MetricExplorerProps {
  onMetricSelect?: (metric: Metric) => void
  onQueryCreate?: (query: string) => void
}

const MOCK_METRICS: Metric[] = [
  { name: 'cpu_usage', type: 'gauge', unit: 'percent', help: 'CPU usage percentage', labels: ['host', 'instance', 'job'], cardinality: 150, lastSeen: Date.now(), sampleCount: 1000000 },
  { name: 'memory_usage', type: 'gauge', unit: 'bytes', help: 'Memory usage in bytes', labels: ['host', 'instance', 'job'], cardinality: 150, lastSeen: Date.now(), sampleCount: 1000000 },
  { name: 'disk_usage', type: 'gauge', unit: 'bytes', help: 'Disk usage in bytes', labels: ['host', 'instance', 'mount'], cardinality: 300, lastSeen: Date.now(), sampleCount: 500000 },
  { name: 'network_rx_bytes', type: 'counter', unit: 'bytes', help: 'Network received bytes', labels: ['host', 'interface'], cardinality: 50, lastSeen: Date.now(), sampleCount: 2000000 },
  { name: 'network_tx_bytes', type: 'counter', unit: 'bytes', help: 'Network transmitted bytes', labels: ['host', 'interface'], cardinality: 50, lastSeen: Date.now(), sampleCount: 2000000 },
  { name: 'http_requests_total', type: 'counter', unit: 'requests', help: 'Total HTTP requests', labels: ['method', 'status', 'path'], cardinality: 500, lastSeen: Date.now(), sampleCount: 5000000 },
  { name: 'http_request_duration_seconds', type: 'histogram', unit: 'seconds', help: 'HTTP request duration', labels: ['method', 'status', 'path'], cardinality: 500, lastSeen: Date.now(), sampleCount: 5000000 },
  { name: 'up', type: 'gauge', unit: 'boolean', help: 'Target up status', labels: ['instance', 'job'], cardinality: 100, lastSeen: Date.now(), sampleCount: 100000 },
  { name: 'process_cpu_seconds_total', type: 'counter', unit: 'seconds', help: 'Total process CPU time', labels: ['instance', 'job'], cardinality: 100, lastSeen: Date.now(), sampleCount: 100000 },
  { name: 'process_resident_memory_bytes', type: 'gauge', unit: 'bytes', help: 'Resident memory size', labels: ['instance', 'job'], cardinality: 100, lastSeen: Date.now(), sampleCount: 100000 },
  { name: 'go_goroutines', type: 'gauge', unit: 'goroutines', help: 'Number of goroutines', labels: ['instance', 'job'], cardinality: 100, lastSeen: Date.now(), sampleCount: 100000 },
  { name: 'go_memstats_alloc_bytes', type: 'gauge', unit: 'bytes', help: 'Bytes allocated', labels: ['instance', 'job'], cardinality: 100, lastSeen: Date.now(), sampleCount: 100000 },
]

export function MetricExplorer({ onMetricSelect, onQueryCreate }: MetricExplorerProps) {
  const [metrics, setMetrics] = useState<Metric[]>(MOCK_METRICS)
  const [selectedMetric, setSelectedMetric] = useState<Metric | null>(null)
  const [viewMode, setViewMode] = useState<'grid' | 'list' | 'detail'>('grid')
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState<string[]>([])
  const [sortBy, setSortBy] = useState<'name' | 'cardinality' | 'sampleCount' | 'lastSeen'>('name')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('asc')
  const [timeRange, setTimeRange] = useState<[Date, Date]>([new Date(Date.now() - 3600000), new Date()])
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

  const handleMetricClick = useCallback((metric: Metric) => {
    setSelectedMetric(metric)
    setViewMode('detail')
    onMetricSelect?.(metric)
  }, [onMetricSelect])

  const handleCreateQuery = useCallback((metric: Metric) => {
    const query = `${metric.name}{job="node"}`
    onQueryCreate?.(query)
  }, [onQueryCreate])

  const handleBack = useCallback(() => {
    setSelectedMetric(null)
    setViewMode('grid')
  }, [])

  const typeColors: Record<string, string> = {
    counter: 'badge--info',
    gauge: 'badge--success',
    histogram: 'badge--warning',
    summary: 'badge--neutral',
  }

  if (viewMode === 'detail' && selectedMetric) {
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
            <MetricCard title="Type" value={selectedMetric.type} color="info" />
            <MetricCard title="Unit" value={selectedMetric.unit} color="primary" />
            <MetricCard title="Cardinality" value={selectedMetric.cardinality} color="warning" />
            <MetricCard title="Sample Count" value={selectedMetric.sampleCount} color="success" />
          </div>

          <div className="card mb-6">
            <div className="card__header">
              <h3 className="font-semibold">Labels</h3>
            </div>
            <div className="card__body">
              <div className="flex flex-wrap gap-2">
                {selectedMetric.labels.map(label => (
                  <span key={label} className="badge badge--neutral">{label}</span>
                ))}
              </div>
            </div>
          </div>

          <div className="card mb-6">
            <div className="card__header flex items-center justify-between">
              <h3 className="font-semibold">Query Examples</h3>
            </div>
            <div className="card__body space-y-2">
              {[
                `${selectedMetric.name}`,
                `rate(${selectedMetric.name}[5m])`,
                `sum by (host) (rate(${selectedMetric.name}[5m]))`,
                `topk(10, rate(${selectedMetric.name}[5m]))`,
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

          <ChartPanel
            title={`${selectedMetric.name} over time`}
            type="line"
            data={{
              labels: Array.from({ length: 24 }, (_, i) => `${i}:00`),
              datasets: [{
                label: selectedMetric.name,
                data: Array.from({ length: 24 }, () => Math.random() * 100),
              }],
            }}
            height={300}
            timeRange={timeRange}
            onTimeRangeChange={setTimeRange}
          />
        </div>
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
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`btn btn--sm ${showFilters ? 'btn--primary' : 'btn--ghost'}`}
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
            </svg>
            Filters
          </button>
        </div>
      </div>

      {showFilters && (
        <div className="bg-secondary border-b border-border p-4 animate-slide-in">
          <div className="flex flex-wrap gap-4">
            <div>
              <label className="label">Type</label>
              <div className="flex gap-2">
                {['counter', 'gauge', 'histogram', 'summary'].map(t => (
                  <label key={t} className={`badge ${typeFilter.includes(t) ? 'badge--info' : 'badge--neutral'} cursor-pointer`}>
                    <input
                      type="checkbox"
                      className="sr-only"
                      checked={typeFilter.includes(t)}
                      onChange={() => setTypeFilter(prev => prev.includes(t) ? prev.filter(x => x !== t) : [...prev, t])}
                    />
                    {t}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Sort By</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as any)}
                className="input input--sm"
              >
                <option value="name">Name</option>
                <option value="cardinality">Cardinality</option>
                <option value="sampleCount">Sample Count</option>
                <option value="lastSeen">Last Seen</option>
              </select>
            </div>
            <div>
              <label className="label">Direction</label>
              <select
                value={sortDir}
                onChange={(e) => setSortDir(e.target.value as any)}
                className="input input--sm"
              >
                <option value="asc">Ascending</option>
                <option value="desc">Descending</option>
              </select>
            </div>
          </div>
        </div>
      )}

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
                  <span className={`badge ${typeColors[metric.type]}`}>{metric.type}</span>
                </div>
                <p className="text-xs text-muted mb-3 line-clamp-2">{metric.help}</p>
                <div className="flex items-center justify-between text-xs text-muted">
                  <span>{metric.cardinality} series</span>
                  <span>{metric.sampleCount.toLocaleString()} samples</span>
                </div>
                <div className="mt-2 flex flex-wrap gap-1">
                  {metric.labels.slice(0, 3).map(label => (
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
          <DataTable
            data={filteredMetrics}
            columns={[
              { key: 'name', header: 'Metric Name', accessor: (m: Metric) => m.name, sortable: true },
              { key: 'type', header: 'Type', accessor: (m: Metric) => m.type, sortable: true, render: (v: string) => <span className={`badge ${typeColors[v]}`}>{v}</span> },
              { key: 'unit', header: 'Unit', accessor: (m: Metric) => m.unit },
              { key: 'cardinality', header: 'Cardinality', accessor: (m: Metric) => m.cardinality, sortable: true, align: 'right' },
              { key: 'sampleCount', header: 'Samples', accessor: (m: Metric) => m.sampleCount, sortable: true, align: 'right', render: (v: number) => v.toLocaleString() },
              { key: 'labels', header: 'Labels', accessor: (m: Metric) => m.labels.join(', '), render: (v: string, m: Metric) => (
                <div className="flex flex-wrap gap-1">
                  {m.labels.slice(0, 3).map(l => <span key={l} className="badge badge--neutral text-xs">{l}</span>)}
                </div>
              )},
              { key: 'actions', header: 'Actions', accessor: (m: Metric) => m, render: (_: any, m: Metric) => (
                <div className="flex gap-1">
                  <button className="btn btn--ghost btn--sm" onClick={() => handleMetricClick(m)}>View</button>
                  <button className="btn btn--ghost btn--sm" onClick={() => handleCreateQuery(m)}>Query</button>
                </div>
              )},
            ]}
            pageSize={20}
            emptyMessage="No metrics found"
          />
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