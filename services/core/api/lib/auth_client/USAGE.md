# Auth Client Usage Guide

The auth client provides a complete interface to the auth service API, allowing the core service to manage endpoints, policies, users, and perform authentication/authorization checks.

## Initialization

```go
import authClient "mynute-go/services/core/api/lib/auth_client"

client := authClient.NewAuthClient()
// Uses AUTH_SERVICE_URL env var, defaults to http://localhost:4001
```

## Health Check

```go
if !client.IsAvailable() {
    log.Fatal("Auth service is not available")
}
```

## Endpoint Management

### Fetch All Endpoints
```go
endpoints, err := client.FetchEndpoints()
if err != nil {
    log.Fatal(err)
}
for _, ep := range endpoints {
    fmt.Printf("%s %s -> %s\n", ep.Method, ep.Path, ep.ControllerName)
}
```

### Create Endpoint
```go
endpoint, err := client.CreateEndpoint(authClient.CreateEndpointRequest{
    ControllerName:   "GetUsers",
    Description:      "Retrieve all users",
    Method:           "GET",
    Path:             "/api/users",
    DenyUnauthorized: true,
    NeedsCompanyId:   false,
})
```

### Get Endpoint by ID
```go
endpoint, err := client.GetEndpoint(endpointID)
```

### Update Endpoint
```go
description := "Updated description"
endpoint, err := client.UpdateEndpoint(endpointID, authClient.UpdateEndpointRequest{
    Description: &description,
})
```

### Delete Endpoint
```go
err := client.DeleteEndpoint(endpointID)
```

## Policy Management

### Fetch All Policies
```go
policies, err := client.FetchPolicies()
if err != nil {
    log.Fatal(err)
}
```

### Create Policy
```go
conditions := json.RawMessage(`{
    "logic_type": "AND",
    "children": [
        {
            "leaf": {
                "check_type": "user_type",
                "operator": "equals",
                "value": "employee"
            }
        }
    ]
}`)

policy, err := client.CreatePolicy(authClient.CreatePolicyRequest{
    Name:        "Employee Access",
    Description: "Allow employees to access this endpoint",
    Effect:      "Allow",
    EndPointID:  endpointID,
    Conditions:  conditions,
})
```

### Get Policy by ID
```go
policy, err := client.GetPolicy(policyID)
```

### Update Policy
```go
newName := "Updated Policy Name"
policy, err := client.UpdatePolicy(policyID, authClient.UpdatePolicyRequest{
    Name: &newName,
})
```

### Delete Policy
```go
err := client.DeletePolicy(policyID)
```

## User Management

### Get Client by Email
```go
client, err := client.GetClientByEmail("user@example.com")
```

### Get Client by ID
```go
client, err := client.GetClientByID(clientID)
```

### Get Employee by Email
```go
employee, err := client.GetEmployeeByEmail("employee@example.com")
```

### Get Employee by ID
```go
employee, err := client.GetEmployeeByID(employeeID)
```

### Get Admin by ID
```go
admin, err := client.GetAdminByID(adminID)
```

### List All Admins
```go
admins, err := client.ListAdmins()
```

## Authentication

### Validate User Token
```go
result, err := client.ValidateToken("eyJhbGciOiJIUzI1NiIs...")
if err != nil {
    log.Printf("Invalid token: %v", err)
    return
}
if result.Valid {
    fmt.Printf("Token valid for user: %s\n", result.Email)
}
```

### Validate Admin Token
```go
result, err := client.ValidateAdminToken("eyJhbGciOiJIUzI1NiIs...")
if err != nil {
    log.Printf("Invalid admin token: %v", err)
    return
}
```

## Authorization

### Check Access
```go
response, err := client.CheckAccess(authClient.AccessCheckRequest{
    UserID:    userID,
    Method:    "POST",
    Path:      "/api/appointments",
    CompanyID: &companyID, // Optional
})

if err != nil {
    log.Printf("Access check failed: %v", err)
    return
}

if response.Allowed {
    fmt.Println("Access granted")
} else {
    fmt.Printf("Access denied: %s\n", response.Reason)
}
```

## Error Handling

All methods return errors that include:
- Network errors
- HTTP status errors with response body
- JSON decoding errors

```go
endpoint, err := client.GetEndpoint(id)
if err != nil {
    if strings.Contains(err.Error(), "status 404") {
        log.Println("Endpoint not found")
    } else if strings.Contains(err.Error(), "request failed") {
        log.Println("Network error")
    } else {
        log.Printf("Unknown error: %v", err)
    }
    return
}
```

## Environment Configuration

```bash
# Set custom auth service URL
export AUTH_SERVICE_URL=http://auth-service:4001

# Or in .env file
AUTH_SERVICE_URL=http://auth-service:4001
```

## Complete Example

```go
package main

import (
    "log"
    authClient "mynute-go/services/core/api/lib/auth_client"
)

func main() {
    // Initialize client
    client := authClient.NewAuthClient()
    
    // Check availability
    if !client.IsAvailable() {
        log.Fatal("Auth service is not available")
    }
    
    // Fetch all endpoints
    endpoints, err := client.FetchEndpoints()
    if err != nil {
        log.Fatalf("Failed to fetch endpoints: %v", err)
    }
    
    log.Printf("Found %d endpoints", len(endpoints))
    
    // Check user access
    userID := uuid.MustParse("...")
    companyID := uuid.MustParse("...")
    
    access, err := client.CheckAccess(authClient.AccessCheckRequest{
        UserID:    userID,
        Method:    "POST",
        Path:      "/api/appointments",
        CompanyID: &companyID,
    })
    
    if err != nil {
        log.Fatalf("Access check failed: %v", err)
    }
    
    if access.Allowed {
        log.Println("User has access to create appointments")
    } else {
        log.Printf("Access denied: %s", access.Reason)
    }
}
```
