# TODO App - Frontend Client

A high-performance vanilla JavaScript frontend for the TODO App REST API.

## 🎯 Features

✅ **Pure Vanilla JavaScript** - Raw speed and performance, no frameworks
✅ **Real-time Todo Management** - Create, read, update, delete todos
✅ **Responsive Design** - Works beautifully on desktop and mobile
✅ **Dark Theme** - Modern, eye-friendly interface
✅ **Todo Filtering** - View all, pending, or completed todos
✅ **XSS Protection** - HTML escaping for security
✅ **Error Handling** - User-friendly error messages
✅ **Success Notifications** - Feedback for user actions
✅ **Zero Dependencies** - No npm packages required

## 📁 Project Structure

```
client/
├── index.html      # Main HTML structure
├── styles.css      # Complete styling (dark theme)
├── app.js          # Main application logic
├── .htaccess       # Server routing config
└── README.md       # This file
```

## 🚀 Quick Start

### Prerequisites

- Backend API running on `http://localhost:8080/api`
- Bun (fastest and simplest option)

### Running the Frontend

#### Option 1: Bun (Recommended - Fastest)

Install Bun if you don't have it:
```bash
curl -fsSL https://bun.sh/install | bash
```

Run the frontend:
```bash
cd /Users/alexandrustefanescu/Desktop/todo-app/client
bun --serve index.html
```

Then open: `http://localhost:3000`

**That's it!** Bun starts automatically on port 3000 with zero configuration.

#### Option 2: Live Server (VS Code Extension)

1. Install "Live Server" extension in VS Code
2. Right-click `index.html` → "Open with Live Server"

#### Option 3: Node.js http-server

```bash
cd /Users/alexandrustefanescu/Desktop/todo-app/client
npx http-server -p 3000
```

Then open: `http://localhost:3000`

## 📡 API Integration

The frontend communicates with the backend API:

```
Base URL: http://localhost:8080/api
```

### Endpoints Used

| Method | Endpoint | Action |
|--------|----------|--------|
| GET | `/todos` | Fetch all todos |
| POST | `/todos` | Create new todo |
| GET | `/todos/{id}` | Get single todo |
| PUT | `/todos/{id}` | Update todo |
| DELETE | `/todos/{id}` | Delete todo |

## 🎮 Usage

### Creating a Todo

1. Enter a title (required)
2. Optionally add a description
3. Click "Add Todo"

### Updating a Todo

1. Click the "✏️ Edit" button on any todo
2. Modify the title, description, or status
3. Click "Save Changes"

### Completing a Todo

- Click the checkbox next to a todo to mark it as completed
- The todo will be struck through

### Deleting a Todo

1. Click the "🗑️ Delete" button
2. Confirm the deletion

### Filtering Todos

Use the filter buttons at the top:
- **All** - Show all todos
- **Pending** - Show incomplete todos
- **Completed** - Show completed todos

## 🛡️ Security Features

- **HTML Escaping** - Prevents XSS attacks
- **Input Validation** - Validates required fields
- **CORS Handling** - Backend should allow frontend origin
- **Error Messages** - Sanitized error display

## 📦 Performance Optimizations

- No framework overhead
- Direct DOM manipulation
- Minimal JavaScript bundle
- Efficient event delegation
- CSS animations with GPU acceleration
- Responsive images and design

## 🎨 Design Features

### Dark Theme

- Easy on the eyes
- Reduces eye strain
- Modern aesthetic
- Good contrast ratios (WCAG AA)

### Responsive Layout

- Desktop: 2-column layout (form + list)
- Tablet: Adjusted spacing
- Mobile: 1-column layout

### Animations

- Smooth transitions
- Modal fade-in/slide-up
- Success message animations
- Respects `prefers-reduced-motion` preference

## 🧪 Testing

### Test API Connection

Open browser console and run:

```javascript
fetch('http://localhost:8080/api/todos')
  .then(r => r.json())
  .then(d => console.log(d))
  .catch(e => console.error('API Error:', e))
```

### Common Issues

**CORS Error?**
- Make sure backend is running
- Backend should allow requests from `localhost:3000`

**API Connection Failed?**
- Check backend is running on port 8080
- Check `API_BASE_URL` in `app.js`

**Styles Not Loading?**
- Ensure `styles.css` is in the same directory
- Check browser console for 404 errors

## 📱 Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari, Chrome Mobile)

## 🔧 Configuration

Edit `API_BASE_URL` in `app.js` to change the backend URL:

```javascript
const API_BASE_URL = 'http://localhost:8080/api';
```

## 📚 File Breakdown

### index.html
- Semantic HTML5 structure
- Form for creating todos
- Modal for editing todos
- Todo list container
- Filter buttons

### styles.css
- CSS Custom Properties (variables)
- Dark theme styling
- Responsive grid layout
- Button and form styles
- Modal styling
- Animations
- Mobile-first approach

### app.js
- Fetch all todos on page load
- Create, update, delete todo functions
- Render todos to DOM
- Event listeners for user interactions
- Modal management
- Error/success message handling
- HTML escaping for security

## 🚀 Deployment

### Static Hosting (Netlify, Vercel, GitHub Pages)

1. Copy `index.html`, `styles.css`, `app.js` to your hosting
2. Update `API_BASE_URL` to your production API URL
3. Deploy!

### Using with Docker

Build a simple Docker image:

```dockerfile
FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 📖 Code Quality

- Clean, commented code
- Consistent naming conventions
- Modular function structure
- Error handling with try/catch
- Console logging for debugging
- No external dependencies

## 🤝 Contributing

Feel free to extend this frontend:

- Add filtering by date
- Add sorting options
- Add search functionality
- Add tags/categories
- Add recurring todos
- Add todo priority levels

## 📄 License

MIT License - Use freely in your projects

## 🔗 Related Documentation

- [Backend API](../api/README.md)
- [Project Overview](../README.md)
- [Quick Start Guide](../api/QUICKSTART.md)

---

**Built with ❤️ using Vanilla JavaScript**
