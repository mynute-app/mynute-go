package DTO

type Employee struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Branches []Branch  `json:"branches"`
	Services []Service `json:"services"`
	Company  Company   `json:"company"`
}
