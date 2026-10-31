import React, { useState, useCallback, useRef, useEffect } from 'react'
import Editor from '@monaco-editor/react'
import { useQuery } from '../../context/QueryContext'
import { useToast } from '../../components/Common/ToastContainer'
import { PlayIcon, SaveIcon, DownloadIcon, TrashIcon, CopyIcon, HistoryIcon, ClockIcon, ZapIcon, BarChartIcon, ChevronDownIcon, ChevronRightIcon, FilterIcon, SearchIcon } from '../Layout/Sidebar'

const MONACO_OPTIONS = {
  automaticLayout: true,
  minimap: { enabled: false },
  fontSize: 13,
  lineNumbers: 'on' as const,
  wordWrap: 'on' as const,
  tabSize: 2,
  insertSpaces: true,
  bracketPairColorization: { enabled: true },
  guides: { bracketPairs: true },
  renderLineHighlight: 'line' as const,
  scrollBeyondLastLine: false,
  smoothScrolling: true,
  cursorBlinking: 'smooth' as const,
  cursorSmoothCaretAnimation: 'on' as const,
  suggest: {
    showKeywords: true,
    showSnippets: true,
    showFunctions: true,
  },
  quickSuggestions: {
    other: true,
    comments: false,
    strings: false,
  },
  parameterHints: { enabled: true },
  folding: true,
  foldingStrategy: 'indentation' as const,
  showFoldingControls: 'always' as const,
  renderWhitespace: 'selection' as const,
  rulers: [80, 120],
  colorDecorators: true,
  contextmenu: true,
  mouseWheelZoom: true,
  multiCursorModifier: 'ctrlCmd' as const,
  accessibilitySupport: 'auto' as const,
}

const QUERY_LANGUAGE_CONFIG = {
  promql: {
    keywords: [
      'by', 'without', 'on', 'ignoring', 'group_left', 'group_right',
      'bool', 'offset', 'rate', 'irate', 'increase', 'delta', 'idelta',
      'sum', 'avg', 'min', 'max', 'count', 'stddev', 'stdvar',
      'topk', 'bottomk', 'quantile', 'histogram_quantile',
      'abs', 'ceil', 'floor', 'round', 'sqrt', 'exp', 'ln', 'log2', 'log10',
      'clamp', 'clamp_min', 'clamp_max', 'deriv', 'predict_linear',
      'holt_winters', 'label_replace', 'label_join', 'vector', 'scalar',
      'sort', 'sort_desc', 'time', 'timestamp', 'day_of_month', 'day_of_week',
      'day_of_year', 'hour', 'minute', 'month', 'year', 'days_in_month',
      'changes', 'resets', 'idelta', 'deriv',
    ],
    operators: ['+', '-', '*', '/', '%', '^', '=', '!=', '<', '>', '<=', '>=', '=~', '!~'],
    brackets: ['(', ')', '{', '}', '[', ']'],
  },
  helix: {
    keywords: [
      'by', 'without', 'on', 'ignoring', 'group_left', 'group_right',
      'bool', 'offset', 'rate', 'irate', 'increase', 'delta', 'idelta',
      'sum', 'avg', 'min', 'max', 'count', 'stddev', 'stdvar',
      'topk', 'bottomk', 'quantile', 'histogram_quantile',
      'abs', 'ceil', 'floor', 'round', 'sqrt', 'exp', 'ln', 'log2', 'log10',
      'clamp', 'clamp_min', 'clamp_max', 'deriv', 'predict_linear',
      'holt_winters', 'label_replace', 'label_join', 'vector', 'scalar',
      'sort', 'sort_desc', 'time', 'timestamp', 'day_of_month', 'day_of_week',
      'day_of_year', 'hour', 'minute', 'month', 'year', 'days_in_month',
      'changes', 'resets', 'idelta', 'deriv',
      'filter', 'select', 'where', 'group', 'order', 'limit', 'offset',
    ],
    operators: ['+', '-', '*', '/', '%', '^', '=', '!=', '<', '>', '<=', '>=', '=~', '!~', '&&', '||'],
    brackets: ['(', ')', '{', '}', '[', ']'],
  },
}

interface QueryEditorProps {
  value: string
  onChange: (value: string) => void
  onExecute: () => void
  isExecuting: boolean
  language: 'promql' | 'helix'
  onLanguageChange: (lang: 'promql' | 'helix') => void
  theme: 'light' | 'dark'
  showLineNumbers?: boolean
  readOnly?: boolean
  height?: string
  placeholder?: string
}

export function QueryEditor({
  value,
  onChange,
  onExecute,
  isExecuting,
  language,
  onLanguageChange,
  theme,
  showLineNumbers = true,
  readOnly = false,
  height = '200px',
  placeholder = 'Enter your query here...',
}: QueryEditorProps) {
  const editorRef = useRef<any>(null)
  const [isFocused, setIsFocused] = useState(false)
  const [showSuggestions, setShowSuggestions] = useState(false)

  const handleEditorDidMount = useCallback((editor: any, monaco: any) => {
    editorRef.current = editor

    // Add keyboard shortcuts
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      onExecute()
    })

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
      // Save query
      console.log('Save query')
    })

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyF, () => {
      editor.getAction('editor.action.formatDocument').run()
    })

    // Focus editor
    editor.focus()
  }, [onExecute])

  const handleEditorChange = useCallback((value: string | undefined) => {
    if (value !== undefined) {
      onChange(value)
    }
  }, [onChange])

  const handleFormat = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.getAction('editor.action.formatDocument').run()
    }
  }, [])

  const handleClear = useCallback(() => {
    onChange('')
    if (editorRef.current) {
      editorRef.current.focus()
    }
  }, [onChange])

  const handleCopy = useCallback(() => {
    if (value) {
      navigator.clipboard.writeText(value)
    }
  }, [value])

  const handleRun = useCallback(() => {
    if (value.trim() && !isExecuting) {
      onExecute()
    }
  }, [value, isExecuting, onExecute])

  return (
    <div className={`query-editor ${isFocused ? 'query-editor--focused' : ''}`}>
      <div className="query-editor__toolbar">
        <div className="query-editor__toolbar-left">
          <select
            value={language}
            onChange={(e) => onLanguageChange(e.target.value as 'promql' | 'helix')}
            className="query-editor__language-select"
            aria-label="Query language"
          >
            <option value="promql">PromQL</option>
            <option value="helix">HelixQL</option>
          </select>
          <button
            onClick={handleFormat}
            className="query-editor__toolbar-btn"
            title="Format document (Ctrl+Shift+F)"
            aria-label="Format document"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="15" y1="4" x2="21" y2="4" />
              <line x1="9" y1="20" x2="21" y2="20" />
              <line x1="3" y1="8" x2="21" y2="8" />
              <line x1="9" y1="12" x2="3" y2="12" />
              <line x1="9" y1="16" x2="21" y2="16" />
              <line x1="15" y1="4" x2="9" y2="20" />
            </svg>
          </button>
          <button
            onClick={handleClear}
            className="query-editor__toolbar-btn"
            title="Clear editor"
            aria-label="Clear editor"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
          <button
            onClick={handleCopy}
            className="query-editor__toolbar-btn"
            title="Copy to clipboard"
            aria-label="Copy to clipboard"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
              <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
            </svg>
          </button>
        </div>
        <div className="query-editor__toolbar-right">
          <button
            onClick={handleRun}
            disabled={isExecuting || !value.trim()}
            className="query-editor__run-btn"
            title="Run query (Ctrl+Enter)"
            aria-label="Run query"
          >
            {isExecuting ? (
              <>
                <svg className="loading-spinner loading-spinner--sm mr-2" />
                Running...
              </>
            ) : (
              <>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
                  <polygon points="5 3 19 12 5 21 5 3" />
                </svg>
                Run
              </>
            )}
          </button>
        </div>
      </div>

      <div className="query-editor__container">
        <Editor
          height={height}
          language={language === 'promql' ? 'promql' : 'helix'}
          value={value}
          onChange={handleEditorChange}
          onMount={handleEditorDidMount}
          theme={theme === 'dark' ? 'vs-dark' : 'vs'}
          options={{
            ...MONACO_OPTIONS,
            readOnly,
            lineNumbers: showLineNumbers ? 'on' : 'off',
          }}
          loading={
            <div className="flex items-center justify-center h-full">
              <div className="loading-spinner loading-spinner--lg" />
            </div>
          }
        />
        {!value && placeholder && (
          <div className="query-editor__placeholder">
            {placeholder}
          </div>
        )}
      </div>

      <div className="query-editor__footer">
        <div className="query-editor__footer-left">
          <span className="text-xs text-muted">
            {value.split('\n').length} lines
          </span>
          <span className="text-xs text-muted mx-2">•</span>
          <span className="text-xs text-muted">
            {value.length} characters
          </span>
        </div>
        <div className="query-editor__footer-right">
          <span className="text-xs text-muted">
            Press <kbd className="px-1 py-0.5 bg-tertiary rounded text-xs">Ctrl+Enter</kbd> to run
          </span>
        </div>
      </div>
    </div>
  )
}

export { QueryEditor }