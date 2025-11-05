import Joi from 'joi';

// Schemas for request validation
export const createTodoSchema = {
  type: 'object',
  required: ['title'],
  properties: {
    title: { type: 'string', minLength: 1, maxLength: 255 },
    description: { type: ['string', 'null'], maxLength: 5000 },
  },
};

export const updateTodoSchema = {
  type: 'object',
  properties: {
    title: { type: 'string', minLength: 1, maxLength: 255 },
    description: { type: ['string', 'null'], maxLength: 5000 },
    completed: { type: 'boolean' },
  },
};

export const paramsSchema = {
  type: 'object',
  required: ['id'],
  properties: {
    id: { type: 'string', format: 'uuid' },
  },
};

// Joi schemas for advanced validation
export const createTodoJoi = Joi.object({
  title: Joi.string().min(1).max(255).required(),
  description: Joi.string().max(5000).allow(null).optional(),
});

export const updateTodoJoi = Joi.object({
  title: Joi.string().min(1).max(255).optional(),
  description: Joi.string().max(5000).allow(null).optional(),
  completed: Joi.boolean().optional(),
});

// Todo response class
export class Todo {
  constructor(id, title, description, completed, createdAt, updatedAt) {
    this.id = id;
    this.title = title;
    this.description = description;
    this.completed = completed;
    this.created_at = createdAt;
    this.updated_at = updatedAt;
  }

  toJSON() {
    return {
      id: this.id,
      title: this.title,
      description: this.description,
      completed: this.completed,
      created_at: this.created_at,
      updated_at: this.updated_at,
    };
  }
}
