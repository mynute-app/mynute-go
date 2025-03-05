package namespace

import "fmt"

type GeneralStruct struct {
	Name         string
	Model        string
	ModelArr     string
	Changes      string
	Dto          string
	DtoArr       string
	Associations string
	UserData     string
}

type RequestStruct struct {
	Body      string
	Path      string
	AuthToken string
}

var RequestKey = RequestStruct{
	Body:      "req_body",
	Path:      "req_path",
	AuthToken: "req_auth_token",
}

var GeneralKey = GeneralStruct{
	Name:         "name_key",
	Model:        "model_key",
	ModelArr:     "modelArr_key",
	Changes:      "changes_key",
	Dto:          "dto_key",
	DtoArr:       "dtoArr_key",
	Associations: "associations_key",
	UserData:     "user_data",
}

type QueryStruct struct {
	Id        string
	CompanyId string
	BaseURL   string
}

var AppPort = "4000"

var QueryKey = QueryStruct{
	Id:        "id",
	CompanyId: "companyId",
	BaseURL:   fmt.Sprintf("http://localhost:%s", AppPort),
}

type TypeStruct struct {
	Name  string
	Model string
}

var CompanyKey = TypeStruct{
	Name:  "company",
	Model: "company_model",
}

var CompanyTypeKey = TypeStruct{
	Name:  "sector",
	Model: "sector_model",
}

var BranchKey = TypeStruct{
	Name:  "branch",
	Model: "branch_model",
}

var UserKey = TypeStruct{
	Name:  "user",
	Model: "user_model",
}

var ServiceKey = TypeStruct{
	Name:  "service",
	Model: "service_model",
}

var HolidaysKey = TypeStruct{
	Name:  "holidays",
	Model: "holidays_model",
}
