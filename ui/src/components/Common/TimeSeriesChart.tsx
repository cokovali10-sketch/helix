import React, { useState, useCallback, useEffect, useRef, useMemo } from 'react'
import { Line } from 'react-chartjs-2'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler, TimeScale } from 'chart.js'
import 'chartjs-adapter-date-fns'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler, TimeScale)

interface TimeSeriesChartProps {
  title: string
  data: Array<{
    timestamp: number
    value: number
    labels?: Record<string, string>
  }>
  color?: string
  height?: number
  showPoints?: boolean
  showArea?: boolean
  showGrid?: boolean
  showLegend?: boolean
  timeFormat?: 'auto' | 'seconds' | 'minutes' | 'hours' | 'days'
  valueFormat?: 'number' | 'bytes' | 'duration' | 'percent'
  threshold?: number
  thresholdColor?: string
  onPointClick?: (point: { timestamp: number; value: number; labels?: Record<string, string> }) => void
  loading?: boolean
  emptyMessage?: string
}

const TIME_FORMATS = {
  auto: (ts: number) => {
    const date = new Date(ts)
    return date.toLocaleTimeString()
  },
  seconds: (ts: number) => {
    const date = new Date(ts)
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
  },
  minutes: (ts: number) => {
    const date = new Date(ts)
    return `${date.getHours()}:${date.getMinutes()}`
  },
  hours: (ts: number) => {
    const date = new Date(ts)
    return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:00`
  },
  days: (ts: number) => {
    const date = new Date(ts)
    return `${date.getMonth() + 1}/${date.getDate()}`
  },
}

const VALUE_FORMATS = {
  number: (v: number) => v.toLocaleString(),
  bytes: (v: number) => {
    if (v >= 1073741824) return (v / 1073741824).toFixed(2) + ' GB'
    if (v >= 1048576) return (v / 1048576).toFixed(2) + ' MB'
    if (v >= 1024) return (v / 1024).toFixed(2) + ' KB'
    return v.toFixed(0) + ' B'
  },
  duration: (v: number) => {
    if (v >= 3600000) return (v / 3600000).toFixed(2) + ' h'
    if (v >= 60000) return (v / 60000).toFixed(2) + ' m'
    if (v >= 1000) return (v / 1000).toFixed(2) + ' s'
    return v.toFixed(0) + ' ms'
  },
  percent: (v: number) => v.toFixed(2) + '%',
}

export function TimeSeriesChart({
  title,
  data,
  color = '#4f46e5',
  height = 300,
  showPoints = false,
  showArea = false,
  showGrid = true,
  showLegend = false,
  timeFormat = 'auto',
  valueFormat = 'number',
  threshold,
  thresholdColor = '#ef4444',
  onPointClick,
  loading = false,
  emptyMessage = 'No data available',
}: TimeSeriesChartProps) {
  const [hoveredPoint, setHoveredPoint] = useState<{ timestamp: number; value: number; labels?: Record<string, string> } | null>(null)
  const [zoomRange, setZoomRange] = useState<[number, number] | null>(null)

  const chartData = useMemo(() => {
    const sortedData = [...data].sort((a, b) => a.timestamp - b.timestamp)
    
    return {
      labels: sortedData.map(d => d.timestamp),
      datasets: [{
        label: title,
        data: sortedData.map(d => ({ x: d.timestamp, y: d.value })),
        borderColor: color,
        backgroundColor: showArea ? color + '20' : 'transparent',
        borderWidth: 2,
        pointRadius: showPoints ? 3 : 0,
        pointHoverRadius: 5,
        pointBackgroundColor: color,
        pointBorderColor: '#fff',
        pointBorderWidth: 2,
        tension: 0.4,
        fill: showArea,
      }],
    }
  }, [data, title, color, showPoints, showArea])

  const chartOptions = useMemo(() => ({
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    onClick: (event: any, elements: any[]) => {
      if (elements.length > 0 && onPointClick) {
        const index = elements[0].index
        const point = data[index]
        onPointClick(point)
      }
    },
    onHover: (event: any, elements: any[]) => {
      if (elements.length > 0) {
        const index = elements[0].index
        const point = data[index]
        setHoveredPoint(point)
      } else {
        setHoveredPoint(null)
      }
    },
    plugins: {
      legend: {
        display: showLegend,
        position: 'top' as const,
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
        callbacks: {
          title: (items: any[]) => {
            if (items.length > 0) {
              const ts = items[0].parsed.x
              return TIME_FORMATS[timeFormat](ts)
            }
            return ''
          },
          label: (context: any) => {
            const value = context.parsed.y
            return `${context.dataset.label}: ${VALUE_FORMATS[valueFormat](value)}`
          },
        },
      },
    },
    scales: {
      x: {
        type: 'time' as const,
        time: {
          unit: timeFormat === 'days' ? 'day' : timeFormat === 'hours' ? 'hour' : 'minute',
          displayFormats: {
            minute: 'HH:mm',
            hour: 'MMM d HH:mm',
            day: 'MMM d',
          },
        },
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
          callback: (value: any) => VALUE_FORMATS[valueFormat](value),
        },
      },
    },
    elements: {
      line: { tension: 0.4 },
      point: { radius: showPoints ? 3 : 0, hoverRadius: 5 },
    },
    animation: {
      duration: 750,
      easing: 'easeInOutQuart' as const,
    },
  }), [timeFormat, valueFormat, showGrid, showLegend, showPoints, onPointClick, data])

  if (loading) {
    return (
      <div className="card">
        <div className="card__header">
          <div className="skeleton skeleton--text w-32" />
        </div>
        <div className="card__body" style={{ height }}>
          <div className="skeleton h-full rounded" />
        </div>
      </div>
    )
  }

  if (!data.length) {
    return (
      <div className="card">
        <div className="card__header">
          <h3 className="font-semibold">{title}</h3>
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
    <div className="card">
      <div className="card__header flex items-center justify-between">
        <h3 className="font-semibold">{title}</h3>
        {hoveredPoint && (
          <div className="text-sm text-muted">
            {TIME_FORMATS[timeFormat](hoveredPoint.timestamp)}: {VALUE_FORMATS[valueFormat](hoveredPoint.value)}
          </div>
        )}
      </div>
      <div className="card__body" style={{ height }}>
        <Line data={chartData} options={chartOptions} />
        {threshold !== undefined && (
          <div
            className="absolute left-0 right-0 border-t-2 border-dashed"
            style={{
              bottom: `${(threshold / (Math.max(...data.map(d => d.value)) || 1)) * 100}%`,
              borderColor: thresholdColor,
            }}
          />
        )}
      </div>
    </div>
  )
}

export { TimeSeriesChart }