package DTO

type TimeRange struct {
	Start string `json:"start" example:"09:00"`
	End   string `json:"end" example:"17:00"`
}

type UpdateEmployeeSwagger struct {
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Clark"`
}

type CreateEmployee struct {
	CompanyID uint     `gorm:"not null" json:"company_id"`
	Name      string   `gorm:"not null" json:"name" example:"Joseph"`
	Surname   string   `json:"surname" example:"Doe"`
	Role      string   `json:"role" example:"user"`
	Email     string   `gorm:"not null;unique" json:"email" example:"joseph.doe@example.com"`
	Phone     string   `gorm:"not null;unique" json:"phone" example:"+15555555551"`
}

type Employee struct {
	ID             uint               `json:"id" example:"1"`
	UserID         uint               `json:"user_id" example:"1"`
	User           User               `json:"user"`
	CompanyID      uint               `json:"company_id" example:"1"`
	Company        Company            `json:"company" example:"1"`
	Branches       []BranchPopulated  `json:"branches" example:"[]"`
	Services       []ServicePopulated `json:"services" example:"[]"`
	AvailableSlots []TimeRange        `json:"available_slots" example:"[{\"start\":\"09:00\", \"end\":\"17:00\"}]"`
	Appointments   []Appointment      `json:"appointments" example:"[]"`
	Tags           []string           `json:"tag" example:"[\"super-admin\", \"branch-manager\"]"`
}

type EmployeePopulated struct {
	ID        uint `json:"id" example:"1"`
	UserID    uint `json:"user_id" example:"1"`
	CompanyID uint `json:"company_id" example:"1"`
}
