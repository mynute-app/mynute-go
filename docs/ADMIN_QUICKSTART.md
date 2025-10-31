# Admin System - Quick Start Guide

## 🚀 Quick Setup (5 minutes)

### Step 1: Run Migration
```powershell
cd c:\code\mynute-go
.\scripts\migrate.ps1 -action up
```

### Step 2: Seed Admin Data
```powershell
go run cmd/seed-admin/main.go
```

### Step 3: Login
```bash
curl -X POST http://localhost:3000/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@mynute.com","password":"Admin@123456"}'
```

**Done!** You now have a working admin system.

---

## 📝 Essential Commands

### Login
```bash
curl -X POST http://localhost:3000/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"EMAIL","password":"PASSWORD"}'
```

### Get My Info
```bash
curl http://localhost:3000/admin/auth/me \
  -H "X-Auth-Token: YOUR_TOKEN"
```

### Create New Admin (Superadmin only)
```bash
curl -X POST http://localhost:3000/admin/create \
  -H "Content-Type: application/json" \
  -H "X-Auth-Token: YOUR_TOKEN" \
  -d '{
    "name": "New Admin",
    "email": "newadmin@example.com",
    "password": "SecurePass123!",
    "is_active": true,
    "roles": ["support"]
  }'
```

### List All Admins
```bash
curl http://localhost:3000/admin/list \
  -H "X-Auth-Token: YOUR_TOKEN"
```

---

## 🔑 Default Credentials

After running the seeder:

- **Email:** `admin@mynute.com`
- **Password:** `Admin@123456`

⚠️ **Change immediately in production!**

---

## 🛡️ Admin Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| `superadmin` | Full system access | Can create/edit/delete admins and roles |
| `support` | Customer support | Read-only access to all tenant data |
| `auditor` | Compliance/auditing | Read-only access for auditing purposes |

---

## 🔧 Environment Variables

Add to `.env`:

```env
DEFAULT_ADMIN_EMAIL=admin@yourcompany.com
DEFAULT_ADMIN_PASSWORD=YourSecurePassword123!
DEFAULT_ADMIN_NAME=System Administrator
JWT_SECRET=your-secret-key-here
```

---

## 📚 Full Documentation

See [ADMIN_SYSTEM.md](./ADMIN_SYSTEM.md) for complete documentation.

---

## ⚠️ Security Checklist

- [ ] Changed default admin password
- [ ] Set strong JWT_SECRET in .env
- [ ] Using HTTPS in production
- [ ] Limited superadmin role to trusted users only
- [ ] Implemented admin action logging

---

## 🐛 Troubleshooting

**Problem:** Can't login
- Check email/password are correct
- Verify admin is active: `UPDATE public.admins SET is_active = true WHERE email = 'your@email.com'`

**Problem:** "Superadmin role required"
- Verify role assignment:
  ```sql
  SELECT a.email, r.name 
  FROM public.admins a
  JOIN public.admin_role_admins ara ON a.id = ara.admin_id
  JOIN public.role_admins r ON ara.role_admin_id = r.id;
  ```

**Problem:** Token expired
- Login again or use `/admin/auth/refresh`

---

## 📞 Support

For detailed documentation, see:
- [ADMIN_SYSTEM.md](./ADMIN_SYSTEM.md) - Complete admin system documentation
- [MIGRATIONS.md](./MIGRATIONS.md) - Database migration guide
