export interface HealthResponse {
  status: string;
  service: string;
  version: string;
}

export interface DashboardOverview {
  services: ServiceStatus[];
  timestamp: string;
}

export interface ServiceStatus {
  name: string;
  url: string;
  status: 'healthy' | 'unhealthy' | 'unknown';
  responseTimeMs?: number;
}

export interface EventSummary {
  total_events: number;
  events_by_type: Record<string, number>;
  latest_event_time: string | null;
}

export interface ApiError {
  error: string;
  message: string;
}
