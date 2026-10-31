import React, { useState, useCallback, useEffect } from 'react'
import { Editor } from '@monaco-editor/react'
import { useQuery } from '../../context/QueryContext'
import { useToast } from '../../components/Common/ToastContainer'
import { ChevronDownIcon, PlayIcon, SaveIcon, DownloadIcon, TrashIcon, CopyIcon, HistoryIcon, ClockIcon, ZapIcon, BarChartIcon } from '../Layout/Sidebar'

const queryExamples = [
  { label: 'Rate of CPU usage', query: 'rate(cpu_usage[30s])' },
  { label: 'Memory usage by host', query: 'memory_usage{job="node"}' },
  { label: '95th percentile latency', query: 'histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))' },
  { label: 'Top 10 CPU consumers', query: 'topk(10, rate(cpu_usage[1m]))' },
  { label: 'Error rate by service', query: 'rate(http_requests_total{status=~"5.."}[5m])' },
  { label: 'Disk usage percentage', query: '(disk_used / disk_total) * 100' },
]

const TIME_RANGES = [
  { label: 'Last 5 minutes', value: 5 * 60 * 1000 },
  { label: 'Last 15 minutes', value: 15 * 60 * 1000 },
  { label: 'Last 1 hour', value: 60 * 60 * 1000 },
  { label: 'Last 6 hours', value: 6 * 60 * 60 * 1000 },
  { label: 'Last 24 hours', value: 24 * 60 * 60 * 1000 },
  { label: 'Last 7 days', value: 7 * 24 * 60 * 60 * 1000 },
]

const TENANTS = [
  { id: 'tenant-1', name: 'SOC Team Alpha' },
  { id: 'tenant-2', name: 'SOC Team Beta' },
  { id: 'tenant-3', name: 'Platform Team' },
]

interface QueryPageProps {}

export function QueryPage() {
  const { currentQuery, setCurrentQuery, executeQuery, isExecuting, currentResult, error, history, savedQueries, editorLanguage, setEditorLanguage } = useQuery()
  const { addToast } = useToast()
  const [query, setQuery] = useState('')
  const [showExamples, setShowExamples] = useState(false)
  const [showHistory, setShowHistory] = useState(false)
  const [showSavedQueries, setShowSavedQueries] = useState(false)
  const [selectedTenant, setSelectedTenant] = useState('tenant-1')
  const [selectedTimeRange, setSelectedTimeRange] = useState(TIME_RANGES[2].value)
  const [editorTheme, setEditorTheme] = useState('vs-dark')
  const [resultView, setResultView] = useState<'table' | 'chart' | 'raw'>('table')
  const editorRef = React.useRef<any>(null)

  // Sync with context
  useEffect(() => {
    if (currentQuery !== query) {
      setQuery(currentQuery)
    }
  }, [currentQuery])

  const handleQueryChange = useCallback((value: string) => {
    setQuery(value)
    setCurrentQuery(value)
  }, [setCurrentQuery])

  const handleExecute = useCallback(async () => {
    if (!query.trim()) return
    try {
      await executeQuery({ 
        expression: query.trim(), 
        tenant: selectedTenant,
        timeRange: [Date.now() - selectedTimeRange, Date.now()]
      })
      addToast({ type: 'success', title: 'Query executed successfully' })
    } catch (err) {
      // Error already handled in context
    }
  }, [query, selectedTenant, selectedTimeRange, executeQuery, addToast])

  const handleSave = useCallback(() => {
    if (!query.trim()) return
    const name = prompt('Enter a name for this query:')
    if (name) {
      // saveQuery logic would go here
      addToast({ type: 'success', title: 'Query saved', message: `"${name}" has been saved` })
    }
  }, [query, addToast])

  const handleExampleClick = useCallback((example: { label: string; query: string }) => {
    setQuery(example.query)
    setCurrentQuery(example.query)
    setShowExamples(false)
  }, [setCurrentQuery])

  const handleHistoryClick = useCallback((item: any) => {
    setQuery(item.expression)
    setCurrentQuery(item.expression)
    setShowHistory(false)
  }, [setCurrentQuery])

  const handleFormat = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.getAction('editor.action.formatDocument').run()
    }
  }, [])

  const editorOptions = {
    theme: editorTheme,
    language: editorLanguage === 'promql' ? 'promql' : 'helix',
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 13,
    lineNumbers: 'on',
    wordWrap: 'on',
    tabSize: 2,
    insertSpaces: true,
    bracketPairColorization: { enabled: true },
    guides: { bracketPairs: true },
    renderLineHighlight: 'line',
    scrollBeyondLastLine: false,
    smoothScrolling: true,
    cursorBlinking: 'smooth',
    cursorSmoothCaretAnimation: 'on',
  }

  return (
    <div className="h-full flex flex-col bg-primary">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-3 p-4 border-b border-border bg-secondary">
        <div className="flex items-center gap-2">
          <select
            value={selectedTenant}
            onChange={(e) => setSelectedTenant(e.target.value)}
            className="input input--sm w-auto"
            aria-label="Select tenant"
          >
            {TENANTS.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
          </select>
          
          <select
            value={selectedTimeRange}
            onChange={(e) => setSelectedTimeRange(Number(e.target.value))}
            className="input input--sm w-auto"
            aria-label="Select time range"
          >
            {TIME_RANGES.map(t => <option key={t.value} value={t.value}>{t.label}</option>)}
          </select>
        </div>

        <div className="flex-1 max-w-4xl mx-4">
          <div className="relative">
            <form onSubmit={(e) => { e.preventDefault(); handleExecute(); }} className="relative">
              <Editor
                ref={editorRef}
                height="120px"
                value={query}
                onChange={handleQueryChange}
                options={editorOptions}
                defaultLanguage={editorLanguage === 'promql' ? 'promql' : 'helix'}
              />
              <div className="absolute bottom-2 right-2 flex items-center gap-1 p-1 bg-secondary/90 backdrop-blur rounded-lg border border-border">
                <select
                  value={editorLanguage}
                  onChange={(e) => setEditorLanguage(e.target.value as 'promql' | 'helix')}
                  className="input input--sm text-xs py-1 px-2 bg-transparent border-0 focus:ring-0"
                  aria-label="Query language"
                >
                  <option value="promql">PromQL</option>
                  <option value="helix">HelixQL</option>
                </select>
                <button
                  onClick={() => setEditorTheme(editorTheme === 'vs-dark' ? 'vs' : 'vs-dark')}
                  className="btn btn--ghost btn--sm p-1.5"
                  title="Toggle theme"
                >
                  {editorTheme === 'vs-dark' ? (
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
                  ) : (
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>
                  )}
                </button>
                <button onClick={handleFormat} className="btn btn--ghost btn--sm p-1.5" title="Format (Shift+Alt+F)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="15" y1="4" x2="21" y2="4"/><line x1="9" y1="20" x2="21" y2="20"/><line x1="3" y1="8" x2="21" y2="8"/><line x1="9" y1="12" x2="3" y2="12"/><line x1="9" y1="16" x2="21" y2="16"/><line x1="15" y1="4" x2="9" y2="20"/></svg>
                </button>
              </div>
            </form>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowExamples(!showExamples)}
            className="btn btn--ghost btn--sm"
          >
            <ZapIcon className="w-4 h-4" />
            <span className="hidden sm:inline">Examples</span>
          </button>
          <button
            onClick={() => setShowHistory(!showHistory)}
            className="btn btn--ghost btn--sm"
          >
            <HistoryIcon className="w-4 h-4" />
            <span className="hidden sm:inline">History</span>
          </button>
          <button
            onClick={() => setShowSavedQueries(!showSavedQueries)}
            className="btn btn--ghost btn--sm"
          >
            <SaveIcon className="w-4 h-4" />
            <span className="hidden sm:inline">Saved</span>
          </button>
          <button
            onClick={handleExecute}
            disabled={isExecuting || !query.trim()}
            className="btn btn--primary"
          >
            {isExecuting ? (
              <>
                <svg className="loading-spinner loading-spinner--sm mr-2" />
                Executing...
              </>
            ) : (
              <>
                <PlayIcon className="w-4 h-4 mr-2" />
                Run
              </>
            )}
          </button>
        </div>
      </div>

      {/* Dropdown Panels */}
      {showExamples && (
        <div className="absolute top-full left-0 right-0 bg-secondary border-t border-border z-20 p-4 animate-slide-in">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold">Query Examples</h3>
              <button onClick={() => setShowExamples(false)} className="text-muted hover:text-primary">Close</button>
            </div>
            <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
              {queryExamples.map((ex, i) => (
                <button
                  key={i}
                  onClick={() => handleExampleClick(ex)}
                  className="p-3 bg-secondary border border-border rounded-lg hover:bg-tertiary transition-colors text-left group"
                >
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-medium text-sm">{ex.label}</span>
                    <span className="text-xs text-muted">Click to use</span>
                  </div>
                  <code className="text-xs font-mono text-primary bg-code px-2 py-1 rounded block truncate">{ex.query}</code>
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {showHistory && (
        <div className="absolute top-full left-0 right-0 bg-secondary border-t border-border z-20 p-4 animate-slide-in max-h-96 overflow-y-auto">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold">Query History</h3>
              <button onClick={() => setShowHistory(false)} className="text-muted hover:text-primary">Close</button>
            </div>
            {history.length === 0 ? (
              <div className="text-center py-8 text-muted">No queries yet</div>
            ) : (
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {history.slice(0, 20).map((item, i) => (
                  <button
                    key={i}
                    onClick={() => handleHistoryClick(item)}
                    className="w-full p-3 bg-secondary border border-border rounded-lg hover:bg-tertiary transition-colors text-left"
                  >
                    <div className="flex items-start justify-between gap-2 mb-1">
                      <span className="font-medium text-sm truncate">{item.expression}</span>
                      <span className={`text-xs px-2 py-0.5 rounded ${item.error ? 'bg-error-light text-error' : 'bg-success-light text-success'}`}>
                        {item.error ? 'Error' : item.resultType}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 text-xs text-muted">
                      <span className="font-mono">{item.rowCount} rows</span>
                      <span>•</span>
                      <span>{item.executionTimeMs.toFixed(1)}ms</span>
                      <span>•</span>
                      <span>{new Date(item.timestamp).toLocaleTimeString()}</span>
                      {item.tenant && <span>• {item.tenant}</span>}
                    </div>
                    {item.error && <p className="text-xs text-error mt-1">{item.error}</p>}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {showSavedQueries && (
        <div className="absolute top-full left-0 right-0 bg-secondary border-t border-border z-20 p-4 animate-slide-in">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold">Saved Queries</h3>
              <button onClick={() => setShowSavedQueries(false)} className="text-muted hover:text-primary">Close</button>
            </div>
            {savedQueries.length === 0 ? (
              <div className="text-center py-8 text-muted">No saved queries</div>
            ) : (
              <div className="space-y-2">
                {savedQueries.map((sq, i) => (
                  <button
                    key={i}
                    onClick={() => { setQuery(sq.expression); setCurrentQuery(sq.expression); setShowSavedQueries(false); }}
                    className="w-full p-3 bg-secondary border border-border rounded-lg hover:bg-tertiary transition-colors text-left"
                  >
                    <div className="flex items-start justify-between gap-2 mb-1">
                      <span className="font-medium text-sm truncate">{sq.name}</span>
                      <span className="text-xs text-muted">{new Date(sq.updatedAt).toLocaleDateString()}</span>
                    </div>
                    <code className="text-xs font-mono text-primary bg-code px-2 py-1 rounded block truncate">{sq.expression}</code>
                    {sq.tags.length && (
                      <div className="flex gap-1 mt-1 flex-wrap">
                        {sq.tags.map((tag: string) => (
                          <span key={tag} className="badge badge--neutral text-xs">{tag}</span>
                        ))}
                      </div>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export { QueryPage }