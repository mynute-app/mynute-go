# Mynute Admin Panel

A lightweight, modern admin panel built with:
- **Preact** - Fast 3KB React alternative
- **TypeScript** - Type safety without build step
- **Preact Signals** - Reactive state management
- **HTM** - JSX-like syntax in the browser
- **Tailwind CSS** - Utility-first styling
- **Playwright** - Comprehensive E2E testing
- **Live Reload** - Automatic browser refresh on file changes

## ⚡ Quick Start

```bash
# Start the Go backend (on Windows use PowerShell)
# For live reload, set APP_ENV=dev
$env:APP_ENV="dev"
go run main.go

# Backend starts on port 4000 (default, configurable via APP_PORT)

# Open browser and navigate to:
http://localhost:4000/admin

# Edit .ts files → Save → Browser reloads automatically! 🎉
```

**How it works:**
- Backend serves admin panel at `/admin` (maps to `./admin` directory)
- Admin panel makes API calls to `/api` (relative URL, same origin)
- Live reload watches `./admin/src/` and notifies browser via SSE when `APP_ENV=dev`

## Features

✅ **No build step required** - Run directly in the browser via ESM  
✅ **Auto-reload** - Backend watches files and triggers browser reload  
✅ **TypeScript support** - Full type safety with `.ts` files  
✅ **Modern React patterns** - Hooks, components, routing  
✅ **Preact Signals** - Simple, fast state management  
✅ **Responsive UI** - Built with Tailwind CSS  
✅ **First Admin Setup** - Automatic registration flow for initial admin  
✅ **Full CRUD Operations** - Companies, Clients, Admin users  
✅ **86+ E2E Tests** - Comprehensive Playwright test suite  
✅ **Email Verification** - Admin account verification flow  

## Pages & Features

### 🔐 Authentication
- **Login Page** - Admin authentication with JWT
- **First Admin Registration** - Automatic setup for first superadmin
- **Email Verification** - Required for admin accounts

### 📊 Dashboard (`/`)
- Real-time statistics (Companies, Clients, Admins)
- Quick action buttons for navigation
- System information panel
- Clickable stat cards

### 🏢 Companies Management (`/companies`)
- Grid view with search functionality
- Company cards with details (name, tax ID, status)
- Individual company detail page with tabs:
  - Overview (information & statistics)
  - Branches (locations with address/phone)
  - Employees (staff members)
  - Services (offerings with pricing)
  - Subdomains (registered domains)
- CRUD operations with confirmation dialogs

### 👥 Clients Management (`/clients`)
- Table view with search
- Client avatars with initials
- Details modal showing:
  - Full client information
  - Appointment history
  - Status indicators
- Delete functionality

### 🔑 Admin Users (`/users`)
- List all admin users with roles
- Create/Delete admin accounts
- Role management  

## Structure

```
/admin
  index.html              # Entry point with CDN imports
  package.json            # Dependencies (Playwright, TypeScript)
  playwright.config.ts    # E2E test configuration
  tsconfig.json          # TypeScript configuration
  /docs                  # Documentation
    FEATURES.md          # Feature summary
    QUICKSTART.md        # Quick start guide
    TEST_COVERAGE.md     # Test coverage report
    FIRST_ADMIN_REGISTRATION.md  # First admin setup docs
    LIVE_RELOAD.md       # Live reload system documentation
  /src
    main.ts              # Application bootstrap
    App.ts               # Main app component with routing
    types.ts             # TypeScript type definitions
    global.d.ts          # CDN module declarations
    /components
      Layout.ts          # Main layout wrapper
      Header.ts          # Top navigation bar
      Sidebar.ts         # Side navigation menu
    /pages
      Login.ts           # Login + First admin registration
      Dashboard.ts       # Main dashboard with stats
      Companies.ts       # Companies list view
      CompanyDetail.ts   # Company details with tabs
      Clients.ts         # Clients list with modal
      Users.ts           # Admin users management
    /stores
      authStore.ts       # Authentication state (Signals)
      adminStore.ts      # Admin users CRUD (Signals)
      companyStore.ts    # Companies CRUD (Signals)
      clientStore.ts     # Clients & appointments (Signals)
    /utils
      api.ts             # API client with auth headers
      constants.ts       # App constants
  /tests                 # E2E test suite (86+ tests)
    fixtures.ts          # Test fixtures & auth
    login.spec.ts        # Login tests
    first-admin-registration.spec.ts  # First admin setup tests
    dashboard.spec.ts    # Dashboard tests
    companies.spec.ts    # Companies list tests
    company-detail.spec.ts  # Company detail tests
    clients.spec.ts      # Clients tests
    users.spec.ts        # Admin users tests
    navigation.spec.ts   # Navigation tests
    integration.spec.ts  # Integration & error tests
```

## Getting Started

### First Time Setup

1. **Start your Go backend** (it serves the admin folder):
   ```bash
   go run main.go
   ```

2. **Open in browser**: `http://localhost:3000/admin`

3. **Create first admin**:
   - If no admin exists, you'll see a registration form
   - Fill in your details (name, email, password)
   - Click "Create Admin Account"
   - Check your email for verification link
   - Return to login page

4. **Login**:
   - Enter your email and password
   - You'll be redirected to the dashboard

### Development Workflow

1. **No compilation needed** - Edit `.ts` files and refresh browser
2. **TypeScript checking** - Use VS Code for inline type checking
3. **Browser DevTools** - Full source maps available
4. **Hot reload** - Simple refresh to see changes

### Running Tests

```bash
cd admin

# Install dependencies
npm install

# Run all tests (headless)
npm test

# Run tests in UI mode (recommended)
npm run test:ui

# Run with visible browser
npm run test:headed

# View test report
npm run test:report
```

### API Configuration

Update the API base URL in `src/utils/api.ts`:
```typescript
const API_BASE_URL = '/api'; // Adjust to your backend URL
```

## State Management with Signals

```typescript
import { signal } from '@preact/signals';

// Create signals
const count = signal(0);

// Read value
console.log(count.value); // 0

// Update value
count.value++;

// Components auto-update when signals change
function Counter() {
  return html`<div>Count: ${count.value}</div>`;
}
```

## API Integration

The `api` utility in `src/utils/api.ts` automatically:
- Adds authentication headers
- Handles JSON serialization
- Provides typed responses
- Manages errors

```typescript
// Example usage
const users = await api.get<User[]>('/admin/users');
const created = await api.post<User>('/admin/users', { name: 'John' });
```

## Routing

Uses `preact-router` for client-side routing:

```typescript
html`
  <${Router}>
    <${Login} default />
    <${Dashboard} path="/" />
    <${Companies} path="/companies" />
    <${CompanyDetail} path="/companies/:id" />
    <${Clients} path="/clients" />
    <${Users} path="/users" />
  </${Router}>
`
```

## Testing

**Comprehensive E2E test suite with Playwright:**

- 86+ tests covering all features
- Authenticated page fixtures
- API mocking for isolated tests
- CI/CD integration with GitHub Actions
- ~98% feature coverage

**Test categories:**
- Authentication & first admin setup
- Dashboard functionality
- Companies CRUD operations
- Company detail views with tabs
- Clients management
- Navigation & routing
- Error handling & edge cases
- Data consistency
- Accessibility & UX

See `docs/TEST_COVERAGE.md` for detailed test documentation.

## Documentation

- **README.md** - This file (overview and setup)
- **docs/QUICKSTART.md** - Quick start guide
- **docs/FEATURES.md** - Feature summary and architecture
- **docs/TEST_COVERAGE.md** - Test coverage report
- **docs/FIRST_ADMIN_REGISTRATION.md** - First admin setup feature docs
- **docs/LIVE_RELOAD.md** - Live reload system documentation

## Backend Integration

### Required Endpoints

The admin panel expects these backend endpoints:

**Authentication:**
- `POST /admin/auth/login` - Admin login
- `GET /admin/are_there_any_superadmin` - Check for existing admins
- `POST /admin/first_superadmin` - Create first admin
- `POST /admin/send-verification-code/email/:email` - Send verification email

**Admin Users:**
- `GET /admin` - List all admins
- `POST /admin` - Create admin
- `GET /admin/:id` - Get admin details
- `PUT /admin/:id` - Update admin
- `DELETE /admin/:id` - Delete admin

**Companies:**
- `GET /admin/companies` - List all companies (may need creation)
- `GET /company/:id` - Get company with nested data
- `DELETE /company/:id` - Delete company

**Clients:**
- `GET /admin/clients` - List all clients (may need creation)
- `GET /client/:id` - Get client details
- `GET /client/:id/appointments` - Get appointments
- `DELETE /client/:id` - Delete client

See `docs/FEATURES.md` for complete endpoint list.

## Browser Support

- Modern browsers with ES2020+ support
- Chrome, Firefox, Safari, Edge (latest versions)
