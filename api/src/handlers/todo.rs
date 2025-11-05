use actix_web::{web, HttpResponse};
use sqlx::PgPool;
use uuid::Uuid;
use chrono::Utc;

use crate::models::{CreateTodoRequest, UpdateTodoRequest, TodoResponse, Todo};
use crate::error::ApiError;

/// List all todos
pub async fn list_todos(pool: web::Data<PgPool>) -> Result<HttpResponse, ApiError> {
    let todos = sqlx::query_as::<_, Todo>(
        "SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC"
    )
    .fetch_all(pool.get_ref())
    .await?;

    let response: Vec<TodoResponse> = todos.into_iter().map(|t| t.into()).collect();
    Ok(HttpResponse::Ok().json(response))
}

/// Get a single todo by ID
pub async fn get_todo(
    pool: web::Data<PgPool>,
    id: web::Path<Uuid>,
) -> Result<HttpResponse, ApiError> {
    let todo = sqlx::query_as::<_, Todo>(
        "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1"
    )
    .bind(id.into_inner())
    .fetch_one(pool.get_ref())
    .await?;

    Ok(HttpResponse::Ok().json(TodoResponse::from(todo)))
}

/// Create a new todo
pub async fn create_todo(
    pool: web::Data<PgPool>,
    req: web::Json<CreateTodoRequest>,
) -> Result<HttpResponse, ApiError> {
    if req.title.trim().is_empty() {
        return Err(ApiError::BadRequest("Title cannot be empty".to_string()));
    }

    let id = Uuid::new_v4();
    let now = Utc::now();

    let todo = sqlx::query_as::<_, Todo>(
        "INSERT INTO todos (id, title, description, completed, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id, title, description, completed, created_at, updated_at"
    )
    .bind(id)
    .bind(&req.title)
    .bind(&req.description)
    .bind(false)
    .bind(now)
    .bind(now)
    .fetch_one(pool.get_ref())
    .await?;

    Ok(HttpResponse::Created().json(TodoResponse::from(todo)))
}

/// Update a todo
pub async fn update_todo(
    pool: web::Data<PgPool>,
    id: web::Path<Uuid>,
    req: web::Json<UpdateTodoRequest>,
) -> Result<HttpResponse, ApiError> {
    let id = id.into_inner();
    let now = Utc::now();

    // First, check if the todo exists
    let existing = sqlx::query_as::<_, Todo>(
        "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1"
    )
    .bind(id)
    .fetch_optional(pool.get_ref())
    .await?;

    if existing.is_none() {
        return Err(ApiError::NotFound(format!("Todo with id {} not found", id)));
    }

    let existing = existing.unwrap();

    // Update fields, keeping existing values if not provided
    let title = req.title.as_ref().unwrap_or(&existing.title).clone();
    let description = req.description.as_ref().or(existing.description.as_ref()).cloned();
    let completed = req.completed.unwrap_or(existing.completed);

    let todo = sqlx::query_as::<_, Todo>(
        "UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = $4
         WHERE id = $5
         RETURNING id, title, description, completed, created_at, updated_at"
    )
    .bind(title)
    .bind(description)
    .bind(completed)
    .bind(now)
    .bind(id)
    .fetch_one(pool.get_ref())
    .await?;

    Ok(HttpResponse::Ok().json(TodoResponse::from(todo)))
}

/// Delete a todo
pub async fn delete_todo(
    pool: web::Data<PgPool>,
    id: web::Path<Uuid>,
) -> Result<HttpResponse, ApiError> {
    let id = id.into_inner();

    let result = sqlx::query("DELETE FROM todos WHERE id = $1")
        .bind(id)
        .execute(pool.get_ref())
        .await?;

    if result.rows_affected() == 0 {
        return Err(ApiError::NotFound(format!("Todo with id {} not found", id)));
    }

    Ok(HttpResponse::NoContent().finish())
}
