package endpointSeed

import (
	"mynute-go/services/auth/config/db/model"
	resourceSeed "mynute-go/services/core/src/config/db/seed/resource"
	"mynute-go/services/core/src/config/namespace"
)

var CreateClient = &model.EndPoint{
	Path:           "/client",
	Method:         namespace.CreateActionMethod,
	ControllerName: "CreateClient",
	Description:    "Create client",
}

var LoginClient = &model.EndPoint{
	Path:           "/client/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginClientByPassword",
	Description:    "Login client",
}

var LoginClientByEmailCode = &model.EndPoint{
	Path:           "/client/login-with-code",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginClientByEmailCode",
	Description:    "Login client by email code",
}

var SendLoginCodeToClientEmail = &model.EndPoint{
	Path:           "/client/send-login-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendClientLoginValidationCodeByEmail",
	Description:    "Send login code to client email",
}

var ResetClientPasswordByEmail = &model.EndPoint{
	Path:           "/client/reset-password/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "ResetClientPasswordByEmail",
	Description:    "Reset client password by email",
}

var SendClientVerificationCodeByEmail = &model.EndPoint{
	Path:           "/client/send-verification-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendClientVerificationCodeByEmail",
	Description:    "Send verification code to client email",
}

var VerifyClientEmail = &model.EndPoint{
	Path:           "/client/verify-email/:email/:code",
	Method:         namespace.ViewActionMethod,
	ControllerName: "VerifyClientEmail",
	Description:    "Verify client email code",
}

var GetClientByEmail = &model.EndPoint{
	Path:             "/client/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetClientByEmail",
	Description:      "View client by email",
	DenyUnauthorized: false,
	Resource:         resourceSeed.Client,
}

var GetClientById = &model.EndPoint{
	Path:             "/client/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetClientById",
	Description:      "View client by ID",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Client,
}

var UpdateClientById = &model.EndPoint{
	Path:             "/client/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateClientById",
	Description:      "Update client by ID",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Client,
}

var DeleteClientById = &model.EndPoint{
	Path:             "/client/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteClientById",
	Description:      "Delete client by ID",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Client,
}

var UpdateClientImages = &model.EndPoint{
	Path:             "/client/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateClientImages",
	Description:      "Update client design images",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Client,
}

var DeleteClientImage = &model.EndPoint{
	Path:             "/client/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteClientImage",
	Description:      "Delete client design images",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Client,
}

