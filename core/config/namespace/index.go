package namespace

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

var QueryKey = QueryStruct{
	Id:        "id",
	CompanyId: "companyId",
	BaseURL:   "http://localhost:3000",
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
	Name:  "company_type",
	Model: "company_type_model",
}

var BranchKey = TypeStruct{
	Name:  "branch",
	Model: "branch_model",
}

var EmployeeKey = TypeStruct{
	Name:  "employee",
	Model: "employee_model",
}

var ServiceKey = TypeStruct{
	Name:  "service",
	Model: "service_model",
}

var HolidaysKey = TypeStruct{
	Name:  "holidays",
	Model: "holidays_model",
}