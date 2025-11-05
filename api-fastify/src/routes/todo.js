import {
  listTodos,
  getTodo,
  createTodo,
  updateTodo,
  deleteTodo,
} from '../handlers/todo.js';
import {
  createTodoSchema,
  updateTodoSchema,
  paramsSchema,
} from '../models/todo.js';

export async function todoRoutes(fastify) {
  // Get all todos
  fastify.get('/', {
    schema: {
      response: {
        200: {
          type: 'array',
          items: {
            type: 'object',
            properties: {
              id: { type: 'string' },
              title: { type: 'string' },
              description: { type: ['string', 'null'] },
              completed: { type: 'boolean' },
              created_at: { type: 'string', format: 'date-time' },
              updated_at: { type: 'string', format: 'date-time' },
            },
          },
        },
      },
    },
    handler: listTodos,
  });

  // Create a new todo
  fastify.post('/', {
    schema: {
      body: createTodoSchema,
      response: {
        201: {
          type: 'object',
          properties: {
            id: { type: 'string' },
            title: { type: 'string' },
            description: { type: ['string', 'null'] },
            completed: { type: 'boolean' },
            created_at: { type: 'string', format: 'date-time' },
            updated_at: { type: 'string', format: 'date-time' },
          },
        },
      },
    },
    handler: createTodo,
  });

  // Get a specific todo
  fastify.get('/:id', {
    schema: {
      params: paramsSchema,
      response: {
        200: {
          type: 'object',
          properties: {
            id: { type: 'string' },
            title: { type: 'string' },
            description: { type: ['string', 'null'] },
            completed: { type: 'boolean' },
            created_at: { type: 'string', format: 'date-time' },
            updated_at: { type: 'string', format: 'date-time' },
          },
        },
      },
    },
    handler: getTodo,
  });

  // Update a todo
  fastify.put('/:id', {
    schema: {
      params: paramsSchema,
      body: updateTodoSchema,
      response: {
        200: {
          type: 'object',
          properties: {
            id: { type: 'string' },
            title: { type: 'string' },
            description: { type: ['string', 'null'] },
            completed: { type: 'boolean' },
            created_at: { type: 'string', format: 'date-time' },
            updated_at: { type: 'string', format: 'date-time' },
          },
        },
      },
    },
    handler: updateTodo,
  });

  // Delete a todo
  fastify.delete('/:id', {
    schema: {
      params: paramsSchema,
      response: {
        204: {},
      },
    },
    handler: deleteTodo,
  });
}
