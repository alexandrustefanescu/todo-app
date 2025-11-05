"""
SQLAlchemy database models for Todo application
"""

from datetime import datetime
from typing import Optional

import uuid
from sqlalchemy import Column, String, Boolean, DateTime, Text
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class Todo(Base):
    """
    Todo database model

    Attributes:
        id: UUID primary key (auto-generated)
        title: Todo title (required)
        description: Todo description (optional)
        completed: Completion status (default: False)
        created_at: Timestamp when todo was created
        updated_at: Timestamp when todo was last updated
    """

    __tablename__ = "todos"

    id = Column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        nullable=False,
    )
    title = Column(String(255), nullable=False, index=True)
    description = Column(Text, nullable=True)
    completed = Column(Boolean, default=False, nullable=False, index=True)
    created_at = Column(
        DateTime(timezone=True),
        default=datetime.utcnow,
        nullable=False,
        index=True,
    )
    updated_at = Column(
        DateTime(timezone=True),
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False,
    )

    def __repr__(self) -> str:
        return f"<Todo(id={self.id}, title={self.title}, completed={self.completed})>"
