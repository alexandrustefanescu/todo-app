"""
Pydantic schemas for request/response validation
"""

from datetime import datetime
from typing import Optional
from uuid import UUID

from pydantic import BaseModel, Field, field_validator


class TodoCreate(BaseModel):
    """Schema for creating a new todo"""

    title: str = Field(..., min_length=1, max_length=255, description="Todo title")
    description: Optional[str] = Field(
        None,
        max_length=5000,
        description="Optional todo description"
    )

    @field_validator('title')
    @classmethod
    def title_not_empty(cls, v):
        if not v or not v.strip():
            raise ValueError('Title cannot be empty or whitespace')
        return v.strip()

    @field_validator('description')
    @classmethod
    def description_cleanup(cls, v):
        if v is not None:
            v = v.strip()
            if v == "":
                return None
        return v

    class Config:
        json_schema_extra = {
            "example": {
                "title": "Learn FastAPI",
                "description": "Study FastAPI documentation and best practices"
            }
        }


class TodoUpdate(BaseModel):
    """Schema for updating a todo (all fields optional)"""

    title: Optional[str] = Field(
        None,
        min_length=1,
        max_length=255,
        description="Updated todo title"
    )
    description: Optional[str] = Field(
        None,
        max_length=5000,
        description="Updated todo description"
    )
    completed: Optional[bool] = Field(
        None,
        description="Completion status"
    )

    @field_validator('title')
    @classmethod
    def title_not_empty(cls, v):
        if v is not None and (not v or not v.strip()):
            raise ValueError('Title cannot be empty or whitespace')
        return v.strip() if v else v

    @field_validator('description')
    @classmethod
    def description_cleanup(cls, v):
        if v is not None:
            v = v.strip()
            if v == "":
                return None
        return v

    class Config:
        json_schema_extra = {
            "example": {
                "title": "Learn FastAPI Advanced",
                "description": "Study advanced FastAPI patterns",
                "completed": True
            }
        }


class TodoResponse(BaseModel):
    """Schema for todo responses"""

    id: UUID
    title: str
    description: Optional[str] = None
    completed: bool
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True
        json_schema_extra = {
            "example": {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "title": "Learn FastAPI",
                "description": "Study FastAPI documentation",
                "completed": False,
                "created_at": "2025-11-05T10:30:00Z",
                "updated_at": "2025-11-05T10:30:00Z"
            }
        }
