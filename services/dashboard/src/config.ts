export interface Config {
  port: number;
  host: string;
  logLevel: string;
  analyzerUrl: string;
  gatewayUrl: string;
}

export function loadConfig(): Config {
  return {
    port: parseInt(process.env.DASHBOARD_PORT || '8003', 10),
    host: process.env.DASHBOARD_HOST || '0.0.0.0',
    logLevel: process.env.LOG_LEVEL || 'INFO',
    analyzerUrl: process.env.ANALYZER_URL || 'http://localhost:8001',
    gatewayUrl: process.env.GATEWAY_URL || 'http://localhost:8002',
  };
}
