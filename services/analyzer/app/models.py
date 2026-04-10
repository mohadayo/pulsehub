from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime
from enum import Enum


class EventType(str, Enum):
    CLICK = "click"
    PAGE_VIEW = "page_view"
    PURCHASE = "purchase"
    SIGNUP = "signup"
    ERROR = "error"


class Event(BaseModel):
    id: Optional[str] = None
    event_type: EventType
    source: str
    payload: dict = Field(default_factory=dict)
    timestamp: Optional[datetime] = None


class EventSummary(BaseModel):
    total_events: int
    events_by_type: dict[str, int]
    latest_event_time: Optional[str] = None


class AnalysisResult(BaseModel):
    event_type: str
    count: int
    percentage: float


class HealthResponse(BaseModel):
    status: str
    service: str
    version: str
