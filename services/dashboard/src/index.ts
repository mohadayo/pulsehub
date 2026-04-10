import { app } from './app';
import { loadConfig } from './config';
import { Logger } from './logger';

const config = loadConfig();
const logger = new Logger('PulseHub Dashboard', config.logLevel as 'INFO');

app.listen(config.port, config.host, () => {
  logger.info(`PulseHub Dashboard starting on ${config.host}:${config.port}`);
});
