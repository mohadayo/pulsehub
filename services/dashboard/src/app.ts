import express, { Request, Response } from 'express';
import { loadConfig } from './config';
import { Logger } from './logger';
import type { HealthResponse, DashboardOverview, ServiceStatus, ApiError } from './types';

const config = loadConfig();
const logger = new Logger('PulseHub Dashboard', config.logLevel as 'INFO');

export const app = express();
app.use(express.json());

app.get('/health', (_req: Request, res: Response) => {
  const response: HealthResponse = {
    status: 'ok',
    service: 'PulseHub Dashboard',
    version: '1.0.0',
  };
  res.json(response);
});

async function checkService(name: string, url: string): Promise<ServiceStatus> {
  const start = Date.now();
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 3000);
    const resp = await fetch(`${url}/health`, { signal: controller.signal });
    clearTimeout(timeout);
    const elapsed = Date.now() - start;

    if (resp.ok) {
      logger.info(`Service ${name} is healthy`, { responseTimeMs: elapsed });
      return { name, url, status: 'healthy', responseTimeMs: elapsed };
    }
    logger.warn(`Service ${name} returned non-OK status`, { statusCode: resp.status });
    return { name, url, status: 'unhealthy', responseTimeMs: elapsed };
  } catch (err) {
    const elapsed = Date.now() - start;
    logger.error(`Service ${name} health check failed`, {
      error: err instanceof Error ? err.message : String(err),
    });
    return { name, url, status: 'unhealthy', responseTimeMs: elapsed };
  }
}

app.get('/overview', async (_req: Request, res: Response) => {
  const services = await Promise.all([
    checkService('Analyzer', config.analyzerUrl),
    checkService('Gateway', config.gatewayUrl),
  ]);

  const overview: DashboardOverview = {
    services,
    timestamp: new Date().toISOString(),
  };

  logger.info('Dashboard overview generated', { serviceCount: services.length });
  res.json(overview);
});

app.get('/services', (_req: Request, res: Response) => {
  const services = [
    {
      name: 'Analyzer',
      description: 'Event analysis and aggregation service',
      url: config.analyzerUrl,
      port: 8001,
      language: 'Python',
      endpoints: ['/health', '/events', '/events/summary', '/events/analyze'],
    },
    {
      name: 'Gateway',
      description: 'Event collection gateway service',
      url: config.gatewayUrl,
      port: 8002,
      language: 'Go',
      endpoints: ['/health', '/events', '/stats'],
    },
    {
      name: 'Dashboard',
      description: 'Dashboard BFF service',
      url: `http://localhost:${config.port}`,
      port: config.port,
      language: 'TypeScript',
      endpoints: ['/health', '/overview', '/services'],
    },
  ];
  res.json(services);
});

app.use((_req: Request, res: Response) => {
  const error: ApiError = {
    error: 'not_found',
    message: 'The requested endpoint does not exist',
  };
  res.status(404).json(error);
});
