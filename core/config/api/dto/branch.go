package DTO

type Branch struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Company   Company    `json:"company"`
	Employees []Employee `json:"employees"`
	Services  []Service  `json:"services"`
}