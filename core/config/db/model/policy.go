package model

type PolicyRule struct {
	ID             uint                `gorm:"primaryKey"`
	CompanyID      uint
	CreatedByUser  uint
	Name           string              `gorm:"uniqueIndex"`
	Description    string
	SubjectAttr    string              // e.g., "user_id", "role"
	SubjectValue   string              // e.g., "admin", "123"
	Method         string              // e.g., "GET", "POST"
	Path           string              // e.g., "/branch/:branch_id/employee/:employee_id/services"
	Conditions     []ResourceCondition `gorm:"-" json:"conditions"` // not stored directly if using another table/logic
}

type ResourceCondition struct {
	Attr  string `json:"attr"`
	Value string `json:"value"`
	Op    string `json:"op"` // e.g., "equal", "contains"
}