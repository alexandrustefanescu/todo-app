import pg from 'pg';
import logger from '../utils/logger.js';

const { Pool } = pg;

let pool;

export const getPool = () => {
  if (!pool) {
    throw new Error('Database pool not initialized');
  }
  return pool;
};

export const initializeDatabase = async () => {
  const connectionString = process.env.DATABASE_URL;

  if (!connectionString) {
    throw new Error('DATABASE_URL environment variable is not set');
  }

  pool = new Pool({
    connectionString,
    max: 20,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 2000,
  });

  pool.on('error', (err) => {
    logger.error(err, 'Unexpected error on idle client');
  });

  // Test the connection
  try {
    const client = await pool.connect();
    logger.info('Database connection test successful');
    client.release();
  } catch (err) {
    logger.error(err, 'Failed to connect to database');
    throw err;
  }

  return pool;
};

export const closeDatabase = async () => {
  if (pool) {
    await pool.end();
    logger.info('Database pool closed');
  }
};

export const query = async (text, params) => {
  const client = await pool.connect();
  try {
    return await client.query(text, params);
  } finally {
    client.release();
  }
};
