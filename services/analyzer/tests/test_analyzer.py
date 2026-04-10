import pytest
from app.analyzer import EventAnalyzer
from app.models import Event, EventType


@pytest.fixture
def analyzer():
    return EventAnalyzer()


def test_ingest_assigns_id(analyzer):
    event = Event(event_type=EventType.CLICK, source="web")
    result = analyzer.ingest(event)
    assert result.id is not None


def test_ingest_assigns_timestamp(analyzer):
    event = Event(event_type=EventType.PAGE_VIEW, source="mobile")
    result = analyzer.ingest(event)
    assert result.timestamp is not None


def test_summary_empty(analyzer):
    summary = analyzer.get_summary()
    assert summary.total_events == 0
    assert summary.events_by_type == {}
    assert summary.latest_event_time is None


def test_summary_with_events(analyzer):
    analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    analyzer.ingest(Event(event_type=EventType.PURCHASE, source="api"))

    summary = analyzer.get_summary()
    assert summary.total_events == 3
    assert summary.events_by_type["click"] == 2
    assert summary.events_by_type["purchase"] == 1


def test_analyze_empty(analyzer):
    results = analyzer.analyze()
    assert results == []


def test_analyze_percentages(analyzer):
    analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    analyzer.ingest(Event(event_type=EventType.SIGNUP, source="web"))
    analyzer.ingest(Event(event_type=EventType.ERROR, source="api"))

    results = analyzer.analyze()
    assert len(results) == 3
    click_result = next(r for r in results if r.event_type == "click")
    assert click_result.count == 2
    assert click_result.percentage == 50.0


def test_clear(analyzer):
    analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    analyzer.ingest(Event(event_type=EventType.PURCHASE, source="api"))
    cleared = analyzer.clear()
    assert cleared == 2
    assert analyzer.get_summary().total_events == 0


def test_analyze_sorted_by_count(analyzer):
    for _ in range(5):
        analyzer.ingest(Event(event_type=EventType.CLICK, source="web"))
    for _ in range(3):
        analyzer.ingest(Event(event_type=EventType.PAGE_VIEW, source="web"))
    analyzer.ingest(Event(event_type=EventType.ERROR, source="api"))

    results = analyzer.analyze()
    assert results[0].event_type == "click"
    assert results[1].event_type == "page_view"
    assert results[2].event_type == "error"
