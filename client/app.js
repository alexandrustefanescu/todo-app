// ============================================
// TODO APP - VANILLA JAVASCRIPT
// High Performance Frontend
// ============================================

// API Configuration
const API_BASE_URL = 'http://localhost:8080/api';

// State Management
let allTodos = [];
let currentFilter = 'all';
let editingTodoId = null;

// DOM Elements
const todoForm = document.getElementById('todoForm');
const todosList = document.getElementById('todosList');
const loadingMessage = document.getElementById('loadingMessage');
const emptyMessage = document.getElementById('emptyMessage');
const formError = document.getElementById('formError');
const editModal = document.getElementById('editModal');
const editForm = document.getElementById('editForm');
const editError = document.getElementById('editError');
const filterButtons = document.querySelectorAll('.btn-filter');

// ============================================
// API FUNCTIONS
// ============================================

/**
 * Fetch all todos from the API
 */
async function fetchTodos() {
    try {
        loadingMessage.style.display = 'block';
        emptyMessage.style.display = 'none';
        todosList.innerHTML = '';

        const response = await fetch(`${API_BASE_URL}/todos`);

        if (!response.ok) {
            throw new Error(`Failed to fetch todos: ${response.statusText}`);
        }

        allTodos = await response.json();
        renderTodos();
    } catch (error) {
        console.error('Error fetching todos:', error);
        showErrorMessage('Failed to load todos. Make sure the API is running.', loadingMessage);
    }
}

/**
 * Create a new todo
 */
async function createTodo(title, description) {
    try {
        const response = await fetch(`${API_BASE_URL}/todos`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title.trim(),
                description: description.trim() || null,
            }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to create todo');
        }

        const newTodo = await response.json();
        allTodos.unshift(newTodo);
        renderTodos();
        todoForm.reset();
        clearErrorMessage(formError);
        showSuccessMessage('Todo created successfully!');
    } catch (error) {
        console.error('Error creating todo:', error);
        showErrorMessage(error.message, formError);
    }
}

/**
 * Update a todo
 */
async function updateTodo(id, updates) {
    try {
        const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updates),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to update todo');
        }

        const updatedTodo = await response.json();
        allTodos = allTodos.map(todo => (todo.id === id ? updatedTodo : todo));
        renderTodos();
        closeEditModal();
        showSuccessMessage('Todo updated successfully!');
    } catch (error) {
        console.error('Error updating todo:', error);
        showErrorMessage(error.message, editError);
    }
}

/**
 * Delete a todo
 */
async function deleteTodo(id) {
    try {
        const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to delete todo');
        }

        allTodos = allTodos.filter(todo => todo.id !== id);
        renderTodos();
        showSuccessMessage('Todo deleted successfully!');
    } catch (error) {
        console.error('Error deleting todo:', error);
        showErrorMessage(error.message, loadingMessage);
    }
}

// ============================================
// UI RENDERING FUNCTIONS
// ============================================

/**
 * Filter todos based on current filter state
 */
function getFilteredTodos() {
    if (currentFilter === 'completed') {
        return allTodos.filter(todo => todo.completed);
    }
    if (currentFilter === 'pending') {
        return allTodos.filter(todo => !todo.completed);
    }
    return allTodos;
}

/**
 * Render todos list
 */
function renderTodos() {
    loadingMessage.style.display = 'none';

    const filteredTodos = getFilteredTodos();

    if (filteredTodos.length === 0) {
        emptyMessage.style.display = 'block';
        todosList.innerHTML = '';
        return;
    }

    emptyMessage.style.display = 'none';
    todosList.innerHTML = filteredTodos
        .map(todo => createTodoElement(todo))
        .join('');

    // Attach event listeners to dynamically created elements
    attachTodoEventListeners();
}

/**
 * Create HTML for a single todo item
 */
function createTodoElement(todo) {
    const createdDate = new Date(todo.created_at).toLocaleDateString();
    const statusClass = todo.completed ? 'completed' : 'pending';
    const statusText = todo.completed ? '✓ Completed' : '⏳ Pending';

    return `
        <div class="todo-item ${statusClass}" data-id="${todo.id}">
            <input
                type="checkbox"
                class="todo-checkbox"
                ${todo.completed ? 'checked' : ''}
                data-id="${todo.id}"
            >
            <div class="todo-content">
                <div class="todo-title">${escapeHtml(todo.title)}</div>
                ${
                    todo.description
                        ? `<div class="todo-description">${escapeHtml(todo.description)}</div>`
                        : ''
                }
                <div class="todo-meta">
                    <span class="todo-status ${statusClass}">${statusText}</span>
                    <span>Created: ${createdDate}</span>
                </div>
            </div>
            <div class="todo-actions">
                <button class="btn btn-edit edit-btn" data-id="${todo.id}">
                    ✏️ Edit
                </button>
                <button class="btn btn-danger delete-btn" data-id="${todo.id}">
                    🗑️ Delete
                </button>
            </div>
        </div>
    `;
}

/**
 * Attach event listeners to todo items
 */
function attachTodoEventListeners() {
    // Checkbox change listeners
    document.querySelectorAll('.todo-checkbox').forEach(checkbox => {
        checkbox.addEventListener('change', (e) => {
            const id = e.target.dataset.id;
            const completed = e.target.checked;
            updateTodo(id, { completed });
        });
    });

    // Edit button listeners
    document.querySelectorAll('.edit-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = e.target.dataset.id;
            openEditModal(id);
        });
    });

    // Delete button listeners
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = e.target.dataset.id;
            if (confirm('Are you sure you want to delete this todo?')) {
                deleteTodo(id);
            }
        });
    });
}

// ============================================
// MODAL FUNCTIONS
// ============================================

/**
 * Open edit modal with todo data
 */
function openEditModal(id) {
    editingTodoId = id;
    const todo = allTodos.find(t => t.id === id);

    if (!todo) return;

    document.getElementById('editTitle').value = todo.title;
    document.getElementById('editDescription').value = todo.description || '';
    document.getElementById('editCompleted').checked = todo.completed;
    clearErrorMessage(editError);

    editModal.classList.add('show');
}

/**
 * Close edit modal
 */
function closeEditModal() {
    editingTodoId = null;
    editModal.classList.remove('show');
    editForm.reset();
}

// ============================================
// EVENT LISTENERS
// ============================================

/**
 * Form submission
 */
todoForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;

    if (!title.trim()) {
        showErrorMessage('Title is required', formError);
        return;
    }

    createTodo(title, description);
});

/**
 * Edit form submission
 */
editForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const title = document.getElementById('editTitle').value;
    const description = document.getElementById('editDescription').value;
    const completed = document.getElementById('editCompleted').checked;

    if (!title.trim()) {
        showErrorMessage('Title is required', editError);
        return;
    }

    updateTodo(editingTodoId, {
        title: title.trim(),
        description: description.trim() || null,
        completed,
    });
});

/**
 * Modal close button
 */
document.querySelector('.close-btn').addEventListener('click', closeEditModal);

/**
 * Modal cancel button
 */
document.querySelector('.cancel-btn').addEventListener('click', closeEditModal);

/**
 * Close modal when clicking outside
 */
editModal.addEventListener('click', (e) => {
    if (e.target === editModal) {
        closeEditModal();
    }
});

/**
 * Filter buttons
 */
filterButtons.forEach(btn => {
    btn.addEventListener('click', (e) => {
        filterButtons.forEach(b => b.classList.remove('active'));
        e.target.classList.add('active');
        currentFilter = e.target.dataset.filter;
        renderTodos();
    });
});

// ============================================
// UTILITY FUNCTIONS
// ============================================

/**
 * Escape HTML to prevent XSS
 */
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;',
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}

/**
 * Show error message
 */
function showErrorMessage(message, element) {
    element.textContent = message;
    element.classList.add('show');
}

/**
 * Clear error message
 */
function clearErrorMessage(element) {
    element.textContent = '';
    element.classList.remove('show');
}

/**
 * Show success message
 */
function showSuccessMessage(message) {
    const successMsg = document.createElement('div');
    successMsg.className = 'success-message';
    successMsg.textContent = message;
    document.body.insertBefore(successMsg, document.body.firstChild);

    setTimeout(() => {
        successMsg.remove();
    }, 3000);
}

// ============================================
// INITIALIZATION
// ============================================

/**
 * Initialize the app
 */
function initializeApp() {
    console.log('🚀 TODO App initialized');
    console.log('📡 API Base URL:', API_BASE_URL);
    fetchTodos();
}

// Start the app
initializeApp();

// Refresh todos every 30 seconds (optional)
// setInterval(fetchTodos, 30000);
