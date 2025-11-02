# Admin Dashboard - Complete Feature Guide

## ✅ Implemented Features

### 0. **Authentication & Setup**

#### Login Page (`/admin`)
- Admin login form with email and password
- JWT token-based authentication
- Error handling and validation
- Loading states

**Accessing the admin panel:**
1. Start backend: `go run main.go` (default port 4000)
2. Navigate to: `http://localhost:4000/admin`
3. Backend serves admin panel from `./admin` directory
4. Admin panel communicates with API at `/api` (same origin, no CORS)

#### First Admin Registration (Automatic)
**Triggered when no superadmin exists:**
- Checks backend for existing superadmin via `GET /admin/are_there_any_superadmin`
- Shows registration form if no admin exists
- **Registration Form Fields:**
  - First Name (required)
  - Last Name (required)
  - Email (required, validated)
  - Password (required, min 8 characters)
  - Confirm Password (required, must match)
- **Validation:**
  - Client-side: Password match, length, required fields
  - Server-side: Email uniqueness, password strength
- **Post-Registration Flow:**
  1. Creates first admin with superadmin role
  2. Automatically sends verification email
  3. Shows success message with email confirmation
  4. Redirects to login page
- **Security:**
  - Email verification required before login
  - Password hashing on backend
  - Account starts as inactive/unverified

### 1. **Enhanced Dashboard** (`/`)
- Real-time statistics for Companies, Clients, and Admins
- Quick action buttons to navigate to different sections
- System information panel (Database Status, API Version, Environment, Uptime)
- Interactive stat cards with click-to-navigate
- Icons and color coding for visual clarity
- Responsive grid layout

### 2. **Companies Management** (`/companies`)
**List View:**
- Grid layout showing all companies
- Search functionality (by name or tax ID)
- Company cards showing:
  - Legal name and trade name
  - Tax ID
  - Creation date
  - Active status
- Actions: View Details, Delete

**Detail View:** (`/companies/:id`)
- Comprehensive company information
- Tabbed interface for different data:
  - **Overview Tab**: Company details and statistics
  - **Branches Tab**: List of all branches with address and phone
  - **Employees Tab**: All employees (marked if owner)
  - **Services Tab**: Services with pricing and duration
  - **Subdomains Tab**: Registered subdomains
- Delete company action (with cascade warning)
- Individual delete actions for nested resources

### 3. **Clients Management** (`/clients`)
- Table view of all registered clients
- Search functionality (by name, email, phone)
- Client information display with avatar
- Actions: View Details, Delete
- **Client Details Modal**:
  - Full client information
  - List of appointments
  - Appointment status indicators

### 4. **Admin Users** (`/users`)
- List of all admin users with roles
- Create new admin accounts
- Delete admin functionality
- Role display and management
- User status indicators

## 🏗️ Architecture

### Technology Stack

**Frontend (No Build Required):**
- **Preact 10.19.3** - Fast 3KB React alternative via ESM
- **TypeScript** - Type safety without compilation
- **Preact Signals 1.2.2** - Reactive state management
- **HTM 3.1.1** - JSX-like syntax using template literals
- **Tailwind CSS** - Utility-first styling via CDN
- **Preact Router 4.1.2** - Client-side routing

**Testing:**
- **Playwright** - E2E testing framework
- **86+ Tests** - Comprehensive coverage
- **GitHub Actions** - Automated CI/CD

### Data Stores (Preact Signals)

All state management uses Preact Signals for reactive updates:

- **authStore.ts** - Admin authentication and session management
  - Login/logout functionality
  - JWT token storage
  - User session persistence
  - Authentication state

- **adminStore.ts** - Admin user management
  - CRUD operations for admin users
  - List fetching
  - User creation and deletion
  - Role management

- **companyStore.ts** - Company CRUD operations
  - Fetch all companies
  - Fetch company by ID with nested data
  - Create company (placeholder)
  - Update company (placeholder)
  - Delete company with cascade warning
  - Search and filtering

- **clientStore.ts** - Client management and appointments
  - Fetch all clients
  - Fetch client by ID
  - Fetch client appointments
  - Delete client
  - Search and filtering

### Pages

Detailed page components:

- **Login.ts** - Admin login and first admin registration
  - Conditional rendering based on admin existence
  - Registration form with validation
  - Email verification flow
  - Success/error states

- **Dashboard.ts** - Main overview page
  - Stats display with Signals
  - Quick action buttons
  - System information panel
  - Loading states

- **Companies.ts** - Companies list
  - Grid layout with cards
  - Search functionality
  - Delete with confirmation
  - Loading and error states

- **CompanyDetail.ts** - Company details with tabs
  - Overview tab (info + stats)
  - Branches tab (locations)
  - Employees tab (staff)
  - Services tab (offerings)
  - Subdomains tab (domains)
  - Tab state management

- **Clients.ts** - Clients list with modal
  - Table view with search
  - Client details modal
  - Appointments display
  - Delete functionality

- **Users.ts** - Admin users
  - List all admin users
  - Create/delete operations
  - Role management

### Components

Reusable UI components:

- **Layout.ts** - Main layout wrapper with sidebar and header
- **Header.ts** - Top navigation bar with user info and logout
- **Sidebar.ts** - Side navigation menu with active route highlighting

### Utilities

- **api.ts** - Centralized API client
  - Automatic auth header injection
  - JSON serialization
  - Typed responses
  - Error handling
  - Request/response interceptors

- **constants.ts** - Application constants
  - API endpoints
  - Route paths
  - Configuration values

### API Integration
All stores use the centralized `api.ts` utility which:
- Automatically adds authentication headers
- Handles JSON serialization
- Provides typed responses
- Manages errors

## 🔌 Backend Integration

### Authentication Endpoints
- `POST /admin/auth/login` - Admin login (returns JWT in X-Auth-Token header)
- `POST /admin/auth/logout` - Admin logout
- `GET /admin/are_there_any_superadmin` - Check if superadmin exists
- `POST /admin/first_superadmin` - Create first admin account
- `POST /admin/send-verification-code/email/:email` - Send verification email
- `GET /admin/verify-email/:email/:code` - Verify admin email

### Admin User Endpoints (Existing)
- `GET /admin` - List all admins
- `POST /admin` - Create admin
- `GET /admin/:id` - Get admin details
- `PUT /admin/:id` - Update admin
- `DELETE /admin/:id` - Delete admin

### Companies Endpoints
- `GET /admin/companies` - **⚠️ May need to be created** - List all companies for admin view
- `GET /company/:id` - Get company details with nested data (branches, employees, services, subdomains)
- `PUT /company/:id` - Update company
- `DELETE /company/:id` - Delete company (cascades to nested resources)

### Clients Endpoints
- `GET /admin/clients` - **⚠️ May need to be created** - List all clients across all companies
- `GET /client/:id` - Get client details
- `GET /client/:id/appointments` - Get client appointments
- `DELETE /client/:id` - Delete client

### Nested Resource Endpoints (May need creation)
Individual delete endpoints for nested resources:
- `DELETE /branch/:id` - Delete individual branch
- `DELETE /employee/:id` - Delete individual employee
- `DELETE /service/:id` - Delete individual service
- `DELETE /subdomain/:id` - Delete subdomain

## ⚠️ Backend Endpoints Status

### ✅ Already Implemented (in admin.go)
- Authentication endpoints
- Admin user CRUD
- Email verification
- Password reset

### ⚠️ May Need Creation
These endpoints are used by the frontend but may not exist yet:

1. **`GET /admin/companies`** - List all companies for admin view
   - Should return all companies across all schemas
   - Include basic company info (name, tax_id, created_at, etc.)

2. **`GET /admin/clients`** - List all clients across companies
   - Should aggregate clients from all companies
   - Include client basic info (name, email, phone, created_at)

3. **Individual nested resource deletes:**
   - `DELETE /branch/:id`
   - `DELETE /employee/:id`
   - `DELETE /service/:id`
   - `DELETE /subdomain/:id`

## 🎨 UI/UX Features

- **Responsive Design**: Works on desktop, tablet, and mobile
- **Real-time Search**: Instant filtering on Companies and Clients pages
- **Modal Dialogs**: Client details shown in overlay modal with backdrop
- **Confirmation Dialogs**: All delete actions require browser confirmation
- **Loading States**: Visual indicators during API calls
- **Error Handling**: User-friendly error messages from API
- **Navigation**: Active route highlighting in sidebar
- **Icons**: Emoji icons for visual clarity and appeal
- **Color Coding**: Status badges (Active, Confirmed, Cancelled)
- **Empty States**: Helpful messages when no data exists
- **Form Validation**: Client-side validation with error messages
- **Required Field Indicators**: Asterisks (*) for required fields
- **Password Hints**: Helper text for password requirements
- **Success Feedback**: Visual confirmation after successful actions
- **Hover Effects**: Cards and buttons with transition effects
- **Grid Layouts**: Responsive grid system for cards
- **Table Views**: Clean, organized data tables
- **Badges**: Status indicators and count badges on tabs
- **Avatars**: Client initials in colored circles

## 🧪 Testing

### Comprehensive E2E Test Suite (86+ Tests)

**Test Files:**
1. **login.spec.ts** (6 tests) - Login form and authentication
2. **first-admin-registration.spec.ts** (11 tests) - First admin setup
3. **dashboard.spec.ts** (14 tests) - Dashboard functionality
4. **companies.spec.ts** (13 tests) - Companies list
5. **company-detail.spec.ts** (17 tests) - Company details with tabs
6. **clients.spec.ts** (20 tests) - Clients management
7. **users.spec.ts** (existing) - Admin users
8. **navigation.spec.ts** (7 tests) - Routing and navigation
9. **integration.spec.ts** (24 tests) - Integration and edge cases

**Coverage Areas:**
- ✅ Authentication and first admin setup
- ✅ All CRUD operations
- ✅ Search and filtering
- ✅ Modal interactions
- ✅ Tab switching
- ✅ Delete confirmations
- ✅ Loading states
- ✅ Error handling
- ✅ Empty states
- ✅ Navigation flows
- ✅ Data consistency
- ✅ Accessibility
- ✅ Edge cases

**CI/CD Integration:**
- Automated testing on GitHub Actions
- Runs on push to main/admin branches
- Runs on pull requests
- Test artifacts saved for review

See `TEST_COVERAGE.md` for detailed test documentation.

## 🚀 Potential Enhancements

### High Priority
1. **Edit Functionality** - Update company/client information inline
2. **Appointments View** - Dedicated page for all appointments with calendar
3. **Real-time Updates** - WebSocket integration for live data
4. **Audit Logs** - Track all admin actions with timestamps

### Medium Priority
5. **Analytics Dashboard** - Charts and graphs for business insights (Chart.js/Recharts)
6. **Export Data** - CSV/PDF export functionality for reports
7. **Bulk Actions** - Select multiple items for bulk delete/update
8. **Advanced Filters** - Filter by date range, status, categories
9. **Pagination** - Handle large datasets efficiently
10. **Settings Page** - System configuration and preferences

### Nice to Have
11. **Two-Factor Authentication** - Enhanced security for admins
12. **Profile Pictures** - Upload and manage admin avatars
13. **Dark Mode** - Theme switcher for dark/light modes
14. **Email Templates** - Customize verification and notification emails
15. **Role Permissions** - Granular permission system beyond superadmin
16. **Activity Dashboard** - Real-time activity feed
17. **Keyboard Shortcuts** - Power user features
18. **Search Suggestions** - Autocomplete in search fields
19. **Drag & Drop** - Reorder items or upload files
20. **Mobile App** - Native mobile app version

## 📝 Usage Guide

### First Time Setup
1. Start Go backend
2. Navigate to `/admin`
3. Create first admin account (auto-prompted)
4. Verify email
5. Login

### Daily Usage
1. **Login** as admin
2. **View Dashboard** for overview and statistics
3. **Navigate to Companies** to manage tenants
4. **Click "View Details"** on a company to see:
   - Company information
   - Branches (locations)
   - Employees (staff members)
   - Services (offerings)
   - Subdomains (registered domains)
5. **Navigate to Clients** to see all registered clients
6. **Click "View Details"** on a client to see:
   - Client information
   - Appointment history
   - Status indicators
7. **Manage Admin Users** in Users section
8. **Use Delete buttons** with caution (confirmation required)

### Best Practices
- ✅ Always verify before deleting (actions are permanent)
- ✅ Use search to find specific items quickly
- ✅ Check client appointments before deleting clients
- ✅ Review company data (branches, employees) before deletion
- ✅ Keep admin accounts secure with strong passwords
- ✅ Verify email addresses for all new admins
- ✅ Regular backups recommended before bulk operations

## 🔐 Security Features

### Implemented
- ✅ JWT token-based authentication
- ✅ Email verification required for admin accounts
- ✅ Password hashing on backend (bcrypt)
- ✅ Protected routes (authentication required)
- ✅ Confirmation dialogs for destructive actions
- ✅ HTTPS recommended for production
- ✅ Input sanitization
- ✅ SQL injection prevention (GORM ORM)
- ✅ XSS protection (Preact auto-escaping)
- ✅ CSRF tokens (backend implementation)

### Recommended
- ⚠️ Rate limiting on authentication endpoints
- ⚠️ CAPTCHA for registration/login
- ⚠️ Two-factor authentication for superadmin
- ⚠️ Session timeout configuration
- ⚠️ IP whitelisting for admin access
- ⚠️ Regular security audits
- ⚠️ Monitor failed login attempts
- ⚠️ Automated vulnerability scanning

## 📚 Documentation

Complete documentation available:

- **README.md** - Project overview and quick setup
- **docs/QUICKSTART.md** - Quick start guide with examples
- **docs/FEATURES.md** - This file (complete feature guide)
- **docs/TEST_COVERAGE.md** - Detailed test coverage report
- **docs/FIRST_ADMIN_REGISTRATION.md** - First admin setup documentation
- **tests/README.md** - Testing documentation

## 🛠️ Development Notes

### Live Reload System

**Automatic browser refresh on file changes:**

The admin panel includes a sophisticated live reload system that eliminates manual browser refreshes during development.

**Backend Implementation:**
- File watcher middleware in `core/src/middleware/livereload.go`
- Recursively monitors `admin/src/` directory
- Detects changes to any `.ts` file (new, modified, deleted)
- Computes MD5 hash of all TypeScript files
- Only runs when `APP_ENV=dev`

**Frontend Implementation:**
- Script in `admin/index.html` before `</body>`
- Primary: Server-Sent Events (SSE) connection to `/admin/dev/watch`
- Fallback: Polling `/admin/dev/hash` endpoint every second
- Automatically reloads browser when changes detected

**What triggers reload:**
- Editing existing `.ts` files
- Creating new `.ts` files
- Deleting `.ts` files  
- Adding new folders with `.ts` files
- Any modification in `admin/src/` directory tree

**How to use:**
1. Start backend: `export APP_ENV=dev && go run main.go`
2. Open browser: `http://localhost:4000/admin`
3. Edit any file in `admin/src/`
4. Save → Browser reloads automatically! 🎉

**Console indicators:**
- Backend: `🔄 Live reload enabled - watching ./admin/src`
- Browser: `🔄 File changed: ./admin/src`

### No Build Step
- TypeScript files run directly in browser
- No webpack, vite, or bundler needed
- Instant refresh to see changes
- Full source maps for debugging

### State Management
- Preact Signals for reactive state
- Simple, fast, and efficient
- No Redux or MobX complexity
- Automatic UI updates

### Styling
- Tailwind CSS via CDN
- Utility-first approach
- No CSS preprocessing
- Responsive by default

### Type Safety
- Full TypeScript support
- Type checking in VS Code
- Type definitions for all stores
- Typed API responses

## 🌐 Browser Support

**Minimum Requirements:**
- Modern browsers with ES2020+ support
- Chrome 87+ (December 2020)
- Firefox 78+ (June 2020)
- Safari 14+ (September 2020)
- Edge 87+ (December 2020)

**Not Supported:**
- Internet Explorer (any version)
- Legacy browsers without ES Modules support

## 📊 Performance

- **Initial Load**: ~50KB (Preact + HTM + Signals via CDN)
- **Page Size**: Minimal (no bundled JavaScript)
- **Load Time**: Fast (CDN cached dependencies)
- **Runtime**: Efficient (Preact's small size and Signals)

---

Last Updated: November 2, 2025  
Version: 1.0.0
