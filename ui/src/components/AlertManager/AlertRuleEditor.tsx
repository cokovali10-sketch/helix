import React, { useState, useCallback, useEffect } from 'react'
import { useToast } from '../../components/Common/ToastContainer'

interface AlertRuleEditorProps {
  initialRule?: AlertRule
  onSave: (rule: AlertRule) => void
  onCancel: () => void
  availableMetrics?: string[]
}

interface AlertRule {
  id: string
  name: string
  expression: string
  duration: string
  severity: 'critical' | 'warning' | 'info'
  labels: Record<string, string>
  annotations: Record<string, string>
  enabled: boolean
  evaluationInterval: string
  forDuration: string
}

const SEVERITY_OPTIONS = [
  { value: 'info', label: 'Info', color: 'text-info' },
  { value: 'warning', label: 'Warning', color: 'text-warning' },
  { value: 'critical', label: 'Critical', color: 'text-error' },
]

const DURATION_OPTIONS = [
  { value: '1m', label: '1 minute' },
  { value: '5m', label: '5 minutes' },
  { value: '10m', label: '10 minutes' },
  { value: '15m', label: '15 minutes' },
  { value: '30m', label: '30 minutes' },
  { value: '1h', label: '1 hour' },
  { value: '2h', label: '2 hours' },
  { value: '6h', label: '6 hours' },
  { value: '12h', label: '12 hours' },
  { value: '24h', label: '24 hours' },
]

const EVALUATION_INTERVALS = [
  { value: '10s', label: '10 seconds' },
  { value: '15s', label: '15 seconds' },
  { value: '30s', label: '30 seconds' },
  { value: '1m', label: '1 minute' },
  { value: '5m', label: '5 minutes' },
]

const COMMON_EXPRESSIONS = [
  { label: 'High CPU', expression: 'cpu_usage > 90', description: 'CPU usage above 90%' },
  { label: 'High Memory', expression: 'memory_usage > 90', description: 'Memory usage above 90%' },
  { label: 'Disk Full', expression: 'disk_usage > 95', description: 'Disk usage above 95%' },
  { label: 'Service Down', expression: 'up == 0', description: 'Service is down' },
  { label: 'High Error Rate', expression: 'error_rate > 5', description: 'Error rate above 5%' },
  { label: 'High Latency', expression: 'latency_p99 > 1000', description: 'P99 latency above 1000ms' },
  { label: 'Low Availability', expression: 'availability < 99.9', description: 'Availability below 99.9%' },
  { label: 'High Request Rate', expression: 'rate(http_requests_total[5m]) > 1000', description: 'Request rate above 1000/s' },
]

export function AlertRuleEditor({ initialRule, onSave, onCancel, availableMetrics = [] }: AlertRuleEditorProps) {
  const { addToast } = useToast()
  const [rule, setRule] = useState<AlertRule>(initialRule || {
    id: crypto.randomUUID(),
    name: '',
    expression: '',
    duration: '5m',
    severity: 'warning',
    labels: {},
    annotations: {},
    enabled: true,
    evaluationInterval: '15s',
    forDuration: '5m',
  })
  const [showExpressionHelp, setShowExpressionHelp] = useState(false)
  const [showLabelEditor, setShowLabelEditor] = useState(false)
  const [showAnnotationEditor, setShowAnnotationEditor] = useState(false)
  const [newLabelKey, setNewLabelKey] = useState('')
  const [newLabelValue, setNewLabelValue] = useState('')
  const [newAnnotationKey, setNewAnnotationKey] = useState('')
  const [newAnnotationValue, setNewAnnotationValue] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = useCallback(() => {
    const newErrors: Record<string, string> = {}
    
    if (!rule.name.trim()) {
      newErrors.name = 'Name is required'
    }
    
    if (!rule.expression.trim()) {
      newErrors.expression = 'Expression is required'
    }
    
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }, [rule])

  const handleSave = useCallback(() => {
    if (!validate()) {
      addToast({ type: 'error', title: 'Validation failed', message: 'Please fix the errors' })
      return
    }
    
    onSave(rule)
    addToast({ type: 'success', title: 'Alert rule saved' })
  }, [rule, validate, onSave, addToast])

  const handleAddLabel = useCallback(() => {
    if (newLabelKey && newLabelValue) {
      setRule(prev => ({
        ...prev,
        labels: { ...prev.labels, [newLabelKey]: newLabelValue },
      }))
      setNewLabelKey('')
      setNewLabelValue('')
    }
  }, [newLabelKey, newLabelValue])

  const handleRemoveLabel = useCallback((key: string) => {
    setRule(prev => {
      const { [key]: _, ...rest } = prev.labels
      return { ...prev, labels: rest }
    })
  }, [])

  const handleAddAnnotation = useCallback(() => {
    if (newAnnotationKey && newAnnotationValue) {
      setRule(prev => ({
        ...prev,
        annotations: { ...prev.annotations, [newAnnotationKey]: newAnnotationValue },
      }))
      setNewAnnotationKey('')
      setNewAnnotationValue('')
    }
  }, [newAnnotationKey, newAnnotationValue])

  const handleRemoveAnnotation = useCallback((key: string) => {
    setRule(prev => {
      const { [key]: _, ...rest } = prev.annotations
      return { ...prev, annotations: rest }
    })
  }, [])

  const handleExpressionSelect = useCallback((expr: string) => {
    setRule(prev => ({ ...prev, expression: expr }))
    setShowExpressionHelp(false)
  }, [])

  return (
    <div className="space-y-6">
      {/* Basic Info */}
      <div className="card">
        <div className="card__header">
          <h3 className="font-semibold">Basic Information</h3>
        </div>
        <div className="card__body space-y-4">
          <div>
            <label className="label" htmlFor="name">Alert Name</label>
            <input
              id="name"
              type="text"
              className={`input ${errors.name ? 'border-error' : ''}`}
              value={rule.name}
              onChange={(e) => setRule({ ...rule, name: e.target.value })}
              placeholder="e.g., HighCPUUsage"
            />
            {errors.name && <p className="text-xs text-error mt-1">{errors.name}</p>}
          </div>

          <div>
            <label className="label" htmlFor="expression">Expression</label>
            <div className="relative">
              <input
                id="expression"
                type="text"
                className={`input font-mono ${errors.expression ? 'border-error' : ''}`}
                value={rule.expression}
                onChange={(e) => setRule({ ...rule, expression: e.target.value })}
                placeholder="e.g., cpu_usage > 90"
              />
              <button
                onClick={() => setShowExpressionHelp(!showExpressionHelp)}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-tertiary text-muted hover:text-primary"
                title="Show examples"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <circle cx="12" cy="12" r="10" />
                  <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
                  <line x1="12" y1="17" x2="12.01" y2="17" />
                </svg>
              </button>
            </div>
            {errors.expression && <p className="text-xs text-error mt-1">{errors.expression}</p>}
            
            {showExpressionHelp && (
              <div className="mt-2 p-3 bg-tertiary rounded-lg border border-border">
                <p className="text-xs font-semibold mb-2">Common Expressions</p>
                <div className="grid gap-2 sm:grid-cols-2">
                  {COMMON_EXPRESSIONS.map((ex, i) => (
                    <button
                      key={i}
                      onClick={() => handleExpressionSelect(ex.expression)}
                      className="p-2 bg-secondary rounded text-left hover:bg-primary-light transition-colors"
                    >
                      <p className="text-xs font-medium">{ex.label}</p>
                      <code className="text-xs font-mono text-primary">{ex.expression}</code>
                      <p className="text-xs text-muted mt-1">{ex.description}</p>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="label" htmlFor="severity">Severity</label>
              <select
                id="severity"
                className="input"
                value={rule.severity}
                onChange={(e) => setRule({ ...rule, severity: e.target.value as any })}
              >
                {SEVERITY_OPTIONS.map(opt => (
                  <option key={opt.value} value={opt.value}>{opt.label}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="label" htmlFor="duration">For Duration</label>
              <select
                id="duration"
                className="input"
                value={rule.forDuration}
                onChange={(e) => setRule({ ...rule, forDuration: e.target.value })}
              >
                {DURATION_OPTIONS.map(opt => (
                  <option key={opt.value} value={opt.value}>{opt.label}</option>
                ))}
              </select>
            </div>
          </div>

          <div>
            <label className="label" htmlFor="evalInterval">Evaluation Interval</label>
            <select
              id="evalInterval"
              className="input"
              value={rule.evaluationInterval}
              onChange={(e) => setRule({ ...rule, evaluationInterval: e.target.value })}
            >
              {EVALUATION_INTERVALS.map(opt => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enabled"
              className="w-4 h-4 rounded border-border text-primary focus:ring-primary"
              checked={rule.enabled}
              onChange={(e) => setRule({ ...rule, enabled: e.target.checked })}
            />
            <label className="label" htmlFor="enabled">Enabled</label>
          </div>
        </div>
      </div>

      {/* Labels */}
      <div className="card">
        <div className="card__header flex items-center justify-between">
          <h3 className="font-semibold">Labels</h3>
          <button
            onClick={() => setShowLabelEditor(!showLabelEditor)}
            className="btn btn--ghost btn--sm"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <line x1="12" y1="5" x2="12" y2="19" />
              <line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Add Label
          </button>
        </div>
        <div className="card__body">
          {Object.keys(rule.labels).length === 0 ? (
            <p className="text-sm text-muted">No labels defined</p>
          ) : (
            <div className="space-y-2">
              {Object.entries(rule.labels).map(([key, value]) => (
                <div key={key} className="flex items-center justify-between p-2 bg-tertiary rounded">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-sm text-primary">{key}</span>
                    <span className="text-muted">=</span>
                    <span className="font-mono text-sm">{value}</span>
                  </div>
                  <button
                    onClick={() => handleRemoveLabel(key)}
                    className="p-1 rounded hover:bg-error-light text-muted hover:text-error transition-colors"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18" />
                      <line x1="6" y1="6" x2="18" y2="18" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}

          {showLabelEditor && (
            <div className="mt-3 p-3 bg-tertiary rounded-lg border border-border">
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Key"
                  value={newLabelKey}
                  onChange={(e) => setNewLabelKey(e.target.value)}
                  className="input input--sm flex-1"
                />
                <input
                  type="text"
                  placeholder="Value"
                  value={newLabelValue}
                  onChange={(e) => setNewLabelValue(e.target.value)}
                  className="input input--sm flex-1"
                />
                <button className="btn btn--primary btn--sm" onClick={handleAddLabel}>Add</button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Annotations */}
      <div className="card">
        <div className="card__header flex items-center justify-between">
          <h3 className="font-semibold">Annotations</h3>
          <button
            onClick={() => setShowAnnotationEditor(!showAnnotationEditor)}
            className="btn btn--ghost btn--sm"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-1">
              <line x1="12" y1="5" x2="12" y2="19" />
              <line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Add Annotation
          </button>
        </div>
        <div className="card__body">
          {Object.keys(rule.annotations).length === 0 ? (
            <p className="text-sm text-muted">No annotations defined</p>
          ) : (
            <div className="space-y-2">
              {Object.entries(rule.annotations).map(([key, value]) => (
                <div key={key} className="flex items-center justify-between p-2 bg-tertiary rounded">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-sm text-primary">{key}</span>
                    <span className="text-muted">=</span>
                    <span className="text-sm">{value}</span>
                  </div>
                  <button
                    onClick={() => handleRemoveAnnotation(key)}
                    className="p-1 rounded hover:bg-error-light text-muted hover:text-error transition-colors"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18" />
                      <line x1="6" y1="6" x2="18" y2="18" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}

          {showAnnotationEditor && (
            <div className="mt-3 p-3 bg-tertiary rounded-lg border border-border">
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Key"
                  value={newAnnotationKey}
                  onChange={(e) => setNewAnnotationKey(e.target.value)}
                  className="input input--sm flex-1"
                />
                <input
                  type="text"
                  placeholder="Value"
                  value={newAnnotationValue}
                  onChange={(e) => setNewAnnotationValue(e.target.value)}
                  className="input input--sm flex-1"
                />
                <button className="btn btn--primary btn--sm" onClick={handleAddAnnotation}>Add</button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center justify-end gap-3">
        <button className="btn btn--secondary" onClick={onCancel}>Cancel</button>
        <button className="btn btn--primary" onClick={handleSave}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
            <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" />
            <polyline points="17 21 17 13 7 13 7 21" />
            <polyline points="7 3 7 8 15 8" />
          </svg>
          Save Alert Rule
        </button>
      </div>
    </div>
  )
}

export { AlertRuleEditor }