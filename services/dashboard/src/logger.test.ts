import { Logger } from './logger';

describe('Logger', () => {
  let consoleSpy: jest.SpyInstance;

  beforeEach(() => {
    consoleSpy = jest.spyOn(console, 'log').mockImplementation();
  });

  afterEach(() => {
    consoleSpy.mockRestore();
  });

  it('should log info messages', () => {
    const logger = new Logger('test-service', 'INFO');
    logger.info('test message');
    expect(consoleSpy).toHaveBeenCalledTimes(1);
    const logged = JSON.parse(consoleSpy.mock.calls[0][0]);
    expect(logged.level).toBe('INFO');
    expect(logged.message).toBe('test message');
    expect(logged.service).toBe('test-service');
  });

  it('should not log debug when level is INFO', () => {
    const logger = new Logger('test-service', 'INFO');
    logger.debug('debug message');
    expect(consoleSpy).not.toHaveBeenCalled();
  });

  it('should log debug when level is DEBUG', () => {
    const logger = new Logger('test-service', 'DEBUG');
    logger.debug('debug message');
    expect(consoleSpy).toHaveBeenCalledTimes(1);
  });

  it('should include metadata in log output', () => {
    const logger = new Logger('test-service', 'INFO');
    logger.info('test', { key: 'value' });
    const logged = JSON.parse(consoleSpy.mock.calls[0][0]);
    expect(logged.key).toBe('value');
  });

  it('should log errors', () => {
    const logger = new Logger('test-service', 'INFO');
    logger.error('something broke', { code: 500 });
    const logged = JSON.parse(consoleSpy.mock.calls[0][0]);
    expect(logged.level).toBe('ERROR');
    expect(logged.code).toBe(500);
  });

  it('should include timestamp', () => {
    const logger = new Logger('test-service', 'INFO');
    logger.info('test');
    const logged = JSON.parse(consoleSpy.mock.calls[0][0]);
    expect(logged.timestamp).toBeDefined();
  });
});
