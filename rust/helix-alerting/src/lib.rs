//! Helix Alerting - DSL, Evaluation, Deduplication, and Routing
//! 
//! This module provides the alerting system for the Helix observability platform.
//! Includes a domain-specific language for alert expressions, evaluation engine,
//! de-duplication logic, and routing to notification channels.

use std::collections::BTreeMap;
use std::fmt;
use std::sync::{Arc, Mutex, RwLock};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

// ---------------------------------------------------------------------------
// Alert Types
// ---------------------------------------------------------------------------

/// Alert severity levels
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum AlertSeverity {
    Info,
    Warning,
    Critical,
    Emergency,
}

impl fmt::Display for AlertSeverity {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            AlertSeverity::Info => write!(f, "info"),
            AlertSeverity::Warning => write!(f, "warning"),
            AlertSeverity::Critical => write!(f, "critical"),
            AlertSeverity::Emergency => write!(f, "emergency"),
        }
    }
}

/// Alert status
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum AlertStatus {
    Firing,
    Resolved,
    Acknowledged,
    Suppressed,
    Silenced,
}

impl fmt::Display for AlertStatus {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            AlertStatus::Firing => write!(f, "firing"),
            AlertStatus::Resolved => write!(f, "resolved"),
            AlertStatus::Acknowledged => write!(f, "acknowledged"),
            AlertStatus::Suppressed => write!(f, "suppressed"),
            AlertStatus::Silenced => write!(f, "silenced"),
        }
    }
}

/// Alert represents a single alert instance
#[derive(Debug, Clone)]
pub struct Alert {
    pub id: String,
    pub name: String,
    pub expression: String,
    pub severity: AlertSeverity,
    pub status: AlertStatus,
    pub labels: BTreeMap<String, String>,
    pub annotations: BTreeMap<String, String>,
    pub starts_at: i64,
    pub ends_at: Option<i64>,
    pub updated_at: i64,
    pub fingerprint: String,
    pub generator_url: String,
    pub value: f64,
}

impl Alert {
    pub fn new(name: String, expression: String, severity: AlertSeverity) -> Self {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        Alert {
            id: format!("alert-{}-{}", name, now),
            name,
            expression,
            severity,
            status: AlertStatus::Firing,
            labels: BTreeMap::new(),
            annotations: BTreeMap::new(),
            starts_at: now,
            ends_at: None,
            updated_at: now,
            fingerprint: String::new(),
            generator_url: String::new(),
            value: 0.0,
        }
    }
    
    pub fn with_labels(mut self, labels: BTreeMap<String, String>) -> Self {
        self.labels = labels;
        self.fingerprint = self.compute_fingerprint();
        self
    }
    
    pub fn with_annotations(mut self, annotations: BTreeMap<String, String>) -> Self {
        self.annotations = annotations;
        self
    }
    
    pub fn with_value(mut self, value: f64) -> Self {
        self.value = value;
        self
    }
    
    fn compute_fingerprint(&self) -> String {
        let mut parts: Vec<String> = self.labels.iter()
            .map(|(k, v)| format!("{}={}", k, v))
            .collect();
        parts.sort();
        format!("{}:{}", self.name, parts.join(","))
    }
    
    pub fn is_firing(&self) -> bool {
        self.status == AlertStatus::Firing
    }
    
    pub fn is_resolved(&self) -> bool {
        self.status == AlertStatus::Resolved
    }
    
    pub fn duration(&self) -> Duration {
        let end = self.ends_at.unwrap_or_else(|| {
            SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64
        });
        Duration::from_secs((end - self.starts_at).max(0) as u64)
    }
}

// ---------------------------------------------------------------------------
// Alert Rule
// ---------------------------------------------------------------------------

/// Alert rule definition
#[derive(Debug, Clone)]
pub struct AlertRule {
    pub id: String,
    pub name: String,
    pub expression: String,
    pub duration: Duration,
    pub labels: BTreeMap<String, String>,
    pub annotations: BTreeMap<String, String>,
    pub severity: AlertSeverity,
    pub enabled: bool,
    pub created_at: i64,
    pub updated_at: i64,
    pub evaluation_interval: Duration,
}

impl AlertRule {
    pub fn new(name: String, expression: String) -> Self {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        AlertRule {
            id: format!("rule-{}", name),
            name,
            expression,
            duration: Duration::from_secs(300), // 5 minutes default
            labels: BTreeMap::new(),
            annotations: BTreeMap::new(),
            severity: AlertSeverity::Warning,
            enabled: true,
            created_at: now,
            updated_at: now,
            evaluation_interval: Duration::from_secs(15),
        }
    }
    
    pub fn with_duration(mut self, duration: Duration) -> Self {
        self.duration = duration;
        self
    }
    
    pub fn with_severity(mut self, severity: AlertSeverity) -> Self {
        self.severity = severity;
        self
    }
    
    pub fn with_labels(mut self, labels: BTreeMap<String, String>) -> Self {
        self.labels = labels;
        self
    }
    
    pub fn with_annotations(mut self, annotations: BTreeMap<String, String>) -> Self {
        self.annotations = annotations;
        self
    }
    
    pub fn with_evaluation_interval(mut self, interval: Duration) -> Self {
        self.evaluation_interval = interval;
        self
    }
}

// ---------------------------------------------------------------------------
// Alert DSL
// ---------------------------------------------------------------------------

/// Alert DSL parser
pub struct AlertDSL {
    // In production, would contain parser state
}

impl AlertDSL {
    pub fn new() -> Self {
        AlertDSL {}
    }
    
    /// Parse an alert expression
    pub fn parse(&self, expression: &str) -> Result<AlertExpression, String> {
        // Simplified parser - in production, would use full PromQL-like parser
        let expr = expression.trim();
        
        if expr.contains(">") {
            let parts: Vec<&str> = expr.split('>').collect();
            if parts.len() == 2 {
                return Ok(AlertExpression {
                    metric: parts[0].trim().to_string(),
                    operator: ComparisonOp::GreaterThan,
                    threshold: parts[1].trim().parse().unwrap_or(0.0),
                });
            }
        }
        
        if expr.contains("<") {
            let parts: Vec<&str> = expr.split('<').collect();
            if parts.len() == 2 {
                return Ok(AlertExpression {
                    metric: parts[0].trim().to_string(),
                    operator: ComparisonOp::LessThan,
                    threshold: parts[1].trim().parse().unwrap_or(0.0),
                });
            }
        }
        
        if expr.contains("==") {
            let parts: Vec<&str> = expr.split("==").collect();
            if parts.len() == 2 {
                return Ok(AlertExpression {
                    metric: parts[0].trim().to_string(),
                    operator: ComparisonOp::Equal,
                    threshold: parts[1].trim().parse().unwrap_or(0.0),
                });
            }
        }
        
        Err(format!("Unable to parse expression: {}", expression))
    }
    
    /// Evaluate an alert expression against data
    pub fn evaluate(&self, expr: &AlertExpression, value: f64) -> bool {
        match expr.operator {
            ComparisonOp::GreaterThan => value > expr.threshold,
            ComparisonOp::LessThan => value < expr.threshold,
            ComparisonOp::Equal => (value - expr.threshold).abs() < f64::EPSILON,
            ComparisonOp::NotEqual => (value - expr.threshold).abs() >= f64::EPSILON,
            ComparisonOp::GreaterThanOrEqual => value >= expr.threshold,
            ComparisonOp::LessThanOrEqual => value <= expr.threshold,
        }
    }
}

/// Comparison operators
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum ComparisonOp {
    GreaterThan,
    LessThan,
    Equal,
    NotEqual,
    GreaterThanOrEqual,
    LessThanOrEqual,
}

/// Parsed alert expression
#[derive(Debug, Clone)]
pub struct AlertExpression {
    pub metric: String,
    pub operator: ComparisonOp,
    pub threshold: f64,
}

// ---------------------------------------------------------------------------
// Alert Evaluator
// ---------------------------------------------------------------------------

/// Alert evaluator
pub struct AlertEvaluator {
    rules: RwLock<BTreeMap<String, AlertRule>>,
    alerts: RwLock<BTreeMap<String, Alert>>,
    dsl: AlertDSL,
}

impl AlertEvaluator {
    pub fn new() -> Self {
        AlertEvaluator {
            rules: RwLock::new(BTreeMap::new()),
            alerts: RwLock::new(BTreeMap::new()),
            dsl: AlertDSL::new(),
        }
    }
    
    /// Add a rule
    pub fn add_rule(&self, rule: AlertRule) {
        let mut rules = self.rules.write().unwrap();
        rules.insert(rule.id.clone(), rule);
    }
    
    /// Remove a rule
    pub fn remove_rule(&self, id: &str) {
        let mut rules = self.rules.write().unwrap();
        rules.remove(id);
    }
    
    /// Get all rules
    pub fn rules(&self) -> Vec<AlertRule> {
        let rules = self.rules.read().unwrap();
        rules.values().cloned().collect()
    }
    
    /// Evaluate all rules
    pub fn evaluate_all(&self, data: &BTreeMap<String, f64>) -> Vec<Alert> {
        let rules = self.rules.read().unwrap();
        let mut new_alerts = Vec::new();
        
        for rule in rules.values() {
            if !rule.enabled {
                continue;
            }
            
            // Parse expression
            let expr = match self.dsl.parse(&rule.expression) {
                Ok(e) => e,
                Err(_) => continue,
            };
            
            // Get metric value
            let value = data.get(&expr.metric).copied().unwrap_or(0.0);
            
            // Evaluate
            let firing = self.dsl.evaluate(&expr, value);
            
            if firing {
                let alert = Alert::new(rule.name.clone(), rule.expression.clone(), rule.severity)
                    .with_labels(rule.labels.clone())
                    .with_annotations(rule.annotations.clone())
                    .with_value(value);
                new_alerts.push(alert);
            }
        }
        
        new_alerts
    }
    
    /// Get all firing alerts
    pub fn firing_alerts(&self) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        alerts.values().filter(|a| a.is_firing()).cloned().collect()
    }
    
    /// Get all resolved alerts
    pub fn resolved_alerts(&self) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        alerts.values().filter(|a| a.is_resolved()).cloned().collect()
    }
}

// ---------------------------------------------------------------------------
// Alert Deduplicator
// ---------------------------------------------------------------------------

/// Alert deduplicator
pub struct AlertDeduplicator {
    alerts: RwLock<BTreeMap<String, Alert>>,
    silence_rules: RwLock<Vec<SilenceRule>>,
}

/// Silence rule for suppressing alerts
#[derive(Debug, Clone)]
pub struct SilenceRule {
    pub id: String,
    pub matchers: BTreeMap<String, String>,
    pub starts_at: i64,
    pub ends_at: i64,
    pub created_by: String,
    pub comment: String,
}

impl AlertDeduplicator {
    pub fn new() -> Self {
        AlertDeduplicator {
            alerts: RwLock::new(BTreeMap::new()),
            silence_rules: RwLock::new(Vec::new()),
        }
    }
    
    /// Check if alert should be deduplicated
    pub fn deduplicate(&self, alert: &Alert) -> bool {
        let alerts = self.alerts.read().unwrap();
        
        // Check if an identical alert is already firing
        for existing in alerts.values() {
            if existing.fingerprint == alert.fingerprint && existing.is_firing() {
                return false; // Duplicate
            }
        }
        
        // Check silence rules
        let silences = self.silence_rules.read().unwrap();
        for silence in silences.iter() {
            if self.matches_silence(alert, silence) {
                return false; // Silenced
            }
        }
        
        true // New alert
    }
    
    fn matches_silence(&self, alert: &Alert, silence: &SilenceRule) -> bool {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        if now < silence.starts_at || now > silence.ends_at {
            return false;
        }
        
        for (key, value) in &silence.matchers {
            if alert.labels.get(key) != Some(value) {
                return false;
            }
        }
        true
    }
    
    /// Add an alert
    pub fn add(&self, alert: Alert) {
        let mut alerts = self.alerts.write().unwrap();
        alerts.insert(alert.id.clone(), alert);
    }
    
    /// Add a silence rule
    pub fn add_silence(&self, silence: SilenceRule) {
        let mut silences = self.silence_rules.write().unwrap();
        silences.push(silence);
    }
    
    /// Remove a silence rule
    pub fn remove_silence(&self, id: &str) {
        let mut silences = self.silence_rules.write().unwrap();
        silences.retain(|s| s.id != id);
    }
    
    /// Get all silences
    pub fn silences(&self) -> Vec<SilenceRule> {
        let silences = self.silence_rules.read().unwrap();
        silences.clone()
    }
}

// ---------------------------------------------------------------------------
// Alert Router
// ---------------------------------------------------------------------------

/// Alert router for routing alerts to receivers
pub struct AlertRouter {
    routes: RwLock<Vec<Route>>,
    receivers: RwLock<BTreeMap<String, Receiver>>,
    groups: RwLock<BTreeMap<String, AlertGroup>>,
}

/// Route definition
#[derive(Debug, Clone)]
pub struct Route {
    pub id: String,
    pub matchers: BTreeMap<String, String>,
    pub receiver: String,
    pub continue_routing: bool,
    pub group_by: Vec<String>,
    pub group_wait: Duration,
    pub group_interval: Duration,
    pub repeat_interval: Duration,
    pub child_routes: Vec<Route>,
}

/// Receiver definition
#[derive(Debug, Clone)]
pub struct Receiver {
    pub id: String,
    pub name: String,
    pub receiver_type: ReceiverType,
    pub config: BTreeMap<String, String>,
}

/// Receiver types
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum ReceiverType {
    Email,
    Slack,
    Webhook,
    PagerDuty,
    OpsGenie,
    Telegram,
    Discord,
}

impl fmt::Display for ReceiverType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ReceiverType::Email => write!(f, "email"),
            ReceiverType::Slack => write!(f, "slack"),
            ReceiverType::Webhook => write!(f, "webhook"),
            ReceiverType::PagerDuty => write!(f, "pagerduty"),
            ReceiverType::OpsGenie => write!(f, "opsgenie"),
            ReceiverType::Telegram => write!(f, "telegram"),
            ReceiverType::Discord => write!(f, "discord"),
        }
    }
}

/// Alert group for grouping related alerts
#[derive(Debug, Clone)]
pub struct AlertGroup {
    pub name: String,
    pub labels: BTreeMap<String, String>,
    pub alerts: Vec<Alert>,
    pub created_at: i64,
    pub updated_at: i64,
}

impl AlertRouter {
    pub fn new() -> Self {
        AlertRouter {
            routes: RwLock::new(Vec::new()),
            receivers: RwLock::new(BTreeMap::new()),
            groups: RwLock::new(BTreeMap::new()),
        }
    }
    
    /// Add a route
    pub fn add_route(&self, route: Route) {
        let mut routes = self.routes.write().unwrap();
        routes.push(route);
    }
    
    /// Add a receiver
    pub fn add_receiver(&self, receiver: Receiver) {
        let mut receivers = self.receivers.write().unwrap();
        receivers.insert(receiver.id.clone(), receiver);
    }
    
    /// Route an alert
    pub fn route(&self, alert: &Alert) -> Vec<String> {
        let routes = self.routes.read().unwrap();
        let mut matched_receivers = Vec::new();
        
        for route in routes.iter() {
            if self.matches(alert, route) {
                matched_receivers.push(route.receiver.clone());
                
                if !route.continue_routing {
                    break;
                }
            }
        }
        
        matched_receivers
    }
    
    fn matches(&self, alert: &Alert, route: &Route) -> bool {
        for (key, value) in &route.matchers {
            if alert.labels.get(key) != Some(value) {
                return false;
            }
        }
        true
    }
    
    /// Group alerts by labels
    pub fn group_alerts(&self, alerts: &[Alert], group_by: &[String]) -> BTreeMap<String, AlertGroup> {
        let mut groups: BTreeMap<String, AlertGroup> = BTreeMap::new();
        
        for alert in alerts {
            let mut group_key = String::new();
            for label in group_by {
                if let Some(value) = alert.labels.get(label) {
                    group_key.push_str(&format!("{}={};", label, value));
                }
            }
            
            let group = groups.entry(group_key.clone()).or_insert_with(|| AlertGroup {
                name: group_key.clone(),
                labels: BTreeMap::new(),
                alerts: Vec::new(),
                created_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64,
                updated_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64,
            });
            
            group.alerts.push(alert.clone());
            group.updated_at = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        }
        
        groups
    }
    
    /// Get all receivers
    pub fn receivers(&self) -> Vec<Receiver> {
        let receivers = self.receivers.read().unwrap();
        receivers.values().cloned().collect()
    }
    
    /// Get all routes
    pub fn routes(&self) -> Vec<Route> {
        let routes = self.routes.read().unwrap();
        routes.clone()
    }
}

// ---------------------------------------------------------------------------
// Alert Manager (Composite)
// ---------------------------------------------------------------------------

/// Alert manager combining evaluator, deduplicator, and router
pub struct AlertManager {
    evaluator: AlertEvaluator,
    deduplicator: AlertDeduplicator,
    router: AlertRouter,
    history: AlertHistory,
    rules: RwLock<BTreeMap<String, AlertRule>>,
}

/// Alert history
pub struct AlertHistory {
    alerts: RwLock<Vec<Alert>>,
    max_size: usize,
}

impl AlertHistory {
    pub fn new(max_size: usize) -> Self {
        AlertHistory {
            alerts: RwLock::new(Vec::new()),
            max_size,
        }
    }
    
    pub fn add(&self, alert: Alert) {
        let mut alerts = self.alerts.write().unwrap();
        alerts.push(alert);
        if alerts.len() > self.max_size {
            alerts.remove(0);
        }
    }
    
    pub fn recent(&self, limit: usize) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        if alerts.len() <= limit {
            alerts.clone()
        } else {
            alerts[alerts.len() - limit..].to_vec()
        }
    }
    
    pub fn by_severity(&self, severity: AlertSeverity) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        alerts.iter().filter(|a| a.severity == severity).cloned().collect()
    }
    
    pub fn firing(&self) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        alerts.iter().filter(|a| a.is_firing()).cloned().collect()
    }
    
    pub fn resolved(&self) -> Vec<Alert> {
        let alerts = self.alerts.read().unwrap();
        alerts.iter().filter(|a| a.is_resolved()).cloned().collect()
    }
    
    pub fn count(&self) -> usize {
        let alerts = self.alerts.read().unwrap();
        alerts.len()
    }
}

impl AlertManager {
    pub fn new() -> Self {
        AlertManager {
            evaluator: AlertEvaluator::new(),
            deduplicator: AlertDeduplicator::new(),
            router: AlertRouter::new(),
            history: AlertHistory::new(10000),
            rules: RwLock::new(BTreeMap::new()),
        }
    }
    
    /// Add a rule
    pub fn add_rule(&self, rule: AlertRule) {
        let mut rules = self.rules.write().unwrap();
        rules.insert(rule.id.clone(), rule.clone());
        self.evaluator.add_rule(rule);
    }
    
    /// Remove a rule
    pub fn remove_rule(&self, id: &str) {
        let mut rules = self.rules.write().unwrap();
        rules.remove(id);
        self.evaluator.remove_rule(id);
    }
    
    /// Get all rules
    pub fn rules(&self) -> Vec<AlertRule> {
        let rules = self.rules.read().unwrap();
        rules.values().cloned().collect()
    }
    
    /// Evaluate all rules against data
    pub fn evaluate(&self, data: &BTreeMap<String, f64>) -> Vec<Alert> {
        let new_alerts = self.evaluator.evaluate_all(data);
        let mut fired = Vec::new();
        
        for alert in new_alerts {
            if self.deduplicator.deduplicate(&alert) {
                self.deduplicator.add(alert.clone());
                self.history.add(alert.clone());
                fired.push(alert.clone());
                
                // Route the alert
                let receivers = self.router.route(&alert);
                for receiver in receivers {
                    println!("Routing alert {} to {}", alert.name, receiver);
                }
            }
        }
        
        fired
    }
    
    /// Get firing alerts
    pub fn firing_alerts(&self) -> Vec<Alert> {
        self.history.firing()
    }
    
    /// Get resolved alerts
    pub fn resolved_alerts(&self) -> Vec<Alert> {
        self.history.resolved()
    }
    
    /// Get alert history
    pub fn history(&self) -> &AlertHistory {
        &self.history
    }
    
    /// Add a receiver
    pub fn add_receiver(&self, receiver: Receiver) {
        self.router.add_receiver(receiver);
    }
    
    /// Add a route
    pub fn add_route(&self, route: Route) {
        self.router.add_route(route);
    }
    
    /// Add a silence rule
    pub fn add_silence(&self, silence: SilenceRule) {
        self.deduplicator.add_silence(silence);
    }
    
    /// Remove a silence rule
    pub fn remove_silence(&self, id: &str) {
        self.deduplicator.remove_silence(id);
    }
    
    /// Get all silences
    pub fn silences(&self) -> Vec<SilenceRule> {
        self.deduplicator.silences()
    }
}

// ---------------------------------------------------------------------------
// Default Alert Rules
// ---------------------------------------------------------------------------

/// Default alert rules for common scenarios
pub fn default_alert_rules() -> Vec<AlertRule> {
    vec![
        AlertRule::new("HighCPUUsage".to_string(), "cpu_usage > 90".to_string())
            .with_duration(Duration::from_secs(300))
            .with_severity(AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("severity".to_string(), "warning".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "High CPU usage detected".to_string())])),
        
        AlertRule::new("HighMemoryUsage".to_string(), "memory_usage > 90".to_string())
            .with_duration(Duration::from_secs(300))
            .with_severity(AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("severity".to_string(), "warning".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "High memory usage detected".to_string())])),
        
        AlertRule::new("DiskSpaceLow".to_string(), "disk_usage > 95".to_string())
            .with_duration(Duration::from_secs(600))
            .with_severity(AlertSeverity::Critical)
            .with_labels(BTreeMap::from([("severity".to_string(), "critical".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "Disk space critically low".to_string())])),
        
        AlertRule::new("ServiceDown".to_string(), "up == 0".to_string())
            .with_duration(Duration::from_secs(60))
            .with_severity(AlertSeverity::Critical)
            .with_labels(BTreeMap::from([("severity".to_string(), "critical".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "Service is down".to_string())])),
        
        AlertRule::new("HighErrorRate".to_string(), "error_rate > 5".to_string())
            .with_duration(Duration::from_secs(300))
            .with_severity(AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("severity".to_string(), "warning".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "High error rate detected".to_string())])),
        
        AlertRule::new("HighLatency".to_string(), "latency_p99 > 1000".to_string())
            .with_duration(Duration::from_secs(300))
            .with_severity(AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("severity".to_string(), "warning".to_string())]))
            .with_annotations(BTreeMap::from([("summary".to_string(), "High latency detected".to_string())])),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_alert_creation() {
        let alert = Alert::new("TestAlert".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("host".to_string(), "server1".to_string())]))
            .with_value(95.0);
        
        assert_eq!(alert.name, "TestAlert");
        assert_eq!(alert.severity, AlertSeverity::Warning);
        assert_eq!(alert.status, AlertStatus::Firing);
        assert_eq!(alert.value, 95.0);
        assert!(alert.is_firing());
    }
    
    #[test]
    fn test_alert_rule_creation() {
        let rule = AlertRule::new("TestRule".to_string(), "cpu > 90".to_string())
            .with_severity(AlertSeverity::Critical)
            .with_duration(Duration::from_secs(600));
        
        assert_eq!(rule.name, "TestRule");
        assert_eq!(rule.severity, AlertSeverity::Critical);
        assert_eq!(rule.duration.as_secs(), 600);
        assert!(rule.enabled);
    }
    
    #[test]
    fn test_alert_dsl_parse() {
        let dsl = AlertDSL::new();
        
        let expr = dsl.parse("cpu_usage > 90").unwrap();
        assert_eq!(expr.metric, "cpu_usage");
        assert_eq!(expr.operator, ComparisonOp::GreaterThan);
        assert_eq!(expr.threshold, 90.0);
        
        let expr = dsl.parse("memory_usage < 10").unwrap();
        assert_eq!(expr.metric, "memory_usage");
        assert_eq!(expr.operator, ComparisonOp::LessThan);
        assert_eq!(expr.threshold, 10.0);
    }
    
    #[test]
    fn test_alert_dsl_evaluate() {
        let dsl = AlertDSL::new();
        let expr = dsl.parse("cpu_usage > 90").unwrap();
        
        assert!(dsl.evaluate(&expr, 95.0));
        assert!(!dsl.evaluate(&expr, 85.0));
    }
    
    #[test]
    fn test_alert_evaluator() {
        let evaluator = AlertEvaluator::new();
        let rule = AlertRule::new("TestRule".to_string(), "cpu_usage > 90".to_string());
        evaluator.add_rule(rule);
        
        assert_eq!(evaluator.rules().len(), 1);
        
        let mut data = BTreeMap::new();
        data.insert("cpu_usage".to_string(), 95.0);
        
        let alerts = evaluator.evaluate_all(&data);
        assert_eq!(alerts.len(), 1);
    }
    
    #[test]
    fn test_alert_deduplicator() {
        let dedup = AlertDeduplicator::new();
        
        let alert1 = Alert::new("Test".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("host".to_string(), "server1".to_string())]));
        
        let alert2 = Alert::new("Test".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning)
            .with_labels(BTreeMap::from([("host".to_string(), "server1".to_string())]));
        
        assert!(dedup.deduplicate(&alert1));
        dedup.add(alert1);
        assert!(!dedup.deduplicate(&alert2)); // Duplicate
    }
    
    #[test]
    fn test_alert_router() {
        let router = AlertRouter::new();
        
        let receiver = Receiver {
            id: "slack-1".to_string(),
            name: "Slack Channel".to_string(),
            receiver_type: ReceiverType::Slack,
            config: BTreeMap::new(),
        };
        router.add_receiver(receiver);
        
        let route = Route {
            id: "route-1".to_string(),
            matchers: BTreeMap::new(),
            receiver: "slack-1".to_string(),
            continue_routing: false,
            group_by: vec![],
            group_wait: Duration::from_secs(10),
            group_interval: Duration::from_secs(10),
            repeat_interval: Duration::from_secs(3600),
            child_routes: vec![],
        };
        router.add_route(route);
        
        let alert = Alert::new("Test".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning);
        let receivers = router.route(&alert);
        assert_eq!(receivers.len(), 1);
        assert_eq!(receivers[0], "slack-1");
    }
    
    #[test]
    fn test_alert_manager() {
        let manager = AlertManager::new();
        
        let rule = AlertRule::new("TestRule".to_string(), "cpu_usage > 90".to_string());
        manager.add_rule(rule);
        
        assert_eq!(manager.rules().len(), 1);
        
        let mut data = BTreeMap::new();
        data.insert("cpu_usage".to_string(), 95.0);
        
        let alerts = manager.evaluate(&data);
        assert_eq!(alerts.len(), 1);
        assert_eq!(manager.firing_alerts().len(), 1);
    }
    
    #[test]
    fn test_alert_history() {
        let history = AlertHistory::new(100);
        
        let alert = Alert::new("Test".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning);
        history.add(alert);
        
        assert_eq!(history.count(), 1);
        assert_eq!(history.recent(10).len(), 1);
        assert_eq!(history.firing().len(), 1);
    }
    
    #[test]
    fn test_default_alert_rules() {
        let rules = default_alert_rules();
        assert!(!rules.is_empty());
        assert!(rules.len() >= 6);
    }
    
    #[test]
    fn test_silence_rule() {
        let dedup = AlertDeduplicator::new();
        
        let silence = SilenceRule {
            id: "silence-1".to_string(),
            matchers: BTreeMap::from([("host".to_string(), "server1".to_string())]),
            starts_at: 0,
            ends_at: i64::MAX,
            created_by: "admin".to_string(),
            comment: "Maintenance".to_string(),
        };
        dedup.add_silence(silence);
        
        assert_eq!(dedup.silences().len(), 1);
    }
    
    #[test]
    fn test_alert_grouping() {
        let router = AlertRouter::new();
        
        let alerts = vec![
            Alert::new("Alert1".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning)
                .with_labels(BTreeMap::from([("host".to_string(), "server1".to_string())])),
            Alert::new("Alert2".to_string(), "cpu > 90".to_string(), AlertSeverity::Warning)
                .with_labels(BTreeMap::from([("host".to_string(), "server1".to_string())])),
        ];
        
        let groups = router.group_alerts(&alerts, &["host".to_string()]);
        assert_eq!(groups.len(), 1);
    }
}