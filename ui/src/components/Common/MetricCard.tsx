import React, { useState, useCallback, useEffect, useRef } from 'react'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend, Filler } from 'chart.js'
import { Line, Bar } from 'react-chartjs-2'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend, Filler)

interface MetricCardProps {
  title: string
  value: number | string
  unit?: string
  trend?: 'up' | 'down' | 'stable'
  trendValue?: number
  icon?: React.ReactNode
  color?: 'primary' | 'success' | 'warning' | 'error' | 'info'
  onClick?: () => void
  loading?: boolean
  sparklineData?: number[]
}

const colorClasses = {
  primary: 'text-primary bg-primary-light',
  success: 'text-success bg-success-light',
  warning: 'text-warning bg-warning-light',
  error: 'text-error bg-error-light',
  info: 'text-info bg-info-light',
}

export function MetricCard({
  title,
  value,
  unit = '',
  trend,
  trendValue,
  icon,
  color = 'primary',
  onClick,
  loading = false,
  sparklineData,
}: MetricCardProps) {
  const [isHovered, setIsHovered] = useState(false)

  const trendIcons = {
    up: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" />
        <polyline points="17 6 23 6 23 12" />
      </svg>
    ),
    down: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="23 18 13.5 8.5 8.5 13.5 1 6" />
        <polyline points="17 18 23 18 23 12" />
      </svg>
    ),
    stable: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <line x1="5" y1="12" x2="19" y2="12" />
      </svg>
    ),
  }

  const trendColors = {
    up: 'text-success',
    down: 'text-error',
    stable: 'text-muted',
  }

  if (loading) {
    return (
      <div className="card p-4">
        <div className="flex items-center justify-between mb-3">
          <div className="skeleton skeleton--text w-24" />
          <div className="skeleton w-8 h-8 rounded" />
        </div>
        <div className="skeleton skeleton--title w-32 mb-2" />
        <div className="skeleton skeleton--text w-16" />
      </div>
    )
  }

  return (
    <div
      className={`card p-4 transition-all duration-200 ${onClick ? 'cursor-pointer hover:shadow-md' : ''} ${isHovered ? 'shadow-md' : ''}`}
      onClick={onClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1 min-w-0">
          <p className="text-xs font-medium text-muted uppercase tracking-wide">{title}</p>
        </div>
        {icon && (
          <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${colorClasses[color]}`}>
            {icon}
          </div>
        )}
      </div>

      <div className="flex items-end justify-between">
        <div>
          <p className="text-2xl font-bold text-primary">
            {typeof value === 'number' ? value.toLocaleString() : value}
            {unit && <span className="text-sm font-normal text-muted ml-1">{unit}</span>}
          </p>
          {trend && (
            <div className={`flex items-center gap-1 mt-1 text-sm ${trendColors[trend]}`}>
              {trendIcons[trend]}
              {trendValue !== undefined && (
                <span className="font-medium">
                  {trendValue > 0 ? '+' : ''}{trendValue}%
                </span>
              )}
            </div>
          )}
        </div>

        {sparklineData && sparklineData.length > 0 && (
          <div className="w-24 h-12">
            <Line
              data={{
                labels: sparklineData.map((_, i) => i.toString()),
                datasets: [{
                  data: sparklineData,
                  borderColor: `var(--color-${color === 'primary' ? 'primary' : color})`,
                  backgroundColor: 'transparent',
                  borderWidth: 2,
                  pointRadius: 0,
                  tension: 0.4,
                  fill: false,
                }],
              }}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false }, tooltip: { enabled: false } },
                scales: { x: { display: false }, y: { display: false } },
                elements: { line: { tension: 0.4 } },
              }}
            />
          </div>
        )}
      </div>
    </div>
  )
}

export { MetricCard }