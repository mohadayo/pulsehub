import logging
from datetime import datetime, timezone
from uuid import uuid4

from .models import Event, EventSummary, AnalysisResult

logger = logging.getLogger(__name__)


class EventAnalyzer:
    def __init__(self):
        self._events: list[Event] = []

    def ingest(self, event: Event) -> Event:
        if event.id is None:
            event.id = str(uuid4())
        if event.timestamp is None:
            event.timestamp = datetime.now(timezone.utc)
        self._events.append(event)
        logger.info(
            "Ingested event id=%s type=%s source=%s",
            event.id, event.event_type, event.source,
        )
        return event

    def get_summary(self) -> EventSummary:
        events_by_type: dict[str, int] = {}
        for ev in self._events:
            key = ev.event_type.value
            events_by_type[key] = events_by_type.get(key, 0) + 1

        latest = None
        if self._events:
            last = self._events[-1]
            if last.timestamp:
                latest = last.timestamp.isoformat()

        logger.info("Generated summary: %d total events", len(self._events))
        return EventSummary(
            total_events=len(self._events),
            events_by_type=events_by_type,
            latest_event_time=latest,
        )

    def analyze(self) -> list[AnalysisResult]:
        summary = self.get_summary()
        total = summary.total_events
        if total == 0:
            return []

        results = []
        for event_type, count in summary.events_by_type.items():
            pct = round((count / total) * 100, 2)
            results.append(AnalysisResult(
                event_type=event_type,
                count=count,
                percentage=pct,
            ))
        results.sort(key=lambda r: r.count, reverse=True)
        logger.info("Analysis complete: %d event types", len(results))
        return results

    def clear(self) -> int:
        count = len(self._events)
        self._events.clear()
        logger.info("Cleared %d events", count)
        return count
