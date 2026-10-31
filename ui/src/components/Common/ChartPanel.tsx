import React, { useState, useCallback, useEffect, useRef } from 'react'
import { Line, Bar } from 'react-chartjs-2'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend, Filler, TimeScale } from 'chart.js'
import 'chartjs-adapter-date-fns'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend, Filler, TimeScale)

interface ChartPanelProps {
  title: string
  type: 'line' | 'bar' | 'area'
  data: {
    labels: string[]
    datasets: Array<{
      label: string
      data: number[]
      color?: string
      fill?: boolean
    }>
  }
  height?: number
  showLegend?: boolean
  showGrid?: boolean
  timeRange?: [Date, Date]
  onTimeRangeChange?: (range: [Date, Date]) => void
  loading?: boolean
  emptyMessage?: string
  actions?: React.ReactNode
}

const CHART_COLORS = [
  '#4f46e5', '#10b981', '#f59e0b', '#ef4444', '#3b82f6',
  '#8b5cf6', '#ec4899', '#14b8a6', '#f97316', '#6366f1',
]

export function ChartPanel({
  title,
  type,
  data,
  height = 300,
  showLegend = true,
  showGrid = true,
  timeRange,
  onTimeRangeChange,
  loading = false,
  emptyMessage = 'No data available',
  actions,
}: ChartPanelProps) {
  const [chartType, setChartType] = useState(type)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const chartRef = useRef<any>(null)

  const chartData = {
    labels: data.labels,
    datasets: data.datasets.map((ds, i) => ({
      label: ds.label,
      data: ds.data,
      borderColor: ds.color || CHART_COLORS[i % CHART_COLORS.length],
      backgroundColor: chartType === 'bar' 
        ? (ds.color || CHART_COLORS[i % CHART_COLORS.length]) + '80'
        : (ds.color || CHART_COLORS[i % CHART_COLORS.length]) + '20',
      borderWidth: 2,
      fill: chartType === 'area' || ds.fill,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 4,
      pointHoverBackgroundColor: ds.color || CHART_COLORS[i % CHART_COLORS.length],
    })),
  }

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    plugins: {
      legend: {
        display: showLegend,
        position: 'top' as const,
        labels: {
          usePointStyle: true,
          padding: 16,
          font: { size: 12 },
        },
      },
      tooltip: {
        enabled: true,
        mode: 'index' as const,
        intersect: false,
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        titleFont: { size: 12, weight: 'bold' },
        bodyFont: { size: 11 },
        padding: 12,
        cornerRadius: 8,
        displayColors: true,
        callbacks: {
          label: (context: any) => {
            const label = context.dataset.label || ''
            const value = context.parsed.y
            return `${label}: ${value.toLocaleString()}`
          },
        },
      },
    },
    scales: {
      x: {
        display: true,
        grid: { display: showGrid, color: 'rgba(0, 0, 0, 0.05)' },
        ticks: {
          font: { size: 11 },
          maxRotation: 0,
          autoSkip: true,
          maxTicksLimit: 8,
        },
      },
      y: {
        display: true,
        grid: { display: showGrid, color: 'rgba(0, 0, 0, 0.05)' },
        ticks: {
          font: { size: 11 },
          callback: (value: any) => {
            if (value >= 1000000) return (value / 1000000).toFixed(1) + 'M'
            if (value >= 1000) return (value / 1000).toFixed(1) + 'K'
            return value
          },
        },
      },
    },
    elements: {
      line: { tension: 0.4 },
      point: { radius: 0, hoverRadius: 4 },
    },
    animation: {
      duration: 750,
      easing: 'easeInOutQuart' as const,
    },
  }

  if (loading) {
    return (
      <div className="card">
        <div className="card__header flex items-center justify-between">
          <div className="skeleton skeleton--text w-32" />
          <div className="flex gap-2">
            <div className="skeleton w-8 h-8 rounded" />
            <div className="skeleton w-8 h-8 rounded" />
          </div>
        </div>
        <div className="card__body" style={{ height }}>
          <div className="skeleton h-full rounded" />
        </div>
      </div>
    )
  }

  if (!data.labels.length || !data.datasets.length) {
    return (
      <div className="card">
        <div className="card__header flex items-center justify-between">
          <h3 className="font-semibold">{title}</h3>
          {actions}
        </div>
        <div className="card__body flex items-center justify-center" style={{ height }}>
          <div className="empty-state">
            <svg className="empty-state__icon" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M3 3v18h18" />
              <path d="M18 17V9" />
              <path d="M13 17V5" />
              <path d="M8 17v-3" />
            </svg>
            <p className="empty-state__description">{emptyMessage}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`card ${isFullscreen ? 'fixed inset-4 z-50' : ''}`}>
      <div className="card__header flex items-center justify-between">
        <h3 className="font-semibold">{title}</h3>
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1 bg-tertiary rounded-lg p-1">
            <button
              onClick={() => setChartType('line')}
              className={`p-1.5 rounded ${chartType === 'line' ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="Line chart"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
              </svg>
            </button>
            <button
              onClick={() => setChartType('bar')}
              className={`p-1.5 rounded ${chartType === 'bar' ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="Bar chart"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="18" y1="20" x2="18" y2="10" />
                <line x1="12" y1="20" x2="12" y2="4" />
                <line x1="6" y1="20" x2="6" y2="14" />
              </svg>
            </button>
            <button
              onClick={() => setChartType('area')}
              className={`p-1.5 rounded ${chartType === 'area' ? 'bg-primary text-white' : 'text-muted hover:text-primary'}`}
              title="Area chart"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M3 3v18h18" />
                <path d="M18.7 8l-5.1 5.2-2.8-2.7L7 14.3" />
              </svg>
            </button>
          </div>
          <button
            onClick={() => setIsFullscreen(!isFullscreen)}
            className="p-1.5 rounded hover:bg-tertiary text-muted hover:text-primary transition-colors"
            title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
          >
            {isFullscreen ? (
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 1 2-2h3M3 16h3a2 2 0 0 1 2 2v3" />
              </svg>
            ) : (
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3" />
              </svg>
            )}
          </button>
          {actions}
        </div>
      </div>

      <div className="card__body" style={{ height: isFullscreen ? 'calc(100vh - 120px)' : height }}>
        {chartType === 'bar' ? (
          <Bar data={chartData} options={chartOptions} ref={chartRef} />
        ) : (
          <Line data={chartData} options={chartOptions} ref={chartRef} />
        )}
      </div>

      {timeRange && onTimeRangeChange && (
        <div className="card__footer flex items-center justify-between">
          <span className="text-xs text-muted">
            {timeRange[0].toLocaleString()} - {timeRange[1].toLocaleString()}
          </span>
          <div className="flex gap-2">
            <button className="btn btn--ghost btn--sm" onClick={() => {
              const now = new Date()
              onTimeRangeChange([new Date(now.getTime() - 3600000), now])
            }}>1h</button>
            <button className="btn btn--ghost btn--sm" onClick={() => {
              const now = new Date()
              onTimeRangeChange([new Date(now.getTime() - 21600000), now])
            }}>6h</button>
            <button className="btn btn--ghost btn--sm" onClick={() => {
              const now = new Date()
              onTimeRangeChange([new Date(now.getTime() - 86400000), now])
            }}>24h</button>
            <button className="btn btn--ghost btn--sm" onClick={() => {
              const now = new Date()
              onTimeRangeChange([new Date(now.getTime() - 604800000), now])
            }}>7d</button>
          </div>
        </div>
      )}
    </div>
  )
}

export { ChartPanel }