import os


class Settings:
    APP_NAME: str = "PulseHub Analyzer"
    APP_VERSION: str = "1.0.0"
    HOST: str = os.getenv("ANALYZER_HOST", "0.0.0.0")
    PORT: int = int(os.getenv("ANALYZER_PORT", "8001"))
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")
    GATEWAY_URL: str = os.getenv("GATEWAY_URL", "http://localhost:8002")


settings = Settings()
