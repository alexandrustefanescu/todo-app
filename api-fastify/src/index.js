import Fastify from 'fastify';
import cors from '@fastify/cors';
import helmet from '@fastify/helmet';
import { todoRoutes } from './routes/todo.js';
import { initializeDatabase } from './db/connection.js';
import logger from './utils/logger.js';

const fastify = Fastify({
  logger: logger,
  bodyLimit: 1048576, // 1MB
  requestTimeout: 30000,
});

// Register plugins
await fastify.register(cors, {
  origin: '*',
  credentials: false,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization'],
});

await fastify.register(helmet, {
  contentSecurityPolicy: false,
});

// Initialize database
try {
  await initializeDatabase();
  logger.info('Database initialized successfully');
} catch (err) {
  logger.error(err, 'Failed to initialize database');
  process.exit(1);
}

// Register routes
await fastify.register(todoRoutes, { prefix: '/api/todos' });

// Health check endpoint
fastify.get('/health', async (request, reply) => {
  return { status: 'ok' };
});

// 404 handler
fastify.setNotFoundHandler((request, reply) => {
  reply.status(404).send({
    error: 'NOT_FOUND',
    message: 'Endpoint not found',
  });
});

// Error handler
fastify.setErrorHandler((error, request, reply) => {
  logger.error(error, 'Request error');

  if (error.validation) {
    return reply.status(400).send({
      error: 'BAD_REQUEST',
      message: error.message,
    });
  }

  reply.status(error.statusCode || 500).send({
    error: error.code || 'INTERNAL_SERVER_ERROR',
    message: error.message || 'An unexpected error occurred',
  });
});

// Graceful shutdown
const signals = ['SIGINT', 'SIGTERM'];
signals.forEach((signal) => {
  process.on(signal, async () => {
    logger.info(`Received ${signal}, shutting down gracefully`);
    await fastify.close();
    process.exit(0);
  });
});

// Start server
const start = async () => {
  try {
    const port = process.env.PORT || 8080;
    const host = process.env.HOST || '127.0.0.1';

    await fastify.listen({ port, host });
    logger.info(`Server listening on http://${host}:${port}`);
  } catch (err) {
    logger.error(err);
    process.exit(1);
  }
};

start();
