// Helix Web UI - Main application structure
// React + TypeScript observability dashboard

// App - Main application component
export interface AppProps {
  // Dashboard state
  dashboardId?: string;
  // Query editor state
  query?: string;
  // Theme settings
  theme: 'light' | 'dark';
  // View mode
  view: 'dashboard' | 'explorer' | 'alerts';
}

// QueryEditor - Query composition interface
export interface QueryEditorProps {
  value: string;
  onChange: (value: string) => void;
  language: 'promql' | 'helix';
  autoComplete: boolean;
  showFunctions: boolean;
}

// Dashboard - Dashboard container
export interface DashboardProps {
  dashboards: DashboardDef[];
  onSelect: (id: string) => void;
  onCreate: (name: string) => void;
  onDelete: (id: string) => void;
}

// DashboardDef - Dashboard definition
export interface DashboardDef {
  id: string;
  title: string;
  description: string;
  panels: PanelDef[];
  layout: 'grid' | 'sidebar';
  theme?: 'light' | 'dark';
  createdAt: Date;
  updatedAt: Date;
  owner: string;
  isPublic: boolean;
}

// PanelDef - Panel definition for a dashboard
export interface PanelDef {
  id: string;
  title: string;
  type: 'timeseries' | 'stat' | 'heatmap' | 'table' | 'single' | 'html';
  query: string;
  width: number; // in columns (1-24)
  height: number; // in rows
  minimized?: boolean;
  collapsed?: boolean;
  timezone?: string;
  legend?: 'always' | 'hidden' | 'hover';
  yaxis?: 'left' | 'right' | 'center';
  format?: 'short' | 'ms' | 'byte' | 'duration' | 'msprecision';
}

// QueryLanguage - Query language modes
export const QueryLanguage = {
  PromQL: 'promql' as const,
  Helix: 'helix' as const,
} as const;

// QueryResult - Query result type
export interface QueryResult {
  resultType: 'vector' | 'matrix' | 'scalar' | 'string';
  result: any[];
  labels: Record<string, string>;
  error?: string;
  refId?: string;
  stat?: string;
}

// ThemeContext - Theme provider context
export interface ThemeContextValue {
  theme: 'light' | 'dark';
  setTheme: (theme: 'light' | 'dark') => void;
  toggleTheme: () => void;
}

// UseQueryResult - Custom hook result
export interface UseQueryResult<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
  placeholder?: T;
}

// DashboardLayout - Grid layout configuration
export interface DashboardLayout {
  columns: number; // Total grid columns
  rows: number; // Total grid rows
  margin: number; // Margin between panels in pixels
  gutter?: number;
}

// PanelComponent - Base panel component props
export interface PanelComponentProps {
  panel: PanelDef;
  onEdit: (panelId: string, updates: Partial<PanelDef>) => void;
  onDelete: (panelId: string) => void;
  onToggleMinimize?: (panelId: string) => void;
}

// QueryEditorComponent - Query editor component props
export interface QueryEditorComponentProps {
  initialQuery: string;
  onQueryChange: (query: string) => void;
  onExecute: () => Promise<QueryResult>;
  onExecuteAsync?: () => Promise<void>;
  language: typeof QueryLanguage[keyof typeof QueryLanguage];
  showLineNumbers: boolean;
  autoFocus: boolean;
  readOnly: boolean;
}

// ExplorerView - Event explorer props
export interface ExplorerViewProps {
  showSources: boolean;
  showTags: boolean;
  showTimestamps: boolean;
  timeRange: TimeRange;
  onTimeRangeChange: (range: TimeRange) => void;
  showScrollbar: boolean;
  emptyStateMessage?: string;
}

// TimeRange - Time selection range
export interface TimeRange {
  from: Date;
  to: Date;
  preset: '1h' | '24h' | '7d' | '30d' | '90d' | string;
  timezone: string;
  includeTimeOffset: boolean;
}

// Notification - System notification
export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  message: string;
  timestamp: Date;
  duration?: number; // ms before auto-dismiss
  action?: {
    label: string;
    onClick: () => void;
    style?: 'primary' | 'secondary' | 'danger';
  };
}

// AlertManager - Alert management types
export interface Alert {
  id: string;
  label: string;
  severity: 'critical' | 'warning' | 'info' | 'normal';
  status: 'firing' | 'resolved' | 'acknowledged' | 'suppressed';
  startTime: Date;
  endTime?: Date;
  generatorURL: string;
  annotations: Record<string, string>;
  labels: Record<string, string>;
  tags: Record<string, string>;
  updatedAt: Date;
}

// AlertManagerConfig - Alert configuration
export interface AlertManagerConfig {
  groupBy: string[];
  repeatInterval: string; // e.g., "1h", "24h"
  receiver: string;
  routes: Route[];
  inhibitRules: InhibitRule[];
  defaultReceiver?: string;
  sendResolved?: boolean;
}

// Route - Alert routing rule
export interface Route {
  match: MatchCriteria;
  receiver: string;
  continue?: boolean;
  routes?: Route[];
  repeatInterval?: string;
}

// MatchCriteria - Label matching conditions for routing
export interface MatchCriteria {
  alertname?: string;
  severity?: string[];
  source?: string[];
  team?: string[];
  [key: string]: any;
}

// InhibitRule - Alert inhibition rule
export interface InhibitRule {
  source: string;
  target: string;
  equal: string[];
  sourceMatchers?: MatchCriteria;
  targetMatchers?: MatchCriteria;
}

// AlertEvaluation - Result of alert evaluation
export interface AlertEvaluation {
  alertId: string;
  value: string;
  status: 'firing' | 'passing';
  previousStatus: 'firing' | 'passing';
  labels: Record<string, string>;
  annotations: Record<string, string>;
  evaluatedAt: Date;
}

// DSL - Domain specific language for alerts
export const DSL = {
  // Evaluate an alert expression against data
  evaluate: (expression: string, data: any[]): AlertEvaluation => {
    // Placeholder implementation
    return {
      alertId: '',
      value: expression,
      status: 'passing',
      previousStatus: 'passing',
      labels: {},
      annotations: {},
      evaluatedAt: new Date(),
    };
  },

  // Dep-duplicate alerts
  dedup: (alerts: Alert[]): Alert[] => {
    // Simple dedup by ID
    const seen = new Set<string>();
    return alerts.filter(alert => {
      if (seen.has(alert.id)) {
        return false;
      }
      seen.add(alert.id);
      return true;
    });
  },

  // Route alerts based on labels
  route: (alerts: Alert[], config: AlertManagerConfig): Alert[] => {
    // Simple routing - return all alerts for default receiver
    return alerts.filter(alert => {
      // Check if alert matches any route
      return config.routes.some(route => {
        // Simplified match logic
        return true;
      });
    });
  },
}

// Export types for use by other modules
export { 
  AppProps, QueryEditorProps, DashboardProps, DashboardDef, PanelDef, 
  QueryResult, ThemeContextValue, UseQueryResult, DashboardLayout, 
  PanelComponentProps, QueryEditorComponentProps, ExplorerViewProps, 
  TimeRange, Notification, Alert, AlertManagerConfig, Route, 
  InhibitRule, AlertEvaluation, DSL, QueryLanguage 
};

// Helper types for UI operations
export type PanelType = PanelDef['type'];
export type DashboardView = AppProps['view'];
export type QueryLanguageMode = typeof QueryLanguage[keyof typeof QueryLanguage];

// Chart configuration types
export interface ChartConfig {
  type: 'timeseries' | 'bar' | 'stat' | 'heatmap';
  data: ChartData;
  options: ChartOptions;
}

export interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

export interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor?: string | string[];
  borderColor?: string | string[];
  fill?: boolean;
  tension?: number;
}

export interface ChartOptions {
  responsive: boolean;
  maintainAspectRatio: boolean;
  scales: {
    y: {
      beginAtZero?: boolean;
    };
  };
  plugins: {
    legend: {
      display: boolean;
    };
    tooltip: {
      enabled: boolean;
    };
  };
}

// Timeline types
export interface TimelineEvent {
  id: string;
  timestamp: Date;
  source: string;
  level: 'info' | 'warning' | 'error' | 'critical';
  message: string;
  metadata?: Record<string, string>;
}

// Filter types
export interface FilterState {
  query: string;
  hosts: string[];
  regions: string[];
  services: string[];
  severity: string[];
  timeRange: TimeRange;
}

// Pagination types
export interface PaginationInfo {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

// Sort types
export interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}