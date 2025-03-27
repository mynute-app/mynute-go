package model

type PolicyRule struct {
	ID            uint `gorm:"primaryKey"`
	CompanyID     uint
	CreatedByUser uint
	Name          string `gorm:"uniqueIndex"`
	Description   string
	SubjectAttr   string              // e.g., "user_id", "role"
	SubjectValue  string              // e.g., "123", "admin"
	Method        string              // e.g., "GET", "POST"
	Path          string              // e.g., "/branch/:branch_id/employee/:employee_id/services"
	Conditions    []ResourceCondition `json:"conditions"` // not stored directly if using another table/logic
}

type ResourceCondition struct {
	Attr  string `json:"attr"`  // e.g., "branch_id", "employee_id"
	Value string `json:"value"` // e.g., "456", "789"
	Op    string `json:"op"`    // e.g., "equal", "contains"
}
