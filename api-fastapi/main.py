"""
FastAPI Todo Application
A modern, high-performance Todo REST API built with FastAPI and PostgreSQL
"""

import logging
import os
from contextlib import asynccontextmanager
from typing import Optional

import uvicorn
from fastapi import FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, async_sessionmaker

from models import Base, Todo
from schemas import TodoCreate, TodoUpdate, TodoResponse

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Database configuration
DATABASE_URL = os.getenv(
    'DATABASE_URL',
    'postgresql+asyncpg://postgres:password@localhost:5432/todo_db'
)

# SQLAlchemy setup
engine = create_async_engine(
    DATABASE_URL,
    echo=False,
    future=True,
    pool_pre_ping=True,
    pool_size=10,
    max_overflow=20,
)

AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
    autocommit=False,
    autoflush=False,
)


async def get_db():
    """Dependency: Get database session"""
    async with AsyncSessionLocal() as session:
        try:
            yield session
        finally:
            await session.close()


async def init_db():
    """Initialize database tables"""
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    logger.info("Database initialized successfully")


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Lifespan context manager for startup and shutdown events
    """
    # Startup
    logger.info("Starting FastAPI Todo Application")
    await init_db()
    yield
    # Shutdown
    logger.info("Shutting down FastAPI Todo Application")
    await engine.dispose()


# Create FastAPI app
app = FastAPI(
    title="Todo API",
    description="A modern Todo REST API built with FastAPI",
    version="1.0.0",
    lifespan=lifespan,
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ==================== Health Check ====================
@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "ok", "service": "Todo API"}


# ==================== Todo Endpoints ====================

@app.get("/api/todos", response_model=list[TodoResponse])
async def list_todos(session: AsyncSession = None):
    """
    List all todos ordered by creation date (newest first)

    Returns:
        List of all todos
    """
    if session is None:
        async with AsyncSessionLocal() as db:
            result = await db.execute(
                select(Todo).order_by(Todo.created_at.desc())
            )
            todos = result.scalars().all()
            return [TodoResponse.from_orm(todo) for todo in todos]

    result = await session.execute(
        select(Todo).order_by(Todo.created_at.desc())
    )
    todos = result.scalars().all()
    return [TodoResponse.from_orm(todo) for todo in todos]


@app.post("/api/todos", response_model=TodoResponse, status_code=status.HTTP_201_CREATED)
async def create_todo(
    todo_data: TodoCreate,
    session: AsyncSession = None
):
    """
    Create a new todo

    Args:
        todo_data: Todo creation data with title and optional description

    Returns:
        Created todo with generated UUID and timestamps
    """
    if not todo_data.title or not todo_data.title.strip():
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Title is required and cannot be empty"
        )

    if session is None:
        async with AsyncSessionLocal() as db:
            new_todo = Todo(
                title=todo_data.title.strip(),
                description=todo_data.description.strip() if todo_data.description else None,
            )
            db.add(new_todo)
            await db.commit()
            await db.refresh(new_todo)
            return TodoResponse.from_orm(new_todo)

    new_todo = Todo(
        title=todo_data.title.strip(),
        description=todo_data.description.strip() if todo_data.description else None,
    )
    session.add(new_todo)
    await session.commit()
    await session.refresh(new_todo)
    return TodoResponse.from_orm(new_todo)


@app.get("/api/todos/{todo_id}", response_model=TodoResponse)
async def get_todo(
    todo_id: str,
    session: AsyncSession = None
):
    """
    Get a specific todo by ID

    Args:
        todo_id: UUID of the todo to retrieve

    Returns:
        The requested todo or 404 if not found
    """
    if session is None:
        async with AsyncSessionLocal() as db:
            result = await db.execute(select(Todo).where(Todo.id == todo_id))
            todo = result.scalar_one_or_none()
    else:
        result = await session.execute(select(Todo).where(Todo.id == todo_id))
        todo = result.scalar_one_or_none()

    if not todo:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Todo with id {todo_id} not found"
        )

    return TodoResponse.from_orm(todo)


@app.put("/api/todos/{todo_id}", response_model=TodoResponse)
async def update_todo(
    todo_id: str,
    todo_data: TodoUpdate,
    session: AsyncSession = None
):
    """
    Update a todo (partial update supported)

    Args:
        todo_id: UUID of the todo to update
        todo_data: Update data (title, description, completed - all optional)

    Returns:
        Updated todo or 404 if not found
    """
    if session is None:
        async with AsyncSessionLocal() as db:
            result = await db.execute(select(Todo).where(Todo.id == todo_id))
            todo = result.scalar_one_or_none()

            if not todo:
                raise HTTPException(
                    status_code=status.HTTP_404_NOT_FOUND,
                    detail=f"Todo with id {todo_id} not found"
                )

            # Update only provided fields
            if todo_data.title is not None:
                if not todo_data.title.strip():
                    raise HTTPException(
                        status_code=status.HTTP_400_BAD_REQUEST,
                        detail="Title cannot be empty"
                    )
                todo.title = todo_data.title.strip()

            if todo_data.description is not None:
                todo.description = todo_data.description.strip() if todo_data.description else None

            if todo_data.completed is not None:
                todo.completed = todo_data.completed

            await db.commit()
            await db.refresh(todo)
            return TodoResponse.from_orm(todo)

    result = await session.execute(select(Todo).where(Todo.id == todo_id))
    todo = result.scalar_one_or_none()

    if not todo:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Todo with id {todo_id} not found"
        )

    # Update only provided fields
    if todo_data.title is not None:
        if not todo_data.title.strip():
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Title cannot be empty"
            )
        todo.title = todo_data.title.strip()

    if todo_data.description is not None:
        todo.description = todo_data.description.strip() if todo_data.description else None

    if todo_data.completed is not None:
        todo.completed = todo_data.completed

    await session.commit()
    await session.refresh(todo)
    return TodoResponse.from_orm(todo)


@app.delete("/api/todos/{todo_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_todo(
    todo_id: str,
    session: AsyncSession = None
):
    """
    Delete a todo by ID

    Args:
        todo_id: UUID of the todo to delete

    Returns:
        204 No Content on success, 404 if not found
    """
    if session is None:
        async with AsyncSessionLocal() as db:
            result = await db.execute(select(Todo).where(Todo.id == todo_id))
            todo = result.scalar_one_or_none()

            if not todo:
                raise HTTPException(
                    status_code=status.HTTP_404_NOT_FOUND,
                    detail=f"Todo with id {todo_id} not found"
                )

            await db.delete(todo)
            await db.commit()
            return None

    result = await session.execute(select(Todo).where(Todo.id == todo_id))
    todo = result.scalar_one_or_none()

    if not todo:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Todo with id {todo_id} not found"
        )

    await session.delete(todo)
    await session.commit()
    return None


# ==================== Root Endpoint ====================
@app.get("/")
async def root():
    """API root endpoint"""
    return {
        "message": "Welcome to Todo API",
        "version": "1.0.0",
        "docs": "/docs",
        "redoc": "/redoc"
    }


if __name__ == "__main__":
    port = int(os.getenv("PORT", 8080))
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=os.getenv("ENV", "production") == "development",
        log_level="info",
    )
