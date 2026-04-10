import pytest
from fastapi.testclient import TestClient
from app.main import app, event_analyzer


@pytest.fixture(autouse=True)
def clear_state():
    event_analyzer.clear()
    yield
    event_analyzer.clear()


@pytest.fixture
def client():
    return TestClient(app)


def test_health(client):
    resp = client.get("/health")
    assert resp.status_code == 200
    data = resp.json()
    assert data["status"] == "ok"
    assert data["service"] == "PulseHub Analyzer"


def test_post_event(client):
    resp = client.post("/events", json={
        "event_type": "click",
        "source": "test",
        "payload": {"button": "submit"},
    })
    assert resp.status_code == 200
    data = resp.json()
    assert data["event_type"] == "click"
    assert data["id"] is not None


def test_get_summary(client):
    client.post("/events", json={
        "event_type": "click", "source": "test"
    })
    resp = client.get("/events/summary")
    assert resp.status_code == 200
    data = resp.json()
    assert data["total_events"] == 1


def test_analyze(client):
    client.post("/events", json={
        "event_type": "click", "source": "test"
    })
    client.post("/events", json={
        "event_type": "purchase", "source": "test"
    })
    resp = client.get("/events/analyze")
    assert resp.status_code == 200
    data = resp.json()
    assert len(data) == 2


def test_clear_events(client):
    client.post("/events", json={
        "event_type": "click", "source": "test"
    })
    resp = client.delete("/events")
    assert resp.status_code == 200
    assert resp.json()["cleared"] == 1
