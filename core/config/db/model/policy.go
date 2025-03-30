package model

type PolicyRule struct {
	ID                  uint                `gorm:"primaryKey"`
	CompanyID           uint                `json:"company_id"`
	CreatedByEmployeeID uint                `json:"created_by_employee_id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	SubjectAttr         string              `json:"subject_attr"`  // e.g., "user_id", "role"
	SubjectValue        string              `json:"subject_value"` // e.g., "123", "admin"
	Method              string              `json:"method"`        // e.g., "GET", "POST"
	ResourceID          uint                `json:"resource_id"`   // e.g., "/branch/:branch_id/employee/:employee_id/services" -> ID: 10
	Resource            Resource            `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE;" json:"resource"`
	Conditions          []ResourceCondition `json:"conditions"` // not stored directly if using another table/logic
}

type ResourceCondition struct {
	Attr  string `json:"attr"`  // e.g., "branch_id", "employee_id"
	Value string `json:"value"` // e.g., "456", "789"
	Op    string `json:"op"`    // e.g., "equal", "contains"
}

func (PolicyRule) TableName() string {
	return "policy_rules"
}

func (PolicyRule) Indexes() map[string]string {
	return map[string]string{
		"idx_company_resource": "CREATE INDEX idx_company_resource ON policy_rules (company_id, resource_id)",
	}
}
