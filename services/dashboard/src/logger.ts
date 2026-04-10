export type LogLevel = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR';

const LOG_LEVELS: Record<LogLevel, number> = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
};

export class Logger {
  private level: number;
  private service: string;

  constructor(service: string, level: LogLevel = 'INFO') {
    this.service = service;
    this.level = LOG_LEVELS[level] ?? LOG_LEVELS.INFO;
  }

  private log(level: LogLevel, message: string, meta?: Record<string, unknown>): void {
    if (LOG_LEVELS[level] < this.level) return;
    const entry = {
      timestamp: new Date().toISOString(),
      level,
      service: this.service,
      message,
      ...meta,
    };
    console.log(JSON.stringify(entry));
  }

  debug(message: string, meta?: Record<string, unknown>): void {
    this.log('DEBUG', message, meta);
  }

  info(message: string, meta?: Record<string, unknown>): void {
    this.log('INFO', message, meta);
  }

  warn(message: string, meta?: Record<string, unknown>): void {
    this.log('WARN', message, meta);
  }

  error(message: string, meta?: Record<string, unknown>): void {
    this.log('ERROR', message, meta);
  }
}
