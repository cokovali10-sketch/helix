import React, { useState } from 'react'
import { useTheme } from '../../context/ThemeContext'
import { useAuth } from '../../context/AuthContext'
import { useToast } from '../../components/Common/ToastContainer'

interface SettingsPageProps {}

export function SettingsPage() {
  const { theme, toggleTheme } = useTheme()
  const { user } = useAuth()
  const { addToast } = useToast()
  const [activeTab, setActiveTab] = useState<'general' | 'profile' | 'notifications' | 'api' | 'advanced'>('general')
  const [settings, setSettings] = useState({
    autoRefresh: true,
    refreshInterval: 30,
    defaultTimeRange: '1h',
    defaultTenant: 'tenant-1',
    showTimestamps: true,
    showTags: true,
    enableNotifications: true,
    notificationSound: true,
    emailNotifications: false,
    slackNotifications: false,
    pagerdutyNotifications: false,
    apiEndpoint: 'http://localhost:8080',
    apiKey: 'helix-dev-key',
    debugMode: false,
    verboseLogging: false,
    maxQueryResults: 1000,
    queryTimeout: 30,
    enableCaching: true,
    cacheTTL: 300,
  })

  const handleSave = () => {
    addToast({ type: 'success', title: 'Settings saved' })
  }

  const handleReset = () => {
    if (confirm('Reset all settings to defaults?')) {
      setSettings({
        autoRefresh: true,
        refreshInterval: 30,
        defaultTimeRange: '1h',
        defaultTenant: 'tenant-1',
        showTimestamps: true,
        showTags: true,
        enableNotifications: true,
        notificationSound: true,
        emailNotifications: false,
        slackNotifications: false,
        pagerdutyNotifications: false,
        apiEndpoint: 'http://localhost:8080',
        apiKey: 'helix-dev-key',
        debugMode: false,
        verboseLogging: false,
        maxQueryResults: 1000,
        queryTimeout: 30,
        enableCaching: true,
        cacheTTL: 300,
      })
      addToast({ type: 'info', title: 'Settings reset to defaults' })
    }
  }

  const tabs = [
    { id: 'general', label: 'General', icon: '⚙️' },
    { id: 'profile', label: 'Profile', icon: '👤' },
    { id: 'notifications', label: 'Notifications', icon: '🔔' },
    { id: 'api', label: 'API', icon: '🔌' },
    { id: 'advanced', label: 'Advanced', icon: '🔧' },
  ]

  return (
    <div className="h-full flex flex-col bg-primary">
      <div className="flex flex-wrap items-center justify-between gap-4 p-4 border-b border-border bg-secondary">
        <div>
          <h1 className="text-2xl font-bold">Settings</h1>
          <p className="text-sm text-muted">Configure your Helix platform</p>
        </div>
        <div className="flex items-center gap-3">
          <button className="btn btn--secondary" onClick={handleReset}>Reset to Defaults</button>
          <button className="btn btn--primary" onClick={handleSave}>
            <svg className="w-4 h-4 mr-2" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
            Save Changes
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-hidden flex">
        {/* Sidebar tabs */}
        <div className="w-48 border-r border-border bg-secondary p-2 flex-shrink-0 hidden sm:block">
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors mb-1 ${activeTab === tab.id ? 'bg-primary-light text-primary' : 'text-secondary hover:bg-tertiary hover:text-primary'}`}
            >
              <span>{tab.icon}</span>
              <span>{tab.label}</span>
            </button>
          ))}
        </div>

        {/* Mobile tabs */}
        <div className="sm:hidden border-b border-border bg-secondary p-2 flex gap-1 overflow-x-auto">
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`px-3 py-1.5 rounded text-sm font-medium whitespace-nowrap ${activeTab === tab.id ? 'bg-primary text-white' : 'text-muted'}`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content */}
        <div className="flex-1 overflow-auto p-6">
          {activeTab === 'general' && (
            <div className="max-w-2xl space-y-6">
              <section>
                <h2 className="text-lg font-semibold mb-4">General Settings</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Auto-refresh</p>
                      <p className="text-sm text-muted">Automatically refresh data</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.autoRefresh} onChange={e => setSettings({...settings, autoRefresh: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  {settings.autoRefresh && (
                    <div>
                      <label className="label">Refresh Interval (seconds)</label>
                      <input type="number" className="input" value={settings.refreshInterval} onChange={e => setSettings({...settings, refreshInterval: Number(e.target.value)})} min={5} max={300} />
                    </div>
                  )}
                  <div>
                    <label className="label">Default Time Range</label>
                    <select className="input" value={settings.defaultTimeRange} onChange={e => setSettings({...settings, defaultTimeRange: e.target.value})}>
                      <option value="5m">Last 5 minutes</option>
                      <option value="15m">Last 15 minutes</option>
                      <option value="1h">Last 1 hour</option>
                      <option value="6h">Last 6 hours</option>
                      <option value="24h">Last 24 hours</option>
                      <option value="7d">Last 7 days</option>
                    </select>
                  </div>
                  <div>
                    <label className="label">Default Tenant</label>
                    <select className="input" value={settings.defaultTenant} onChange={e => setSettings({...settings, defaultTenant: e.target.value})}>
                      <option value="tenant-1">SOC Team Alpha</option>
                      <option value="tenant-2">SOC Team Beta</option>
                      <option value="tenant-3">Platform Team</option>
                    </select>
                  </div>
                </div>
              </section>

              <section>
                <h2 className="text-lg font-semibold mb-4">Display Settings</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Show Timestamps</p>
                      <p className="text-sm text-muted">Display timestamps in event explorer</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.showTimestamps} onChange={e => setSettings({...settings, showTimestamps: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Show Tags</p>
                      <p className="text-sm text-muted">Display tags in event explorer</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.showTags} onChange={e => setSettings({...settings, showTags: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                </div>
              </section>
            </div>
          )}

          {activeTab === 'profile' && (
            <div className="max-w-2xl space-y-6">
              <section>
                <h2 className="text-lg font-semibold mb-4">Profile Settings</h2>
                <div className="space-y-4">
                  <div>
                    <label className="label">Name</label>
                    <input type="text" className="input" value={user?.name || ''} readOnly />
                  </div>
                  <div>
                    <label className="label">Email</label>
                    <input type="email" className="input" value={user?.email || ''} readOnly />
                  </div>
                  <div>
                    <label className="label">Role</label>
                    <input type="text" className="input" value={user?.role || ''} readOnly />
                  </div>
                  <div>
                    <label className="label">Tenant</label>
                    <input type="text" className="input" value={user?.tenantName || ''} readOnly />
                  </div>
                  <button className="btn btn--primary">Change Password</button>
                </div>
              </section>
            </div>
          )}

          {activeTab === 'notifications' && (
            <div className="max-w-2xl space-y-6">
              <section>
                <h2 className="text-lg font-semibold mb-4">Notification Settings</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Enable Notifications</p>
                      <p className="text-sm text-muted">Receive alert notifications</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.enableNotifications} onChange={e => setSettings({...settings, enableNotifications: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Notification Sound</p>
                      <p className="text-sm text-muted">Play sound for new notifications</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.notificationSound} onChange={e => setSettings({...settings, notificationSound: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Email Notifications</p>
                      <p className="text-sm text-muted">Receive notifications via email</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.emailNotifications} onChange={e => setSettings({...settings, emailNotifications: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Slack Notifications</p>
                      <p className="text-sm text-muted">Receive notifications via Slack</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.slackNotifications} onChange={e => setSettings({...settings, slackNotifications: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">PagerDuty Notifications</p>
                      <p className="text-sm text-muted">Receive notifications via PagerDuty</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.pagerdutyNotifications} onChange={e => setSettings({...settings, pagerdutyNotifications: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                </div>
              </section>
            </div>
          )}

          {activeTab === 'api' && (
            <div className="max-w-2xl space-y-6">
              <section>
                <h2 className="text-lg font-semibold mb-4">API Settings</h2>
                <div className="space-y-4">
                  <div>
                    <label className="label">API Endpoint</label>
                    <input type="text" className="input" value={settings.apiEndpoint} onChange={e => setSettings({...settings, apiEndpoint: e.target.value})} />
                  </div>
                  <div>
                    <label className="label">API Key</label>
                    <input type="password" className="input" value={settings.apiKey} onChange={e => setSettings({...settings, apiKey: e.target.value})} />
                  </div>
                  <div>
                    <label className="label">Max Query Results</label>
                    <input type="number" className="input" value={settings.maxQueryResults} onChange={e => setSettings({...settings, maxQueryResults: Number(e.target.value)})} min={100} max={10000} />
                  </div>
                  <div>
                    <label className="label">Query Timeout (seconds)</label>
                    <input type="number" className="input" value={settings.queryTimeout} onChange={e => setSettings({...settings, queryTimeout: Number(e.target.value)})} min={5} max={300} />
                  </div>
                </div>
              </section>

              <section>
                <h2 className="text-lg font-semibold mb-4">Cache Settings</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Enable Caching</p>
                      <p className="text-sm text-muted">Cache query results</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.enableCaching} onChange={e => setSettings({...settings, enableCaching: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  {settings.enableCaching && (
                    <div>
                      <label className="label">Cache TTL (seconds)</label>
                      <input type="number" className="input" value={settings.cacheTTL} onChange={e => setSettings({...settings, cacheTTL: Number(e.target.value)})} min={60} max={3600} />
                    </div>
                  )}
                </div>
              </section>
            </div>
          )}

          {activeTab === 'advanced' && (
            <div className="max-w-2xl space-y-6">
              <section>
                <h2 className="text-lg font-semibold mb-4">Advanced Settings</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Debug Mode</p>
                      <p className="text-sm text-muted">Enable debug logging</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.debugMode} onChange={e => setSettings({...settings, debugMode: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">Verbose Logging</p>
                      <p className="text-sm text-muted">Enable verbose logging</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" className="sr-only peer" checked={settings.verboseLogging} onChange={e => setSettings({...settings, verboseLogging: e.target.checked})} />
                      <div className="w-11 h-6 bg-tertiary peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                    </label>
                  </div>
                </div>
              </section>

              <section>
                <h2 className="text-lg font-semibold mb-4 text-error">Danger Zone</h2>
                <div className="space-y-4">
                  <div className="p-4 border border-error rounded-lg bg-error-light">
                    <p className="font-medium text-error mb-2">Clear All Data</p>
                    <p className="text-sm text-muted mb-3">This will permanently delete all events, metrics, and dashboards.</p>
                    <button className="btn btn--danger">Clear All Data</button>
                  </div>
                  <div className="p-4 border border-error rounded-lg bg-error-light">
                    <p className="font-medium text-error mb-2">Reset Platform</p>
                    <p className="text-sm text-muted mb-3">This will reset the platform to factory defaults.</p>
                    <button className="btn btn--danger">Reset Platform</button>
                  </div>
                </div>
              </section>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export { SettingsPage }