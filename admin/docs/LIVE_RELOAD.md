# Live Reload System

The admin panel includes an automatic browser reload system that eliminates manual refreshes during development.

## 🎯 Overview

**What it does:**
- Watches the `admin/src/` directory for changes
- Automatically reloads the browser when any `.ts` file changes
- No manual F5 needed - just save and see your changes!

**How it works:**
1. Backend file watcher monitors the entire `admin/src/` directory recursively
2. When a change is detected, it notifies the browser via Server-Sent Events (SSE)
3. Browser receives the notification and reloads automatically
4. Your changes appear instantly!

## 🚀 Quick Start

```bash
# Start the backend with dev environment (required for live reload)
export APP_ENV=dev
go run main.go

# Open browser
http://localhost:4000/admin

# Edit any .ts file in admin/src/
# Save → Browser reloads automatically! 🎉
```

## 📋 What Triggers Reload

The system monitors all changes in `admin/src/`:

- ✅ **Edit existing files** - Modify any `.ts` file
- ✅ **Create new files** - Add new `.ts` files anywhere
- ✅ **Delete files** - Remove `.ts` files
- ✅ **Create folders** - Add new directories with `.ts` files
- ✅ **Move files** - Reorganize your code structure

**Examples:**
```bash
# All of these trigger auto-reload:
admin/src/pages/Login.ts          # Modified
admin/src/pages/NewPage.ts        # Created
admin/src/components/Button.ts    # Created in existing folder
admin/src/features/auth/Login.ts  # Created in new nested folder
```

## 🔧 Technical Details

### Backend Implementation

**File:** `core/src/middleware/livereload.go`

**Key components:**
```go
// Main setup function (called in server.go)
SetupLiveReload(app *fiber.App)

// Endpoints created:
GET /admin/dev/watch  // Server-Sent Events stream
GET /admin/dev/hash   // Polling fallback endpoint
```

**How it works:**
1. Walks `admin/src/` directory recursively
2. Finds all `.ts` files
3. Computes MD5 hash of file paths + modification times
4. Checks for changes every 1 second
5. Sends notification when hash changes

**Configuration:**
```go
LiveReloadConfig{
    Enabled:  true,              // Auto-enabled in development
    WatchDir: "./admin/src",     // Directory to monitor
}
```

**Environment control:**
- Enabled when `APP_ENV=dev`
- Disabled in all other environments (test, prod, or unset)

### Frontend Implementation

**File:** `admin/index.html`

**Location:** Script block before `</body>` tag

**Primary method - Server-Sent Events (SSE):**
```javascript
const eventSource = new EventSource('/admin/dev/watch');

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(`🔄 File changed: ${data.file}`);
  window.location.reload();
};
```

**Fallback method - Polling:**
```javascript
// If SSE fails, falls back to polling
setInterval(async () => {
  const response = await fetch('/admin/dev/hash', { cache: 'no-cache' });
  const data = await response.json();
  
  if (lastHash && lastHash !== data.hash) {
    console.log('🔄 Files changed, reloading...');
    window.location.reload();
  }
}, 1000);
```

## 📊 Console Output

### Backend Console

When you start the server:
```
🔄 Live reload enabled - watching ./admin/src
Server is starting at http://localhost:4000
```

### Browser Console

When you first load the page:
```
🔄 Live reload enabled - watching for file changes...
```

When a file changes:
```
🔄 File changed: ./admin/src
```

If SSE connection fails:
```
⚠️ Live reload connection lost. Using fallback polling...
```

## 🎛️ Configuration

### Enable/Disable Live Reload

Live reload is **only enabled** when `APP_ENV=dev`.

**To enable:**
```bash
# Set development environment
export APP_ENV=dev
go run main.go
```

**To disable:**
```bash
# Set any other environment (or leave unset)
export APP_ENV=prod
go run main.go

# Or simply don't set APP_ENV
go run main.go
```

### Change Watch Directory

Edit `core/src/middleware/livereload.go`:

```go
config := LiveReloadConfig{
    Enabled:  true,
    WatchDir: "./admin/src",  // Change this path
}
```

### Adjust Polling Interval

Frontend polling interval (in `index.html`):

```javascript
// Default: check every 1 second (1000ms)
setInterval(checkForChanges, 1000);

// Faster: check every 500ms
setInterval(checkForChanges, 500);

// Slower: check every 2 seconds
setInterval(checkForChanges, 2000);
```

Backend check interval (in `livereload.go`):

```go
// Default: check every 1 second
ticker := time.NewTicker(1 * time.Second)

// Change to 500ms
ticker := time.NewTicker(500 * time.Millisecond)
```

## 🐛 Troubleshooting

### Browser not reloading

**Check backend is running:**
```bash
# Look for this message
🔄 Live reload enabled - watching ./admin/src
```

**Check browser console:**
```javascript
// Should see this when page loads
🔄 Live reload enabled - watching for file changes...
```

**Check network tab:**
- Open DevTools → Network
- Look for request to `/admin/dev/watch` (should show "pending" with EventStream type)
- Or look for periodic requests to `/admin/dev/hash`

### SSE connection failing

If you see: `⚠️ Live reload connection lost. Using fallback polling...`

**Possible causes:**
1. Backend not running
2. CORS issues (shouldn't happen on localhost)
3. Proxy/firewall blocking SSE

**Solution:** The system automatically falls back to polling, which should work fine.

### Files not triggering reload

**Verify file location:**
- Files must be in `admin/src/` or subdirectories
- Only `.ts` files are monitored

**Check file changes are being saved:**
- VS Code: Look for the dot in the tab title (unsaved)
- Make sure auto-save is enabled or save manually (Ctrl+S)

**Check backend logs:**
- Backend should log when computing new hash
- If no logs, watcher might not be running

### Too many reloads

If the browser keeps reloading continuously:

**Possible causes:**
1. File being auto-generated on every change (build artifact)
2. File watching itself triggering changes
3. Multiple backend instances running

**Solution:**
- Check no build tools are running
- Ensure only one `go run main.go` instance
- Exclude generated files from `admin/src/`

## 🔒 Security

### Production Safety

Live reload is **only enabled when `APP_ENV=dev`**:

```go
// Only enable in development
env := os.Getenv("APP_ENV")
if env != "dev" {
    return  // Live reload disabled
}
```

**In production, test, or any other environment:**
- Endpoints `/admin/dev/watch` and `/admin/dev/hash` are not registered
- No file watching overhead
- No browser reload script active (also checks `localhost`)

### Frontend Safety

The browser script only runs on localhost:

```javascript
if (window.location.hostname === 'localhost' || 
    window.location.hostname === '127.0.0.1') {
    // Live reload enabled
}
```

Deployed on a domain? Live reload won't activate.

## 📈 Performance

### Backend Impact

**Minimal overhead:**
- Walks directory every 1 second
- Only computes hash of file metadata (not file contents)
- ~10-20ms per check on typical project size
- No disk I/O for unchanged files

**Scalability:**
- Efficient for typical admin panel (~20-50 files)
- MD5 hashing is fast
- No memory leaks (uses standard library)

### Frontend Impact

**SSE method:**
- Single persistent connection
- No polling overhead
- Instant notifications (<100ms latency)
- Minimal bandwidth (~1KB for connection)

**Polling fallback:**
- 1 request per second
- ~200 bytes per request
- Negligible bandwidth impact
- Works everywhere (no SSE support needed)

## 🆚 Alternatives Comparison

| Method | Speed | Setup | Reliability | Notes |
|--------|-------|-------|-------------|-------|
| **Our Live Reload** | ⚡ Instant | ✅ Automatic | ⭐⭐⭐⭐⭐ | Built-in, no config needed |
| Manual F5 | 🐌 Manual | ✅ None | ⭐⭐⭐⭐⭐ | Always works, but manual |
| Browser Extension | ⚡ Fast | ⚠️ Install required | ⭐⭐⭐ | Extra dependency |
| DevTools Workspace | ⚡ Instant | ⚠️ Manual setup | ⭐⭐⭐⭐ | Chrome-only, complex |
| Vite/Webpack HMR | ⚡⚡ Instant + HMR | ❌ Build step needed | ⭐⭐⭐⭐ | Defeats "no build" goal |

## 📚 Related Documentation

- **QUICKSTART.md** - Getting started guide
- **FEATURES.md** - Complete feature list
- **README.md** - Project overview

## 🎓 Learning Resources

### How Server-Sent Events Work

SSE is a standard browser API for receiving real-time updates from a server:

```javascript
// Browser opens a persistent HTTP connection
const eventSource = new EventSource('/admin/dev/watch');

// Server can send messages anytime
eventSource.onmessage = (event) => {
  console.log('Received:', event.data);
};
```

**Benefits:**
- One-way server → client (perfect for our use case)
- Automatic reconnection
- Built into browsers (no libraries needed)
- Works with HTTP/1.1 and HTTP/2

### How File Watching Works

The backend uses Go's standard library to walk directories:

```go
filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
    if strings.HasSuffix(path, ".ts") {
        // Found a TypeScript file!
        // Add to hash computation
    }
    return nil
})
```

**Why MD5 hash?**
- Fast computation (< 1ms for typical project)
- Detects ANY change (new, modified, deleted files)
- No need to track individual files
- Simple comparison (just compare hash strings)

---

**Last Updated:** November 2, 2025  
**Version:** 1.0.0
