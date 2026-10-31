import React, { createContext, useContext, useState, useCallback, useEffect } from 'react'

export interface QueryRequest {
  expression: string
  tenant?: string
  timeRange?: [number, number]
  step?: number
}

export interface QueryResult {
  resultType: 'vector' | 'matrix' | 'scalar' | 'string'
  result: MetricSample[]
  stats: QueryStats
  error?: string
}

export interface MetricSample {
  labels: Record<string, string>
  value: number
  timestamp: number
}

export interface QueryStats {
  rowsRead: number
  rowsFiltered: number
  executionTimeMs: number
  cacheHits: number
  cacheMisses: number
}

export interface QueryHistoryItem {
  id: string
  expression: string
  timestamp: number
  tenant?: string
  resultType: string
  rowCount: number
  executionTimeMs: number
  error?: string
}

interface QueryContextType {
  // Current query state
  currentQuery: string
  setCurrentQuery: (query: string) => void
  
  // Query execution
  executeQuery: (request: QueryRequest) => Promise<QueryResult>
  isExecuting: boolean
  currentResult: QueryResult | null
  error: string | null
  
  // Query history
  history: QueryHistoryItem[]
  addToHistory: (item: QueryHistoryItem) => void
  clearHistory: () => void
  
  // Saved queries
  savedQueries: SavedQuery[]
  saveQuery: (query: SavedQuery) => void
  deleteSavedQuery: (id: string) => void
  
  // Query editor state
  editorLanguage: 'promql' | 'helix'
  setEditorLanguage: (lang: 'promql' | 'helix') => void
}

export interface SavedQuery {
  id: string
  name: string
  expression: string
  description?: string
  tags: string[]
  createdAt: number
  updatedAt: number
  tenant?: string
}

const QueryContext = createContext<QueryContextType | undefined>(undefined)

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const [currentQuery, setCurrentQuery] = useState('')
  const [isExecuting, setIsExecuting] = useState(false)
  const [currentResult, setCurrentResult] = useState<QueryResult | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [history, setHistory] = useState<QueryHistoryItem[]>(() => {
    try {
      const stored = localStorage.getItem('helix-query-history')
      return stored ? JSON.parse(stored) : []
    } catch {
      return []
    }
  })
  const [savedQueries, setSavedQueries] = useState<SavedQuery[]>(() => {
    try {
      const stored = localStorage.getItem('helix-saved-queries')
      return stored ? JSON.parse(stored) : []
    } catch {
      return []
    }
  })
  const [editorLanguage, setEditorLanguage] = useState<'promql' | 'helix'>('promql')

  useEffect(() => {
    localStorage.setItem('helix-query-history', JSON.stringify(history.slice(0, 100)))
  }, [history])

  useEffect(() => {
    localStorage.setItem('helix-saved-queries', JSON.stringify(savedQueries))
  }, [savedQueries])

  const executeQuery = useCallback(async (request: QueryRequest): Promise<QueryResult> => {
    setIsExecuting(true)
    setError(null)
    
    try {
      const startTime = performance.now()
      
      // Build query URL
      const params = new URLSearchParams()
      params.append('query', request.expression)
      if (request.tenant) params.append('tenant', request.tenant)
      if (request.timeRange) {
        params.append('start', request.timeRange[0].toString())
        params.append('end', request.timeRange[1].toString())
      }
      if (request.step) params.append('step', request.step.toString())
      
      const response = await fetch(`/api/v1/query?${params.toString()}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('helix-auth-token') || ''}`,
        },
      })
      
      const endTime = performance.now()
      const executionTimeMs = endTime - startTime
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || `Query failed: ${response.status}`)
      }
      
      const data = await response.json()
      
      const result: QueryResult = {
        resultType: data.resultType || 'vector',
        result: data.result || [],
        stats: {
          rowsRead: data.stats?.rowsRead || 0,
          rowsFiltered: data.stats?.rowsFiltered || 0,
          executionTimeMs,
          cacheHits: data.stats?.cacheHits || 0,
          cacheMisses: data.stats?.cacheMisses || 0,
        },
      }
      
      setCurrentResult(result)
      
      // Add to history
      const historyItem: QueryHistoryItem = {
        id: crypto.randomUUID(),
        expression: request.expression,
        tenant: request.tenant,
        resultType: result.resultType,
        rowCount: result.result.length,
        executionTimeMs,
        timestamp: Date.now(),
      }
      
      setHistory(prev => [historyItem, ...prev].slice(0, 100))
      
      return result
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error'
      setError(errorMessage)
      
      // Add failed query to history
      const historyItem: QueryHistoryItem = {
        id: crypto.randomUUID(),
        expression: request.expression,
        tenant: request.tenant,
        resultType: 'error',
        rowCount: 0,
        executionTimeMs: 0,
        timestamp: Date.now(),
        error: errorMessage,
      }
      
      setHistory(prev => [historyItem, ...prev].slice(0, 100))
      
      throw err
    } finally {
      setIsExecuting(false)
    }
  }, [])

  const addToHistory = useCallback((item: QueryHistoryItem) => {
    setHistory(prev => [item, ...prev].slice(0, 100))
  }, [])

  const clearHistory = useCallback(() => {
    setHistory([])
    localStorage.removeItem('helix-query-history')
  }, [])

  const saveQuery = useCallback((query: SavedQuery) => {
    setSavedQueries(prev => {
      const existing = prev.find(q => q.id === query.id)
      if (existing) {
        return prev.map(q => q.id === query.id ? { ...query, updatedAt: Date.now() } : q)
      }
      return [{ ...query, id: query.id || crypto.randomUUID(), createdAt: Date.now(), updatedAt: Date.now() }, ...prev]
    })
  }, [])

  const deleteSavedQuery = useCallback((id: string) => {
    setSavedQueries(prev => prev.filter(q => q.id !== id))
  }, [])

  const setCurrentQuery = useCallback((query: string) => {
    setCurrentQuery(query)
  }, [])

  const setEditorLanguage = useCallback((lang: 'promql' | 'helix') => {
    setEditorLanguage(lang)
  }, [])

  return (
    <QueryContext.Provider value={{
      currentQuery,
      setCurrentQuery,
      executeQuery,
      isExecuting,
      currentResult,
      error,
      history,
      addToHistory,
      clearHistory,
      savedQueries,
      saveQuery,
      deleteSavedQuery,
      editorLanguage,
      setEditorLanguage,
    }}>
      {children}
    </QueryContext.Provider>
  )
}

export function useQuery() {
  const context = useContext(QueryContext)
  if (!context) {
    throw new Error('useQuery must be used within a QueryProvider')
  }
  return context
}