import { v4 as uuidv4 } from 'uuid';
import { query } from '../db/connection.js';
import { Todo } from '../models/todo.js';
import logger from '../utils/logger.js';

// List all todos
export const listTodos = async (request, reply) => {
  try {
    const result = await query(
      'SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC'
    );

    const todos = result.rows.map(
      (row) =>
        new Todo(
          row.id,
          row.title,
          row.description,
          row.completed,
          row.created_at,
          row.updated_at
        )
    );

    return reply.send(todos);
  } catch (err) {
    logger.error(err, 'Error listing todos');
    return reply.status(500).send({
      error: 'INTERNAL_SERVER_ERROR',
      message: 'Failed to fetch todos',
    });
  }
};

// Get a specific todo
export const getTodo = async (request, reply) => {
  try {
    const { id } = request.params;

    // Validate UUID format
    if (!isValidUUID(id)) {
      return reply.status(400).send({
        error: 'BAD_REQUEST',
        message: 'Invalid todo ID format',
      });
    }

    const result = await query(
      'SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1',
      [id]
    );

    if (result.rows.length === 0) {
      return reply.status(404).send({
        error: 'NOT_FOUND',
        message: 'Todo not found',
      });
    }

    const row = result.rows[0];
    const todo = new Todo(
      row.id,
      row.title,
      row.description,
      row.completed,
      row.created_at,
      row.updated_at
    );

    return reply.send(todo);
  } catch (err) {
    logger.error(err, 'Error getting todo');
    return reply.status(500).send({
      error: 'INTERNAL_SERVER_ERROR',
      message: 'Failed to fetch todo',
    });
  }
};

// Create a new todo
export const createTodo = async (request, reply) => {
  try {
    const { title, description } = request.body;

    // Validate title
    if (!title || title.trim() === '') {
      return reply.status(400).send({
        error: 'BAD_REQUEST',
        message: 'Title is required and cannot be empty',
      });
    }

    const id = uuidv4();
    const now = new Date();

    const result = await query(
      `INSERT INTO todos (id, title, description, completed, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING id, title, description, completed, created_at, updated_at`,
      [id, title, description || null, false, now, now]
    );

    const row = result.rows[0];
    const todo = new Todo(
      row.id,
      row.title,
      row.description,
      row.completed,
      row.created_at,
      row.updated_at
    );

    return reply.status(201).send(todo);
  } catch (err) {
    logger.error(err, 'Error creating todo');
    return reply.status(500).send({
      error: 'INTERNAL_SERVER_ERROR',
      message: 'Failed to create todo',
    });
  }
};

// Update a todo
export const updateTodo = async (request, reply) => {
  try {
    const { id } = request.params;
    const { title, description, completed } = request.body;

    // Validate UUID format
    if (!isValidUUID(id)) {
      return reply.status(400).send({
        error: 'BAD_REQUEST',
        message: 'Invalid todo ID format',
      });
    }

    // Check if todo exists
    const existsResult = await query('SELECT id FROM todos WHERE id = $1', [id]);
    if (existsResult.rows.length === 0) {
      return reply.status(404).send({
        error: 'NOT_FOUND',
        message: 'Todo not found',
      });
    }

    const now = new Date();

    // Build dynamic update query
    const updates = [];
    const params = [id];
    let paramCount = 2;

    if (title !== undefined && title !== null) {
      updates.push(`title = $${paramCount}`);
      params.push(title);
      paramCount++;
    }

    if (description !== undefined) {
      updates.push(`description = $${paramCount}`);
      params.push(description);
      paramCount++;
    }

    if (completed !== undefined && completed !== null) {
      updates.push(`completed = $${paramCount}`);
      params.push(completed);
      paramCount++;
    }

    updates.push(`updated_at = $${paramCount}`);
    params.push(now);

    const updateQuery = `
      UPDATE todos
      SET ${updates.join(', ')}
      WHERE id = $1
      RETURNING id, title, description, completed, created_at, updated_at
    `;

    const result = await query(updateQuery, params);

    const row = result.rows[0];
    const todo = new Todo(
      row.id,
      row.title,
      row.description,
      row.completed,
      row.created_at,
      row.updated_at
    );

    return reply.send(todo);
  } catch (err) {
    logger.error(err, 'Error updating todo');
    return reply.status(500).send({
      error: 'INTERNAL_SERVER_ERROR',
      message: 'Failed to update todo',
    });
  }
};

// Delete a todo
export const deleteTodo = async (request, reply) => {
  try {
    const { id } = request.params;

    // Validate UUID format
    if (!isValidUUID(id)) {
      return reply.status(400).send({
        error: 'BAD_REQUEST',
        message: 'Invalid todo ID format',
      });
    }

    // Check if todo exists
    const existsResult = await query('SELECT id FROM todos WHERE id = $1', [id]);
    if (existsResult.rows.length === 0) {
      return reply.status(404).send({
        error: 'NOT_FOUND',
        message: 'Todo not found',
      });
    }

    await query('DELETE FROM todos WHERE id = $1', [id]);

    return reply.status(204).send();
  } catch (err) {
    logger.error(err, 'Error deleting todo');
    return reply.status(500).send({
      error: 'INTERNAL_SERVER_ERROR',
      message: 'Failed to delete todo',
    });
  }
};

// Helper function to validate UUID
function isValidUUID(uuid) {
  const uuidRegex =
    /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(uuid);
}
