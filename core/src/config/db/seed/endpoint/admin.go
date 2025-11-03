package endpointSeed

import (
	"mynute-go/auth/config/db/model"
	"mynute-go/core/src/config/namespace"
)

var AdminLoginByPassword = &model.EndPoint{
	Path:           "/admin/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "AdminLoginByPassword",
	Description:    "Admin login",
}

var AreThereAnyAdmin = &model.EndPoint{
	Path:             "/admin/are_there_any_superadmin",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "AreThereAnyAdmin",
	Description:      "Check if any superadmin exists",
	DenyUnauthorized: false, // Public endpoint for initial setup check
}

var CreateFirstAdmin = &model.EndPoint{
	Path:             "/admin/first_superadmin",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateFirstAdmin",
	Description:      "Create the first superadmin (only works if no superadmin exists)",
	DenyUnauthorized: false, // Public endpoint for initial setup
}

var GetAdminByID = &model.EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminByID",
	Description:      "Get admin by ID",
	DenyUnauthorized: true,
}

var GetAdminByEmail = &model.EndPoint{
	Path:             "/admin/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminByEmail",
	Description:      "Get admin by email",
	DenyUnauthorized: false, // Allow public access for post-login user data fetch
}

var ListAdmins = &model.EndPoint{
	Path:             "/admin",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "ListAdmins",
	Description:      "List all admins",
	DenyUnauthorized: true,
}

var CreateAdmin = &model.EndPoint{
	Path:             "/admin",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAdmin",
	Description:      "Create a new admin",
	DenyUnauthorized: false, // Bootstrap case: allows first admin creation without auth
}

var UpdateAdminByID = &model.EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAdminByID",
	Description:      "Update admin by ID",
	DenyUnauthorized: true,
}

var DeleteAdminByID = &model.EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteAdminByID",
	Description:      "Delete admin by ID",
	DenyUnauthorized: true,
}

var ResetAdminPasswordByEmail = &model.EndPoint{
	Path:             "/admin/reset-password/:email",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "ResetAdminPasswordByEmail",
	Description:      "Reset admin password by email",
	DenyUnauthorized: false, // Public endpoint for password resetmodel.
}

var SendAdminVerificationCodeByEmail = &model.EndPoint{
	Path:             "/admin/send-verification-code/email/:email",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "SendAdminVerificationCodeByEmail",
	Description:      "Send verification code to admin email",
	DenyUnauthorized: false, // Public endpoint for email verification
}

var VerifyAdminEmail = &model.EndPoint{
	Path:             "/admin/verify-email/:email/:code",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "VerifyAdminEmail",
	Description:      "Verify admin email with code",
	DenyUnauthorized: false, // Public endpoint for email verification
}

var ListAdminRoles = &model.EndPoint{
	Path:             "/admin/role",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "ListAdminRoles",
	Description:      "List all admin roles",
	DenyUnauthorized: true,
}

var CreateAdminRole = &model.EndPoint{
	Path:             "/admin/role",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAdminRole",
	Description:      "Create a new admin role",
	DenyUnauthorized: true,
}

var GetAdminRoleByID = &model.EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminRoleByID",
	Description:      "Get admin role by ID",
	DenyUnauthorized: true,
}

var UpdateAdminRoleByID = &model.EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAdminRoleByID",
	Description:      "Update role by ID",
	DenyUnauthorized: true,
}

var DeleteAdminRoleByID = &model.EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteAdminRoleByID",
	Description:      "Delete role by ID",
	DenyUnauthorized: true,
}

