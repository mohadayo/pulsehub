import request from 'supertest';
import { app } from './app';

describe('Dashboard API', () => {
  describe('GET /health', () => {
    it('should return health status', async () => {
      const res = await request(app).get('/health');
      expect(res.status).toBe(200);
      expect(res.body.status).toBe('ok');
      expect(res.body.service).toBe('PulseHub Dashboard');
      expect(res.body.version).toBe('1.0.0');
    });
  });

  describe('GET /services', () => {
    it('should return list of services', async () => {
      const res = await request(app).get('/services');
      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
      expect(res.body.length).toBe(3);
    });

    it('should include all three services', async () => {
      const res = await request(app).get('/services');
      const names = res.body.map((s: { name: string }) => s.name);
      expect(names).toContain('Analyzer');
      expect(names).toContain('Gateway');
      expect(names).toContain('Dashboard');
    });

    it('should include language info for each service', async () => {
      const res = await request(app).get('/services');
      const languages = res.body.map((s: { language: string }) => s.language);
      expect(languages).toContain('Python');
      expect(languages).toContain('Go');
      expect(languages).toContain('TypeScript');
    });
  });

  describe('GET /overview', () => {
    it('should return dashboard overview', async () => {
      const res = await request(app).get('/overview');
      expect(res.status).toBe(200);
      expect(res.body.services).toBeDefined();
      expect(Array.isArray(res.body.services)).toBe(true);
      expect(res.body.timestamp).toBeDefined();
    });

    it('should report unhealthy when services are down', async () => {
      const res = await request(app).get('/overview');
      expect(res.status).toBe(200);
      for (const svc of res.body.services) {
        expect(svc.status).toBe('unhealthy');
      }
    });
  });

  describe('404 handler', () => {
    it('should return 404 for unknown routes', async () => {
      const res = await request(app).get('/nonexistent');
      expect(res.status).toBe(404);
      expect(res.body.error).toBe('not_found');
    });
  });
});
