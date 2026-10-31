import React, { useState, useCallback, useEffect } from 'react'

interface FilterBarProps {
  filters: FilterState
  onChange: (filters: FilterState) => void
  onReset: () => void
  sources?: string[]
  tags?: string[]
  severities?: string[]
  showTimeRange?: boolean
  showSearch?: boolean
  className?: string
}

interface FilterState {
  query: string
  sources: string[]
  tags: string[]
  severity: string[]
  timeRange: number
}

const DEFAULT_SOURCES = ['soc-agent-01', 'soc-agent-02', 'web-server-01', 'db-server-01', 'api-gateway-01']
const DEFAULT_TAGS = ['production', 'staging', 'critical', 'warning', 'info', 'security', 'network', 'database']
const DEFAULT_SEVERITIES = ['critical', 'warning', 'info', 'debug']

const TIME_RANGES = [
  { label: '5m', value: 5 * 60 * 1000 },
  { label: '15m', value: 15 * 60 * 1000 },
  { label: '1h', value: 60 * 60 * 1000 },
  { label: '6h', value: 6 * 60 * 60 * 1000 },
  { label: '24h', value: 24 * 60 * 60 * 1000 },
  { label: '7d', value: 7 * 24 * 60 * 60 * 1000 },
]

export function FilterBar({
  filters,
  onChange,
  onReset,
  sources = DEFAULT_SOURCES,
  tags = DEFAULT_TAGS,
  severities = DEFAULT_SEVERITIES,
  showTimeRange = true,
  showSearch = true,
  className = '',
}: FilterBarProps) {
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [searchInput, setSearchInput] = useState(filters.query)

  useEffect(() => {
    const timer = setTimeout(() => {
      if (searchInput !== filters.query) {
        onChange({ ...filters, query: searchInput })
      }
    }, 300)
    return () => clearTimeout(timer)
  }, [searchInput])

  const handleSourceToggle = useCallback((source: string) => {
    const newSources = filters.sources.includes(source)
      ? filters.sources.filter(s => s !== source)
      : [...filters.sources, source]
    onChange({ ...filters, sources: newSources })
  }, [filters, onChange])

  const handleTagToggle = useCallback((tag: string) => {
    const newTags = filters.tags.includes(tag)
      ? filters.tags.filter(t => t !== tag)
      : [...filters.tags, tag]
    onChange({ ...filters, tags: newTags })
  }, [filters, onChange])

  const handleSeverityToggle = useCallback((severity: string) => {
    const newSeverity = filters.severity.includes(severity)
      ? filters.severity.filter(s => s !== severity)
      : [...filters.severity, severity]
    onChange({ ...filters, severity: newSeverity })
  }, [filters, onChange])

  const handleTimeRangeChange = useCallback((value: number) => {
    onChange({ ...filters, timeRange: value })
  }, [filters, onChange])

  const handleReset = useCallback(() => {
    setSearchInput('')
    onReset()
  }, [onReset])

  const activeFilterCount = filters.sources.length + filters.tags.length + filters.severity.length + (filters.query ? 1 : 0)

  return (
    <div className={`filter-bar ${className}`}>
      <div className="flex flex-wrap items-center gap-3 p-3 bg-secondary border-b border-border">
        {showSearch && (
          <div className="relative flex-1 min-w-[200px] max-w-md">
            <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              placeholder="Search..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="input pl-9 text-sm"
            />
            {searchInput && (
              <button
                onClick={() => setSearchInput('')}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-tertiary text-muted hover:text-primary"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            )}
          </div>
        )}

        {showTimeRange && (
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            {TIME_RANGES.map(tr => (
              <button
                key={tr.value}
                onClick={() => handleTimeRangeChange(tr.value)}
                className={`px-2 py-1 rounded text-xs font-medium transition-colors ${filters.timeRange === tr.value ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              >
                {tr.label}
              </button>
            ))}
          </div>
        )}

        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
          className={`btn btn--sm ${showAdvanced ? 'btn--primary' : 'btn--ghost'}`}
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
          </svg>
          Filters
          {activeFilterCount > 0 && (
            <span className="badge badge--info ml-1">{activeFilterCount}</span>
          )}
        </button>

        {activeFilterCount > 0 && (
          <button
            onClick={handleReset}
            className="btn btn--ghost btn--sm text-muted hover:text-error"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
            Clear
          </button>
        )}
      </div>

      {showAdvanced && (
        <div className="p-4 bg-secondary border-b border-border animate-slide-in">
          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <label className="label">Sources</label>
              <div className="flex flex-wrap gap-2 max-h-32 overflow-y-auto">
                {sources.map(source => (
                  <label
                    key={source}
                    className={`badge ${filters.sources.includes(source) ? 'badge--info' : 'badge--neutral'} cursor-pointer transition-colors`}
                  >
                    <input
                      type="checkbox"
                      className="sr-only"
                      checked={filters.sources.includes(source)}
                      onChange={() => handleSourceToggle(source)}
                    />
                    {source}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Tags</label>
              <div className="flex flex-wrap gap-2 max-h-32 overflow-y-auto">
                {tags.map(tag => (
                  <label
                    key={tag}
                    className={`badge ${filters.tags.includes(tag) ? 'badge--warning' : 'badge--neutral'} cursor-pointer transition-colors`}
                  >
                    <input
                      type="checkbox"
                      className="sr-only"
                      checked={filters.tags.includes(tag)}
                      onChange={() => handleTagToggle(tag)}
                    />
                    {tag}
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="label">Severity</label>
              <div className="flex flex-wrap gap-2">
                {severities.map(severity => (
                  <label
                    key={severity}
                    className={`badge ${filters.severity.includes(severity) ? `badge--${severity === 'critical' ? 'error' : severity === 'warning' ? 'warning' : 'info'}` : 'badge--neutral'} cursor-pointer transition-colors`}
                  >
                    <input
                      type="checkbox"
                      className="sr-only"
                      checked={filters.severity.includes(severity)}
                      onChange={() => handleSeverityToggle(severity)}
                    />
                    {severity}
                  </label>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export { FilterBar }