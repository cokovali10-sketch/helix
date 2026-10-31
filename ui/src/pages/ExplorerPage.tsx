import React, { useState, useCallback, useEffect } from 'react'
import { useQuery } from '../../context/QueryContext'
import { useToast } from '../../components/Common/ToastContainer'
import { SearchIcon, FilterIcon, ChevronDownIcon, ChevronRightIcon, DownloadIcon, EyeIcon } from '../Layout/Sidebar'

interface ExplorerPageProps {}

export function ExplorerPage() {
  const { executeQuery } = useQuery()
  const { addToast } = useToast()
  const [events, setEvents] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState({
    query: '',
    sources: [] as string[],
    tags: [] as string[],
    severity: [] as string[],
    timeRange: 3600000,
  })
  const [showFilters, setShowFilters] = useState(false)
  const [selectedEvent, setSelectedEvent] = useState<any>(null)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(50)
  const [totalEvents, setTotalEvents] = useState(0)

  const availableSources = ['soc-agent-01', 'soc-agent-02', 'web-server-01', 'db-server-01', 'api-gateway-01']
  const availableTags = ['production', 'staging', 'critical', 'warning', 'info', 'security', 'network', 'database']
  const availableSeverities = ['critical', 'warning', 'info', 'debug']

  const handleSearch = useCallback(async () => {
    setLoading(true)
    try {
      // In real implementation, this would call the API
      // const response = await fetch('/api/v1/events', { ... })
      // For now, generate mock data
      const mockEvents = Array.from({ length: 100 }, (_, i) => ({
        id: `evt-${i}`,
        timestamp: Date.now() - Math.random() * 86400000,
        source: availableSources[Math.floor(Math.random() * availableSources.length)],
        severity: availableSeverities[Math.floor(Math.random() * availableSeverities.length)],
        message: `Event message ${i}: ${['Connection established', 'Authentication failed', 'High memory usage', 'Disk space low', 'Network latency spike'][Math.floor(Math.random() * 5)]}`,
        tags: availableTags.slice(0, Math.floor(Math.random() * 3) + 1),
        metrics: { cpu: Math.random() * 100, memory: Math.random() * 100 },
      }))
      setEvents(mockEvents)
      setTotalEvents(mockEvents.length)
    } catch (err) {
      console.error('Search failed:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    handleSearch()
  }, [handleSearch])

  const filteredEvents = events.filter(e => {
    if (filters.query && !e.message.toLowerCase().includes(filters.query.toLowerCase())) return false
    if (filters.sources.length && !filters.sources.includes(e.source)) return false
    if (filters.tags.length && !filters.tags.some(t => e.tags.includes(t))) return false
    if (filters.severity.length && !filters.severity.includes(e.severity)) return false
    return true
  })

  const paginatedEvents = filteredEvents.slice((page - 1) * pageSize, page * pageSize)
  const totalPages = Math.ceil(filteredEvents.length / pageSize)

  const TIME_RANGES = [
    { label: 'Last 5 minutes', value: 5 * 60 * 1000 },
    { label: 'Last 15 minutes', value: 15 * 60 * 1000 },
    { label: 'Last 1 hour', value: 60 * 60 * 1000 },
    { label: 'Last 6 hours', value: 6 * 60 * 60 * 1000 },
    { label: 'Last 24 hours', value: 24 * 60 * 60 * 1000 },
  ]

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center gap-3 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Event Explorer</h1>
          <p className="text-sm text-muted">Search and explore raw events</p>
        </div>
        <div className="flex-1 max-w-xl mx-4">
          <div className="relative">
            <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted" />
            <input
              type="text"
              placeholder="Search events..."
              value={filters.query}
              onChange={e => setFilters(f => ({ ...f, query: e.target.value }))}
              className="input pl-10 w-full"
              placeholder="Search events..."
            />
          </div>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={filters.timeRange}
            onChange={e => setFilters(f => ({ ...f, timeRange: Number(e.target.value) }))}
            className="input input--sm w-auto hidden sm:block"
          >
            {TIME_RANGES.map(t => <option key={t.value} value={t.value}>{t.label}</option>)}
          </select>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`btn btn--sm ${showFilters ? 'btn--primary' : 'btn--ghost'}`}
          >
            <FilterIcon className="w-4 h-4 mr-1" />
            <span className="hidden sm:inline">Filters</span>
            {Object.values(filters).flat().length > 0 && (
              <span className="badge badge--info ml-1">
                {Object.values(filters).flat().length}
              </span>
            )}
          </button>
          <button className="btn btn--ghost btn--sm" title="Export">
            <DownloadIcon className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Filter Panel */}
      {showFilters && (
        <div className="bg-secondary border-b border-border animate-slide-in">
          <div className="p-4 border-b border-border flex items-center justify-between">
            <h3 className="font-semibold">Filters</h3>
            <button onClick={() => setShowFilters(false)} className="text-muted hover:text-primary">Close</button>
          </div>
          <div className="p-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label className="label">Sources</label>
              <div className="flex flex-wrap gap-2">
                {availableSources.map(s => (
                  <label key={s} className={`badge ${filters.sources.includes(s) ? 'badge--info' : 'badge--neutral'} cursor-pointer transition-colors`}>
                    <input type="checkbox" className="mr-1" checked={filters.sources.includes(s)} onChange={e => setFilters(f => ({ ...f, sources: e.target.checked ? [...f.sources, s] : f.sources.filter(x => x !== s) }))} />
                    {s}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Tags</label>
              <div className="flex flex-wrap gap-2">
                {availableTags.map(t => (
                  <label key={t} className={`badge ${filters.tags.includes(t) ? 'badge--warning' : 'badge--neutral'} cursor-pointer transition-colors`}>
                    <input type="checkbox" className="mr-1" checked={filters.tags.includes(t)} onChange={e => setFilters(f => ({ ...f, tags: e.target.checked ? [...f.tags, t] : f.tags.filter(x => x !== t) }))} />
                    {t}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Severity</label>
              <div className="flex gap-2">
                {availableSeverities.map(s => (
                  <label key={s} className={`badge ${filters.severity.includes(s) ? `badge--${s === 'critical' ? 'error' : s === 'warning' ? 'warning' : 'info'}` : 'badge--neutral'} cursor-pointer transition-colors`}>
                    <input type="checkbox" className="mr-1" checked={filters.severity.includes(s)} onChange={e => setFilters(f => ({ ...f, severity: e.target.checked ? [...f.severity, s] : f.severity.filter(x => x !== s) }))} />
                    {s}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Time Range</label>
              <select value={filters.timeRange} onChange={e => setFilters(f => ({ ...f, timeRange: Number(e.target.value) }))} className="input">
                {TIME_RANGES.map(t => <option key={t.value} value={t.value}>{t.label}</option>)}
              </select>
            </div>
          </div>
        </div>
      )}

      <div className="flex-1 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="loading-spinner loading-spinner--lg" />
          </div>
        ) : (
          <>
            <div className="px-4 py-2 border-b border-border bg-secondary flex items-center justify-between text-sm text-muted">
              <span>{filteredEvents.length} of {totalEvents} events</span>
              <div className="flex items-center gap-2">
                <select value={page} onChange={e => setPage(Number(e.target.value))} className="input input--sm w-auto">
                  {Array.from({ length: totalPages }, (_, i) => <option key={i + 1} value={i + 1}>Page {i + 1}</option>)}
                </select>
                <span>of {totalPages} pages</span>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="table w-full">
                <thead>
                  <tr>
                    <th style={{ width: '150px' }}>Time</th>
                    <th style={{ width: '100px' }}>Severity</th>
                    <th style={{ width: '150px' }}>Source</th>
                    <th>Message</th>
                    <th style={{ width: '150px' }}>Tags</th>
                    <th style={{ width: '80px' }}>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {paginatedEvents.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="text-center py-12 text-muted">No events found</td>
                    </tr>
                  ) : (
                    paginatedEvents.map((event, i) => (
                      <tr key={i}>
                        <td className="font-mono text-xs">{new Date(event.timestamp).toLocaleTimeString()}</td>
                        <td>
                          <span className={`badge badge--${event.severity === 'critical' ? 'error' : event.severity === 'warning' ? 'warning' : event.severity === 'info' ? 'info' : 'neutral'}`}>
                            {event.severity}
                          </span>
                        </td>
                        <td className="font-mono text-sm truncate max-w-[150px]">{event.source}</td>
                        <td className="truncate max-w-[300px]">{event.message}</td>
                        <td>
                          <div className="flex flex-wrap gap-1">
                            {event.tags.slice(0, 3).map(t => (
                              <span key={t} className="badge badge--neutral text-xs">{t}</span>
                            ))}
                            {event.tags.length > 3 && <span className="badge badge--neutral text-xs">+{event.tags.length - 3}</span>}
                          </div>
                        </td>
                        <td>
                          <button className="btn btn--ghost btn--sm" title="View details">
                            <EyeIcon className="w-4 h-4" />
                          </button>
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>

            {filteredEvents.length > 0 && (
              <div className="px-4 py-3 border-t border-border bg-secondary flex items-center justify-between">
                <span className="text-sm text-muted">Showing {(page - 1) * 50 + 1} to {Math.min(page * 50, filteredEvents.length)} of {filteredEvents.length} events</span>
                <div className="flex items-center gap-1">
                  <button className="btn btn--ghost btn--sm" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>
                    <ChevronLeftIcon className="w-4 h-4" />
                  </button>
                  <button className="btn btn--ghost btn--sm" onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}>
                    <ChevronRightIcon className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

function SearchIcon({ className }: { className?: string }) {
  return <svg className={className} width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
}
function FilterIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/><path d="M8 6h8M8 12h8M8 18h8"/></svg>
}
function DownloadIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
}
function EyeIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M1 12s4-8 11-8 11 8 11-8"/><circle cx="12" cy="12" r="3"/></svg>
}
function ChevronLeftIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
}
function ChevronRightIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
}
function ChevronDownIcon({ className }: { className?: string }) {
  return <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
}

export { ExplorerPage }