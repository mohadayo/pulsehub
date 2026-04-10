import logging
import sys

from fastapi import FastAPI, HTTPException

from .config import settings
from .models import Event, EventSummary, AnalysisResult, HealthResponse
from .analyzer import EventAnalyzer

logging.basicConfig(
    level=getattr(logging, settings.LOG_LEVEL.upper(), logging.INFO),
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    stream=sys.stdout,
)
logger = logging.getLogger(__name__)

app = FastAPI(title=settings.APP_NAME, version=settings.APP_VERSION)
event_analyzer = EventAnalyzer()


@app.get("/health", response_model=HealthResponse)
async def health():
    return HealthResponse(
        status="ok",
        service=settings.APP_NAME,
        version=settings.APP_VERSION,
    )


@app.post("/events", response_model=Event)
async def ingest_event(event: Event):
    try:
        result = event_analyzer.ingest(event)
        return result
    except Exception as e:
        logger.error("Failed to ingest event: %s", str(e))
        raise HTTPException(status_code=500, detail=f"Ingestion failed: {e}")


@app.get("/events/summary", response_model=EventSummary)
async def get_summary():
    return event_analyzer.get_summary()


@app.get("/events/analyze", response_model=list[AnalysisResult])
async def analyze_events():
    return event_analyzer.analyze()


@app.delete("/events")
async def clear_events():
    count = event_analyzer.clear()
    return {"cleared": count}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host=settings.HOST, port=settings.PORT)
