import React, { useState, useCallback, useRef, useEffect } from 'react'
import { format, subMinutes, subHours, subDays, subWeeks, subMonths, startOfDay, endOfDay, startOfWeek, endOfWeek, startOfMonth, endOfMonth } from 'date-fns'

interface TimeRangePickerProps {
  value: [Date, Date]
  onChange: (range: [Date, Date]) => void
  presets?: TimeRangePreset[]
  showCustom?: boolean
  showTime?: boolean
  timezone?: string
  minDate?: Date
  maxDate?: Date
  className?: string
}

interface TimeRangePreset {
  label: string
  value: [Date, Date] | (() => [Date, Date])
  group?: string
}

const DEFAULT_PRESETS: TimeRangePreset[] = [
  { label: 'Last 5 minutes', value: () => [subMinutes(new Date(), 5), new Date()], group: 'Recent' },
  { label: 'Last 15 minutes', value: () => [subMinutes(new Date(), 15), new Date()], group: 'Recent' },
  { label: 'Last 30 minutes', value: () => [subMinutes(new Date(), 30), new Date()], group: 'Recent' },
  { label: 'Last 1 hour', value: () => [subHours(new Date(), 1), new Date()], group: 'Recent' },
  { label: 'Last 3 hours', value: () => [subHours(new Date(), 3), new Date()], group: 'Recent' },
  { label: 'Last 6 hours', value: () => [subHours(new Date(), 6), new Date()], group: 'Recent' },
  { label: 'Last 12 hours', value: () => [subHours(new Date(), 12), new Date()], group: 'Recent' },
  { label: 'Last 24 hours', value: () => [subDays(new Date(), 1), new Date()], group: 'Recent' },
  { label: 'Last 2 days', value: () => [subDays(new Date(), 2), new Date()], group: 'Recent' },
  { label: 'Last 7 days', value: () => [subDays(new Date(), 7), new Date()], group: 'Recent' },
  { label: 'Last 30 days', value: () => [subDays(new Date(), 30), new Date()], group: 'Recent' },
  { label: 'Today', value: () => [startOfDay(new Date()), endOfDay(new Date())], group: 'Calendar' },
  { label: 'Yesterday', value: () => [startOfDay(subDays(new Date(), 1)), endOfDay(subDays(new Date(), 1))], group: 'Calendar' },
  { label: 'This week', value: () => [startOfWeek(new Date()), endOfWeek(new Date())], group: 'Calendar' },
  { label: 'Last week', value: () => [startOfWeek(subWeeks(new Date(), 1)), endOfWeek(subWeeks(new Date(), 1))], group: 'Calendar' },
  { label: 'This month', value: () => [startOfMonth(new Date()), endOfMonth(new Date())], group: 'Calendar' },
  { label: 'Last month', value: () => [startOfMonth(subMonths(new Date(), 1)), endOfMonth(subMonths(new Date(), 1))], group: 'Calendar' },
]

export function TimeRangePicker({
  value,
  onChange,
  presets = DEFAULT_PRESETS,
  showCustom = true,
  showTime = true,
  timezone = 'UTC',
  minDate,
  maxDate,
  className = '',
}: TimeRangePickerProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [selectedPreset, setSelectedPreset] = useState<string | null>(null)
  const [customStart, setCustomStart] = useState(value[0])
  const [customEnd, setCustomEnd] = useState(value[1])
  const [activeTab, setActiveTab] = useState<'presets' | 'custom'>('presets')
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handlePresetClick = useCallback((preset: TimeRangePreset) => {
    const range = typeof preset.value === 'function' ? preset.value() : preset.value
    onChange(range)
    setSelectedPreset(preset.label)
    setIsOpen(false)
  }, [onChange])

  const handleCustomApply = useCallback(() => {
    onChange([customStart, customEnd])
    setSelectedPreset(null)
    setIsOpen(false)
  }, [customStart, customEnd, onChange])

  const handleClear = useCallback(() => {
    const now = new Date()
    onChange([subHours(now, 1), now])
    setSelectedPreset('Last 1 hour')
  }, [onChange])

  const formatRange = (range: [Date, Date]) => {
    const start = range[0]
    const end = range[1]
    const diff = end.getTime() - start.getTime()
    
    if (diff < 3600000) {
      return `${format(start, 'HH:mm')} - ${format(end, 'HH:mm')}`
    } else if (diff < 86400000) {
      return `${format(start, 'MMM d, HH:mm')} - ${format(end, 'HH:mm')}`
    } else {
      return `${format(start, 'MMM d')} - ${format(end, 'MMM d')}`
    }
  }

  const groupedPresets = presets.reduce((acc, preset) => {
    const group = preset.group || 'Other'
    if (!acc[group]) acc[group] = []
    acc[group].push(preset)
    return acc
  }, {} as Record<string, TimeRangePreset[]>)

  return (
    <div ref={containerRef} className={`relative ${className}`}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-2 bg-secondary border border-border rounded-lg text-sm hover:bg-tertiary transition-colors"
        aria-haspopup="true"
        aria-expanded={isOpen}
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-muted">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
          <line x1="16" y1="2" x2="16" y2="6" />
          <line x1="8" y1="2" x2="8" y2="6" />
          <line x1="3" y1="10" x2="21" y2="10" />
        </svg>
        <span className="font-medium">
          {selectedPreset || formatRange(value)}
        </span>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={`transition-transform ${isOpen ? 'rotate-180' : ''}`}>
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-2 w-80 bg-secondary border border-border rounded-lg shadow-lg z-50 animate-slide-in">
          <div className="flex border-b border-border">
            <button
              onClick={() => setActiveTab('presets')}
              className={`flex-1 px-4 py-2 text-sm font-medium transition-colors ${activeTab === 'presets' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-primary'}`}
            >
              Presets
            </button>
            {showCustom && (
              <button
                onClick={() => setActiveTab('custom')}
                className={`flex-1 px-4 py-2 text-sm font-medium transition-colors ${activeTab === 'custom' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-primary'}`}
              >
                Custom
              </button>
            )}
          </div>

          <div className="max-h-96 overflow-y-auto">
            {activeTab === 'presets' && (
              <div className="p-2">
                {Object.entries(groupedPresets).map(([group, groupPresets]) => (
                  <div key={group} className="mb-3">
                    <p className="px-2 py-1 text-xs font-semibold text-muted uppercase tracking-wide">{group}</p>
                    <div className="space-y-1">
                      {groupPresets.map((preset, i) => (
                        <button
                          key={i}
                          onClick={() => handlePresetClick(preset)}
                          className={`w-full text-left px-3 py-2 rounded text-sm transition-colors ${selectedPreset === preset.label ? 'bg-primary-light text-primary font-medium' : 'hover:bg-tertiary text-secondary'}`}
                        >
                          {preset.label}
                        </button>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {activeTab === 'custom' && showCustom && (
              <div className="p-4 space-y-4">
                <div>
                  <label className="label">Start</label>
                  <input
                    type={showTime ? 'datetime-local' : 'date'}
                    value={format(customStart, showTime ? "yyyy-MM-dd'T'HH:mm" : 'yyyy-MM-dd')}
                    onChange={(e) => setCustomStart(new Date(e.target.value))}
                    className="input"
                    min={minDate?.toISOString().slice(0, 16)}
                    max={maxDate?.toISOString().slice(0, 16)}
                  />
                </div>
                <div>
                  <label className="label">End</label>
                  <input
                    type={showTime ? 'datetime-local' : 'date'}
                    value={format(customEnd, showTime ? "yyyy-MM-dd'T'HH:mm" : 'yyyy-MM-dd')}
                    onChange={(e) => setCustomEnd(new Date(e.target.value))}
                    className="input"
                    min={minDate?.toISOString().slice(0, 16)}
                    max={maxDate?.toISOString().slice(0, 16)}
                  />
                </div>
                <div className="flex gap-2">
                  <button className="btn btn--secondary flex-1" onClick={() => setIsOpen(false)}>
                    Cancel
                  </button>
                  <button className="btn btn--primary flex-1" onClick={handleCustomApply}>
                    Apply
                  </button>
                </div>
              </div>
            )}
          </div>

          <div className="p-2 border-t border-border flex items-center justify-between">
            <button
              onClick={handleClear}
              className="text-xs text-muted hover:text-primary transition-colors"
            >
              Reset to default
            </button>
            <span className="text-xs text-muted">
              {timezone}
            </span>
          </div>
        </div>
      )}
    </div>
  )
}

export { TimeRangePicker }