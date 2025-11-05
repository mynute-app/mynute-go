package endpointSeed

import (
	"mynute-go/services/auth/config/db/model"
	"mynute-go/services/core/src/config/namespace"
)

var BeginAuthProviderCallback = &model.EndPoint{
	Path:           "/auth/oauth/:provider",
	Method:         namespace.ViewActionMethod,
	ControllerName: "BeginAuthProviderCallback",
	Description:    "Begin auth provider callback",
}

var GetAuthCallbackFunction = &model.EndPoint{
	Path:           "/auth/oauth/:provider/callback",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetAuthCallbackFunction",
	Description:    "View auth callback function",
}

var LogoutProvider = &model.EndPoint{
	Path:           "/auth/oauth/logout",
	Method:         namespace.ViewActionMethod,
	ControllerName: "LogoutProvider",
	Description:    "Logout provider",
}

