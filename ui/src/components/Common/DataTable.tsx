import React, { useState, useCallback, useMemo } from 'react'

interface DataTableProps<T> {
  data: T[]
  columns: ColumnDef<T>[]
  pageSize?: number
  sortable?: boolean
  filterable?: boolean
  selectable?: boolean
  onRowClick?: (row: T) => void
  onSelectionChange?: (selected: T[]) => void
  loading?: boolean
  emptyMessage?: string
  className?: string
}

interface ColumnDef<T> {
  key: string
  header: string
  accessor: (row: T) => any
  sortable?: boolean
  filterable?: boolean
  width?: string
  align?: 'left' | 'center' | 'right'
  render?: (value: any, row: T) => React.ReactNode
}

type SortDirection = 'asc' | 'desc' | null

interface SortState {
  column: string | null
  direction: SortDirection
}

export function DataTable<T extends { id?: string }>({
  data,
  columns,
  pageSize = 25,
  sortable = true,
  filterable = true,
  selectable = false,
  onRowClick,
  onSelectionChange,
  loading = false,
  emptyMessage = 'No data available',
  className = '',
}: DataTableProps<T>) {
  const [currentPage, setCurrentPage] = useState(1)
  const [sortState, setSortState] = useState<SortState>({ column: null, direction: null })
  const [filterText, setFilterText] = useState('')
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())

  const filteredData = useMemo(() => {
    if (!filterText) return data
    const lowerFilter = filterText.toLowerCase()
    return data.filter(row =>
      columns.some(col => {
        const value = col.accessor(row)
        return String(value).toLowerCase().includes(lowerFilter)
      })
    )
  }, [data, columns, filterText])

  const sortedData = useMemo(() => {
    if (!sortState.column || !sortState.direction) return filteredData
    const column = columns.find(c => c.key === sortState.column)
    if (!column) return filteredData

    return [...filteredData].sort((a, b) => {
      const aVal = column.accessor(a)
      const bVal = column.accessor(b)
      
      if (aVal === bVal) return 0
      if (aVal === null || aVal === undefined) return 1
      if (bVal === null || bVal === undefined) return -1
      
      const comparison = aVal < bVal ? -1 : aVal > bVal ? 1 : 0
      return sortState.direction === 'asc' ? comparison : -comparison
    })
  }, [filteredData, sortState, columns])

  const paginatedData = useMemo(() => {
    const start = (currentPage - 1) * pageSize
    return sortedData.slice(start, start + pageSize)
  }, [sortedData, currentPage, pageSize])

  const totalPages = Math.ceil(sortedData.length / pageSize)

  const handleSort = useCallback((columnKey: string) => {
    if (!sortable) return
    const column = columns.find(c => c.key === columnKey)
    if (!column?.sortable) return

    setSortState(prev => {
      if (prev.column !== columnKey) {
        return { column: columnKey, direction: 'asc' }
      }
      if (prev.direction === 'asc') {
        return { column: columnKey, direction: 'desc' }
      }
      if (prev.direction === 'desc') {
        return { column: null, direction: null }
      }
      return { column: columnKey, direction: 'asc' }
    })
  }, [sortable, columns])

  const handleSelectAll = useCallback(() => {
    if (selectedRows.size === paginatedData.length) {
      setSelectedRows(new Set())
    } else {
      const newSelected = new Set(selectedRows)
      paginatedData.forEach(row => {
        if (row.id) newSelected.add(row.id)
      })
      setSelectedRows(newSelected)
    }
  }, [selectedRows, paginatedData])

  const handleSelectRow = useCallback((row: T) => {
    if (!row.id) return
    const newSelected = new Set(selectedRows)
    if (newSelected.has(row.id)) {
      newSelected.delete(row.id)
    } else {
      newSelected.add(row.id)
    }
    setSelectedRows(newSelected)
    if (onSelectionChange) {
      const selected = data.filter(r => r.id && newSelected.has(r.id))
      onSelectionChange(selected)
    }
  }, [selectedRows, data, onSelectionChange])

  const handlePageChange = useCallback((page: number) => {
    setCurrentPage(Math.max(1, Math.min(page, totalPages)))
  }, [totalPages])

  const getSortIcon = (columnKey: string) => {
    if (sortState.column !== columnKey) return null
    if (sortState.direction === 'asc') return '↑'
    if (sortState.direction === 'desc') return '↓'
    return null
  }

  if (loading) {
    return (
      <div className="card">
        <div className="card__body p-4">
          <div className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="skeleton h-10 rounded" />
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (sortedData.length === 0) {
    return (
      <div className="card">
        <div className="card__body">
          <div className="empty-state">
            <svg className="empty-state__icon" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <rect x="3" y="3" width="18" height="18" rx="2" />
              <line x1="3" y1="9" x2="21" y2="9" />
              <line x1="9" y1="21" x2="9" y2="9" />
            </svg>
            <p className="empty-state__description">{emptyMessage}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`card ${className}`}>
      {filterable && (
        <div className="card__header flex items-center justify-between">
          <div className="relative flex-1 max-w-sm">
            <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              placeholder="Filter..."
              value={filterText}
              onChange={(e) => setFilterText(e.target.value)}
              className="input pl-9 text-sm"
            />
          </div>
          <div className="text-sm text-muted">
            {sortedData.length} of {data.length} rows
          </div>
        </div>
      )}

      <div className="overflow-x-auto">
        <table className="table">
          <thead>
            <tr>
              {selectable && (
                <th style={{ width: '40px' }}>
                  <input
                    type="checkbox"
                    checked={selectedRows.size === paginatedData.length && paginatedData.length > 0}
                    onChange={handleSelectAll}
                    className="w-4 h-4 rounded border-border"
                  />
                </th>
              )}
              {columns.map(col => (
                <th
                  key={col.key}
                  style={{ width: col.width }}
                  className={`${col.sortable !== false && sortable ? 'cursor-pointer select-none' : ''} ${col.align === 'center' ? 'text-center' : col.align === 'right' ? 'text-right' : ''}`}
                  onClick={() => handleSort(col.key)}
                >
                  <div className="flex items-center gap-1">
                    {col.header}
                    {getSortIcon(col.key) && (
                      <span className="text-primary">{getSortIcon(col.key)}</span>
                    )}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {paginatedData.map((row, rowIndex) => (
              <tr
                key={row.id || rowIndex}
                onClick={() => onRowClick?.(row)}
                className={`${onRowClick ? 'cursor-pointer' : ''} ${selectedRows.has(row.id || '') ? 'bg-primary-light' : ''}`}
              >
                {selectable && (
                  <td>
                    <input
                      type="checkbox"
                      checked={selectedRows.has(row.id || '')}
                      onChange={() => handleSelectRow(row)}
                      className="w-4 h-4 rounded border-border"
                      onClick={(e) => e.stopPropagation()}
                    />
                  </td>
                )}
                {columns.map(col => (
                  <td
                    key={col.key}
                    className={col.align === 'center' ? 'text-center' : col.align === 'right' ? 'text-right' : ''}
                  >
                    {col.render ? col.render(col.accessor(row), row) : col.accessor(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="card__footer flex items-center justify-between">
          <span className="text-sm text-muted">
            Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, sortedData.length)} of {sortedData.length} rows
          </span>
          <div className="flex items-center gap-1">
            <button
              onClick={() => handlePageChange(1)}
              disabled={currentPage === 1}
              className="btn btn--ghost btn--sm"
            >
              ««
            </button>
            <button
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="btn btn--ghost btn--sm"
            >
              «
            </button>
            {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              const page = Math.max(1, Math.min(currentPage - 2 + i, totalPages - 4 + i))
              return (
                <button
                  key={page}
                  onClick={() => handlePageChange(page)}
                  className={`btn btn--sm ${currentPage === page ? 'btn--primary' : 'btn--ghost'}`}
                >
                  {page}
                </button>
              )
            })}
            <button
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="btn btn--ghost btn--sm"
            >
              »
            </button>
            <button
              onClick={() => handlePageChange(totalPages)}
              disabled={currentPage === totalPages}
              className="btn btn--ghost btn--sm"
            >
              »»
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

export { DataTable }