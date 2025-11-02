# Admin Panel - Quick Start Guide

## 🚀 First Time Setup

### Step 1: Start Backend
```bash
# In your mynute-go root directory
go run main.go
```

### Step 2: Access Admin Panel
Open your browser and navigate to:
```
http://localhost:3000/admin
```

### Step 3: Create First Admin (One-time)

If no admin exists, you'll see the **First Admin Registration** form:

1. Fill in the form:
   - First Name
   - Last Name
   - Email
   - Password (min 8 characters)
   - Confirm Password

2. Click "Create Admin Account"

3. Check your email for verification link

4. Click verification link to activate your account

5. Return to login page

### Step 4: Login

Enter your credentials and access the dashboard!

## 📁 Project Structure

```
/admin
  ├── index.html              # Entry point with import maps
  ├── package.json            # Dependencies (Playwright, TypeScript)
  ├── playwright.config.ts    # E2E test configuration
  ├── tsconfig.json          # TypeScript config (no compilation)
  ├── /docs                  # Documentation
  │   ├── FEATURES.md        # Feature summary
  │   ├── QUICKSTART.md      # This file
  │   ├── TEST_COVERAGE.md   # Test coverage report
  │   └── FIRST_ADMIN_REGISTRATION.md
  ├── /src                   # Application source code
  │   ├── main.ts            # Bootstrap
  │   ├── App.ts             # Main component with routing
  │   ├── types.ts           # Type definitions
  │   ├── global.d.ts        # CDN module declarations
  │   ├── /components        # Reusable UI components
  │   │   ├── Layout.ts
  │   │   ├── Header.ts
  │   │   └── Sidebar.ts
  │   ├── /pages             # Page components
  │   │   ├── Login.ts       # Login + First admin registration
  │   │   ├── Dashboard.ts   # Main dashboard
  │   │   ├── Companies.ts   # Companies list
  │   │   ├── CompanyDetail.ts  # Company details
  │   │   ├── Clients.ts     # Clients list
  │   │   └── Users.ts       # Admin users
  │   ├── /stores            # Preact Signals stores
  │   │   ├── authStore.ts   # Authentication
  │   │   ├── adminStore.ts  # Admin users CRUD
  │   │   ├── companyStore.ts  # Companies CRUD
  │   │   └── clientStore.ts   # Clients & appointments
  │   └── /utils             # Utilities
  │       ├── api.ts         # API client
  │       └── constants.ts   # Constants
  └── /tests                 # E2E test suite (86+ tests)
      ├── fixtures.ts        # Test fixtures
      ├── login.spec.ts
      ├── first-admin-registration.spec.ts
      ├── dashboard.spec.ts
      ├── companies.spec.ts
      ├── company-detail.spec.ts
      ├── clients.spec.ts
      ├── users.spec.ts
      ├── navigation.spec.ts
      └── integration.spec.ts
```

## 💻 Development

### Auto-Reload Development (Recommended)

The admin panel includes **automatic browser reload** when you edit files:

**How it works:**
- Backend watches `admin/src/` directory recursively
- Detects changes to ANY `.ts` file (including new files and folders)
- Sends notification to browser via Server-Sent Events (SSE)
- Browser automatically reloads with your changes

**Setup:**
```bash
# 1. Start backend with dev environment (live reload automatic)
export APP_ENV=dev  # Required for live reload
go run main.go

# 2. Open browser
http://localhost:4000/admin

# 3. Edit any .ts file in admin/src/
# Save → Browser reloads automatically! 🎉
```

**Console output:**
- Backend: `🔄 Live reload enabled - watching ./admin/src`
- Browser: `🔄 File changed: ./admin/src` (when files change)

**What triggers reload:**
- ✅ Edit existing `.ts` files
- ✅ Create new `.ts` files
- ✅ Delete `.ts` files
- ✅ Add new folders with `.ts` files
- ✅ Any change in `admin/src/` directory tree

### Local Development

1. **Start Backend:**
   ```bash
   go run main.go
   ```

2. **Open in Browser:**
   ```
   http://localhost:4000/admin
   ```

3. **Edit Code:**
   - Make changes to `.ts` files
   - Browser reloads automatically (or press F5 for manual refresh)
   - No build step needed!

4. **TypeScript Checking:**
   - Use VS Code for inline type checking
   - Red squiggles show errors immediately
   - Hover for type information

### Development Tips

✅ **Browser DevTools** - Full source maps available for debugging  
✅ **Auto-Reload** - Files watched automatically, browser reloads on change  
✅ **Type Safety** - TypeScript catches errors before runtime  
✅ **Signals** - State changes auto-update UI  
✅ **No Build** - Direct `.ts` file execution in browser  

## 🧪 Running Tests

```bash
cd admin

# Install dependencies (first time only)
npm install

# Run all tests (headless mode)
npm test

# Run tests in UI mode (recommended for development)
npm run test:ui

# Run tests with visible browser
npm run test:headed

# Debug a specific test file
npx playwright test tests/login.spec.ts --debug

# View last test report
npm run test:report

# Generate HTML report
npx playwright show-report
```

## 📋 Test Coverage

### Total: 86+ Tests across 8 Test Suites

✅ **Login & First Admin** (17 tests)
- Login form display and validation
- First admin registration flow
- Password validation
- Email verification
- Success/error handling

✅ **Dashboard** (14 tests)
- Stats display (Companies, Clients, Admins)
- Quick action buttons
- System information
- Navigation
- Loading states

✅ **Companies** (13 tests)
- List display with search
- Company cards
- Navigation to details
- Delete functionality
- Empty states

✅ **Company Details** (17 tests)
- All 5 tabs (Overview, Branches, Employees, Services, Subdomains)
- Tab switching
- Data display
- Delete actions
- Error handling

✅ **Clients** (20 tests)
- Table display with search
- Client details modal
- Appointments view
- Delete functionality
- Avatar display

✅ **Admin Users** (existing tests)
- List display
- CRUD operations

✅ **Navigation** (7 tests)
- Sidebar navigation
- Active route highlighting
- Auth persistence
- Responsive layout

✅ **Integration** (24 tests)
- Error handling
- Loading states
- Data consistency
- Accessibility
- Edge cases

## 🎯 Key Features to Try

### 1. Dashboard Overview
- View real-time statistics
- Click stat cards to navigate
- Use quick action buttons

### 2. Companies Management
- Browse companies in grid view
- Search by name or tax ID
- Click "View Details" to see:
  - Company information
  - Branches
  - Employees
  - Services
  - Subdomains

### 3. Clients Management
- View all clients in table
- Search by name, email, or phone
- Click "View Details" to see:
  - Client information
  - Appointment history
  - Status indicators

### 4. Admin Users
- View all admin accounts
- Create new admins
- Manage user roles

## 🏗️ Tech Stack

**Frontend (No Build):**
- **Preact** - 3KB React alternative
- **TypeScript** - Type safety in browser
- **Preact Signals** - Reactive state management
- **HTM** - JSX-like syntax without JSX
- **Tailwind CSS** - Utility-first styling (CDN)
- **Preact Router** - Client-side routing

**Testing:**
- **Playwright** - E2E testing framework
- **TypeScript** - Type-safe tests
- **GitHub Actions** - Automated CI/CD

**State Management:**
- Preact Signals for reactive stores
- No Redux, MobX, or complex state libraries
- Simple, fast, and efficient

## 📖 Additional Resources

- **FEATURES.md** - Detailed feature documentation
- **TEST_COVERAGE.md** - Complete test coverage report
- **FIRST_ADMIN_REGISTRATION.md** - First admin setup guide
- **LIVE_RELOAD.md** - Live reload system deep dive
- **README.md** - Project overview

## 🆘 Troubleshooting

### Can't access /admin
- Make sure Go backend is running
- Check port configuration (default: 3000)
- Verify static file serving is configured

### Tests failing
- Run `npm install` to install dependencies
- Make sure backend is running on localhost:3000
- Check Playwright browser installation: `npx playwright install`

### First admin registration not showing
- Backend must have no existing superadmin users
- Check endpoint: `/admin/are_there_any_superadmin`
- Verify API connection in browser DevTools

### TypeScript errors in editor
- Install recommended VS Code extensions
- Check `tsconfig.json` is present
- Restart VS Code

## 🚀 Next Steps

1. **Explore the Dashboard** - Get familiar with the UI
2. **Run the Tests** - See how E2E testing works
3. **Read FEATURES.md** - Learn about all features
4. **Customize** - Adapt to your specific needs
5. **Add Features** - Extend with new functionality

## 📝 Development Workflow

### How It Works Without Compilation

The admin panel uses **modern browser features** to run TypeScript directly:

1. **Import Maps** - Browser loads dependencies from CDN (esm.sh)
2. **ES Modules** - Native browser support for `import`/`export`
3. **TypeScript as JavaScript** - Browser ignores type annotations
4. **Go Backend** - Serves `.ts` files with `Content-Type: application/javascript`

### Basic Workflow (Manual Refresh)

```bash
# 1. Start backend
go run main.go

# 2. Open browser
http://localhost:3000/admin

# 3. Make changes to .ts files in VS Code
# No compilation needed!

# 4. Save file (Ctrl+S)

# 5. Refresh browser (F5)
# Changes appear instantly!

# 6. Run tests to verify
npm test
```

### Auto-Reload Options

**Option A: Browser DevTools Workspace (Recommended)**
1. Open Chrome DevTools (F12)
2. Sources tab → Add folder to workspace
3. Select `admin/src` directory
4. Now: Edit in VS Code → Save → Browser updates!

**Option B: Live Reload Extension**
- Install [Live Reload](https://chrome.google.com/webstore/detail/livereload/jnihajbhpnppcggbcgedagnkighmdlei) for Chrome
- Watches file changes and auto-refreshes

**Option C: Live Reload Script**
- Add the script from `admin/live-reload.html` to your `index.html` (development only)
- Polls for file changes every second and auto-reloads

### Why No Build Step?

✅ **Faster Development** - No waiting for compilation  
✅ **Simpler Setup** - No webpack/vite/rollup configuration  
✅ **Source Maps** - Full debugging with original source  
✅ **Modern Browsers** - ES2020+ features work natively  
✅ **Type Safety** - VS Code provides inline TypeScript checking  
✅ **Auto-Reload** - Backend watches files and triggers browser reload  

The browser executes the JavaScript parts and ignores TypeScript type annotations!

## ⚡ Quick Commands

```bash
# Start development with auto-reload
export APP_ENV=dev
go run main.go

# Run all tests
cd admin && npm test

# Run tests with UI
cd admin && npm run test:ui

# View test results
cd admin && npm run test:report

# Check TypeScript errors
# (Use VS Code - errors show inline)
```

## 🔧 Live Reload Technical Details

### Backend Implementation

The Go backend includes a file watcher middleware (`core/src/middleware/livereload.go`):

**Features:**
- Recursively watches `admin/src/` directory
- Computes MD5 hash of all `.ts` files
- Detects any changes (new/modified/deleted files)
- Provides two endpoints:
  - `GET /admin/dev/watch` - Server-Sent Events stream (preferred)
  - `GET /admin/dev/hash` - Polling endpoint (fallback)
- Only runs in development (disabled in production)

**How to disable:**
Set any environment other than `dev` (or leave `APP_ENV` unset) to disable live reload.

### Frontend Implementation

The browser connects to the backend watcher (`admin/index.html`):

**Primary method (SSE):**
- Establishes EventSource connection to `/admin/dev/watch`
- Receives real-time notifications when files change
- Instantly reloads the page

**Fallback method (Polling):**
- If SSE fails, polls `/admin/dev/hash` every second
- Compares hash values to detect changes
- Reloads when hash differs

**Code location:** See `<script>` block in `admin/index.html` before `</body>`
