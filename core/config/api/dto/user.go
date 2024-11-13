package DTO

type User struct {
	ID        uint               `json:"id"`
	CompanyID uint               `json:"company_id"`
	Name      string             `json:"name"`
	Surname   string             `json:"surname"`
	Email     string             `json:"email"`
	Password string 		   `json:"password"`
	Phone     string             `json:"phone"`
	Branches  []BranchPopulated  `json:"branches"`
	Services  []ServicePopulated `json:"services"`
	Tag       []string           `json:"tag"`
}

// se usuario é super admin
// // Se sim, Ok
// se usuario tem tag especifica ex: "branch-manager", "company-manager", "service-manager"
// // Se ele tem tag especifica:
// // Verificar se o ID do usuario está salvo no 'permissions' do recurso e do metodo utilizado.
// // ex: PATCH branch/5
// // permissions : { "PATCH": [1, 2, 3, 5], "DELETE": [1, 3, 6] }
// // No caso, o usuario com ID 5 pode fazer PATCH na branch 5.
// // Se sim, Ok
// se não, retornar HTTP 403.

// tag: ["super-admin", "branch-manager"]

type UserPopulated struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}
