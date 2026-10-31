# Helix Python SDK Tests
import unittest
import json
from unittest.mock import patch, MagicMock
from helix import HelixClient, Event, QueryOptions, QueryResult, TenantInfo, Metrics


class TestEvent(unittest.TestCase):
    def test_event_creation(self):
        event = Event(
            timestamp=1700000000,
            source="test-source",
            metrics={"cpu": 50.0},
            tags={"host": "test-host"},
            tenant="test-tenant",
        )
        self.assertEqual(event.timestamp, 1700000000)
        self.assertEqual(event.source, "test-source")
        self.assertEqual(event.metrics["cpu"], 50.0)
        self.assertEqual(event.tags["host"], "test-host")
        self.assertEqual(event.tenant, "test-tenant")

    def test_event_to_dict(self):
        event = Event(
            timestamp=1700000000,
            source="test-source",
            metrics={"cpu": 50.0},
            tags={"host": "test-host"},
            tenant="test-tenant",
        )
        d = event.to_dict()
        self.assertEqual(d["timestamp"], 1700000000)
        self.assertEqual(d["source"], "test-source")
        self.assertEqual(d["metrics"]["cpu"], 50.0)
        self.assertEqual(d["tags"]["host"], "test-host")
        self.assertEqual(d["tenant"], "test-tenant")

    def test_event_from_dict(self):
        data = {
            "timestamp": 1700000000,
            "source": "test-source",
            "metrics": {"cpu": 50.0},
            "tags": {"host": "test-host"},
            "tenant": "test-tenant",
            "id": "evt-123",
        }
        event = Event.from_dict(data)
        self.assertEqual(event.timestamp, 1700000000)
        self.assertEqual(event.source, "test-source")
        self.assertEqual(event.id, "evt-123")

    def test_event_default_tenant(self):
        event = Event(timestamp=1700000000, source="test")
        self.assertEqual(event.tenant, "default")


class TestQueryResult(unittest.TestCase):
    def test_query_result_from_dict(self):
        data = {
            "metric": {"host": "server1"},
            "value": 50.0,
            "timestamp": 1700000000,
        }
        result = QueryResult.from_dict(data)
        self.assertEqual(result.metric["host"], "server1")
        self.assertEqual(result.value, 50.0)
        self.assertEqual(result.timestamp, 1700000000)


class TestTenantInfo(unittest.TestCase):
    def test_tenant_info_from_dict(self):
        data = {
            "id": "tenant-1",
            "name": "Test Tenant",
            "rate_limit": 1000,
            "quota": 10000,
        }
        tenant = TenantInfo.from_dict(data)
        self.assertEqual(tenant.id, "tenant-1")
        self.assertEqual(tenant.name, "Test Tenant")
        self.assertEqual(tenant.rate_limit, 1000)
        self.assertEqual(tenant.quota, 10000)


class TestMetrics(unittest.TestCase):
    def test_metrics_from_dict(self):
        data = {
            "events_received": 1000,
            "events_processed": 950,
            "active_tenants": 2,
            "memory_usage_bytes": 1024000,
        }
        metrics = Metrics.from_dict(data)
        self.assertEqual(metrics.events_received, 1000)
        self.assertEqual(metrics.events_processed, 950)
        self.assertEqual(metrics.active_tenants, 2)
        self.assertEqual(metrics.memory_usage_bytes, 1024000)


class TestHelixClient(unittest.TestCase):
    def setUp(self):
        self.client = HelixClient("http://localhost:8080", "test-key")

    @patch("helix.requests.Session")
    def test_health_check(self, mock_session):
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_session.return_value.request.return_value = mock_response
        self.assertTrue(self.client.health_check())

    @patch("helix.requests.Session")
    def test_ingest_event(self, mock_session):
        mock_response = MagicMock()
        mock_response.json.return_value = {"success": True, "eventId": "evt-123"}
        mock_session.return_value.request.return_value = mock_response
        event = Event(timestamp=1700000000, source="test", metrics={"cpu": 50.0})
        result = self.client.ingest_event(event)
        self.assertTrue(result["success"])

    def test_create_client(self):
        client = create_client("http://localhost:8080", "test-key")
        self.assertEqual(client.base_url, "http://localhost:8080")
        self.assertEqual(client.api_key, "test-key")

    def test_default_client(self):
        self.assertEqual(default_client.base_url, "http://localhost:8080")


if __name__ == "__main__":
    unittest.main()