# Running the TODO App Frontend with Bun

**Bun** is the fastest and simplest way to run the frontend. It's a modern JavaScript runtime that's incredibly fast.

## 🚀 Quick Start with Bun

### 1. Install Bun (One-time only)

```bash
curl -fsSL https://bun.sh/install | bash
```

Verify installation:
```bash
bun --version
```

### 2. Start the Frontend Server

```bash
cd /Users/alexandrustefanescu/Desktop/todo-app/client
bun --serve index.html
```

**That's it!** The server starts on `http://localhost:3000`

### 3. Open in Browser

Go to: `http://localhost:3000`

## ✨ Why Bun?

- ⚡ **Lightning Fast** - Significantly faster than Node.js
- 🔧 **Zero Configuration** - No setup needed
- 📦 **Built-in Server** - No additional packages required
- 🎯 **Modern JavaScript** - Native support for ES modules
- 🚀 **Instant Startup** - Starts in milliseconds

## 🔄 Development Workflow

1. **Start Backend (Terminal 1)**
   ```bash
   cd api
   docker-compose up
   cargo run --release
   ```

2. **Start Frontend (Terminal 2)**
   ```bash
   cd client
   bun --serve index.html
   ```

3. **Open Browser (Terminal 3)**
   ```bash
   open http://localhost:3000
   ```

4. **Edit Files** - Auto-reloads in browser (use Live Server extension for auto-reload, or manually refresh)

## 💡 Tips

### Auto-reload on File Changes

While `bun --serve` doesn't have built-in watch/reload, you can:

**Option A: Use VS Code Live Server**
- Install "Live Server" extension
- Right-click `index.html` → "Open with Live Server"

**Option B: Use Bun with a Watch Script**
Create a `dev.sh` file:
```bash
#!/bin/bash
cd client
bun --serve index.html
```

Then run:
```bash
bash dev.sh
```

### Stop the Server

Press `Ctrl+C` in the terminal

## 🌐 Accessing the App

- **Local**: `http://localhost:3000`
- **API**: `http://localhost:8080/api`

## 📝 Common Commands

```bash
# Start server on port 3000 (default)
bun --serve index.html

# Start server on different port
bun --serve index.html --port 5000

# Check Bun version
bun --version

# Update Bun
bun upgrade
```

## 🐛 Troubleshooting

**Port 3000 already in use?**
```bash
bun --serve index.html --port 5000
```

**CORS error?**
- Make sure backend is running on port 8080
- Check `API_BASE_URL` in `app.js`

**API not responding?**
```bash
# Test API connection
curl http://localhost:8080/api/todos
```

## 🔗 More Info

- [Bun Documentation](https://bun.sh)
- [Frontend README](README.md)
- [Backend Setup](../api/README.md)

---

That's all you need to run the frontend! 🚀
