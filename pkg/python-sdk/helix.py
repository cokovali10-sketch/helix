# Helix Python SDK
# This SDK allows Python applications to interact with the Helix observability platform

import json
import time
import requests
from typing import Dict, List, Optional, Any
from dataclasses import dataclass, field
from datetime import datetime
from urllib.parse import urlencode


@dataclass
class Event:
    """An observability event."""
    timestamp: int
    source: str
    metrics: Dict[str, float] = field(default_factory=dict)
    tags: Dict[str, str] = field(default_factory=dict)
    tenant: str = "default"
    id: Optional[str] = None

    def to_dict(self) -> Dict[str, Any]:
        return {
            "timestamp": self.timestamp,
            "source": self.source,
            "metrics": self.metrics,
            "tags": self.tags,
            "tenant": self.tenant,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Event":
        return cls(
            timestamp=data.get("timestamp", int(time.time())),
            source=data.get("source", "unknown"),
            metrics=data.get("metrics", {}),
            tags=data.get("tags", {}),
            tenant=data.get("tenant", "default"),
            id=data.get("id"),
        )


@dataclass
class QueryOptions:
    """Options for executing a query."""
    expression: str
    tenant: Optional[str] = None
    time_range: Optional[tuple[int, int]] = None
    step: Optional[int] = None


@dataclass
class QueryResult:
    """Result of a single query."""
    metric: Dict[str, str]
    value: float
    timestamp: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "QueryResult":
        return cls(
            metric=data.get("metric", {}),
            value=data.get("value", 0.0),
            timestamp=data.get("timestamp", 0),
        )


@dataclass
class TenantInfo:
    """Tenant information."""
    id: str
    name: str
    rate_limit: int
    quota: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "TenantInfo":
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            rate_limit=data.get("rate_limit", 0),
            quota=data.get("quota", 0),
        )


@dataclass
class Metrics:
    """Platform metrics."""
    events_received: int
    events_processed: int
    active_tenants: int
    memory_usage_bytes: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Metrics":
        return cls(
            events_received=data.get("events_received", 0),
            events_processed=data.get("events_processed", 0),
            active_tenants=data.get("active_tenants", 0),
            memory_usage_bytes=data.get("memory_usage_bytes", 0),
        )


class HelixClient:
    """Main client for interacting with the Helix platform."""

    def __init__(self, base_url: str, api_key: str, timeout: int = 30):
        self.base_url = base_url.rstrip("/")
        self.api_key = api_key
        self.timeout = timeout
        self.session = requests.Session()
        self.session.headers.update({
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}",
        })

    def _request(self, method: str, path: str, **kwargs) -> requests.Response:
        url = f"{self.base_url}{path}"
        response = self.session.request(method, url, timeout=self.timeout, **kwargs)
        response.raise_for_status()
        return response

    def ingest_event(self, event: Event) -> Dict[str, Any]:
        """Ingest an event into the Helix platform."""
        response = self._request("POST", "/v1/events", json=event.to_dict())
        return response.json()

    def ingest_events(self, events: List[Event]) -> Dict[str, Any]:
        """Ingest multiple events."""
        response = self._request("POST", "/v1/events/batch", json=[e.to_dict() for e in events])
        return response.json()

    def query(self, options: QueryOptions) -> List[QueryResult]:
        """Execute a query against the Helix platform."""
        params = {"query": options.expression}
        if options.tenant:
            params["tenant"] = options.tenant
        if options.time_range:
            params["start"] = options.time_range[0]
            params["end"] = options.time_range[1]
        if options.step:
            params["step"] = options.step

        response = self._request("GET", f"/v1/query?{urlencode(params)}")
        data = response.json()
        return [QueryResult.from_dict(r) for r in data.get("result", [])]

    def health_check(self) -> bool:
        """Check the health of the Helix platform."""
        try:
            response = self._request("GET", "/health")
            return response.status_code == 200
        except requests.RequestException:
            return False

    def metrics(self) -> Metrics:
        """Get platform metrics."""
        response = self._request("GET", "/metrics")
        return Metrics.from_dict(response.json())

    def list_tenants(self) -> List[TenantInfo]:
        """List all tenants."""
        response = self._request("GET", "/v1/tenants")
        data = response.json()
        return [TenantInfo.from_dict(t) for t in data.get("tenants", [])]

    def search_events(self, query: str, filters: Optional[Dict[str, str]] = None) -> List[Event]:
        """Search for events with optional filters."""
        params = {"query": query}
        if filters:
            params.update(filters)
        response = self._request("GET", f"/v1/events/search?{urlencode(params)}")
        data = response.json()
        return [Event.from_dict(e) for e in data.get("events", [])]

    def get_event(self, event_id: str) -> Event:
        """Get a single event by ID."""
        response = self._request("GET", f"/v1/events/{event_id}")
        return Event.from_dict(response.json())

    def delete_event(self, event_id: str) -> None:
        """Delete an event by ID."""
        self._request("DELETE", f"/v1/events/{event_id}")

    def create_tenant(self, name: str, rate_limit: int, quota: int) -> TenantInfo:
        """Create a new tenant."""
        response = self._request("POST", "/v1/tenants", json={
            "name": name,
            "rate_limit": rate_limit,
            "quota": quota,
        })
        return TenantInfo.from_dict(response.json())

    def update_tenant(self, tenant_id: str, **kwargs) -> TenantInfo:
        """Update a tenant."""
        response = self._request("PUT", f"/v1/tenants/{tenant_id}", json=kwargs)
        return TenantInfo.from_dict(response.json())

    def delete_tenant(self, tenant_id: str) -> None:
        """Delete a tenant."""
        self._request("DELETE", f"/v1/tenants/{tenant_id}")

    def get_dashboards(self) -> List[Dict[str, Any]]:
        """Get all dashboards."""
        response = self._request("GET", "/v1/dashboards")
        return response.json().get("dashboards", [])

    def create_dashboard(self, title: str, description: str = "", panels: List[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Create a new dashboard."""
        response = self._request("POST", "/v1/dashboards", json={
            "title": title,
            "description": description,
            "panels": panels or [],
        })
        return response.json()

    def get_alerts(self, status: Optional[str] = None) -> List[Dict[str, Any]]:
        """Get alerts, optionally filtered by status."""
        params = {}
        if status:
            params["status"] = status
        response = self._request("GET", f"/v1/alerts?{urlencode(params)}")
        return response.json().get("alerts", [])

    def create_alert_rule(self, name: str, expression: str, duration: int, severity: str, labels: Optional[Dict[str, str]] = None) -> Dict[str, Any]:
        """Create a new alert rule."""
        response = self._request("POST", "/v1/alerts/rules", json={
            "name": name,
            "expression": expression,
            "duration": duration,
            "severity": severity,
            "labels": labels or {},
        })
        return response.json()

    def close(self):
        """Close the client session."""
        self.session.close()


# Convenience functions
def create_client(base_url: str, api_key: str) -> HelixClient:
    """Create a new Helix client."""
    return HelixClient(base_url, api_key)


# Default client for local development
default_client = create_client("http://localhost:8080", "helix-dev-key")


if __name__ == "__main__":
    # Example usage
    client = create_client("http://localhost:8080", "helix-dev-key")
    
    # Ingest an event
    event = Event(
        timestamp=int(time.time()),
        source="soc-agent",
        metrics={"cpu": 50.0, "memory": 75.0},
        tags={"host": "soc-01", "zone": "us-east-1"},
        tenant="soc-team-a",
    )
    try:
        result = client.ingest_event(event)
        print(f"Event ingested: {result}")
    except Exception as e:
        print(f"Failed to ingest event: {e}")
    
    # Execute a query
    try:
        results = client.query(QueryOptions(expression="rate(cpu_usage[30s])", tenant="soc-team-a"))
        print(f"Query results: {results}")
    except Exception as e:
        print(f"Query failed: {e}")
    
    # Health check
    try:
        healthy = client.health_check()
        print(f"Health check: {healthy}")
    except Exception as e:
        print(f"Health check failed: {e}")
    
    client.close()