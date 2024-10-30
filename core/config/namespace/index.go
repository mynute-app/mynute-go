package namespace

type TypeStruct struct {
	Name         string
	Model        string
	ModelArr     string
	Changes      string
	Dto          string
	DtoArr       string
	Associations string
	QueryId      string
	CompanyId    string
	BaseURL      string
}

var GeneralKey = TypeStruct{
	Name:         "name_key",
	Model:        "model_key",
	ModelArr:     "modelArr_key",
	Changes:      "changes_key",
	Dto:          "dto_key",
	DtoArr:       "dtoArr_key",
	Associations: "associations_key",
	QueryId:      "id",
	CompanyId:    "companyId",
	BaseURL:      "http://localhost:3000",
}

var CompanyKey = TypeStruct{
	Name: "company",
	Model: "company_model",
}

var CompanyTypeKey = TypeStruct{
	Name: "company_type",
}

var BranchKey = TypeStruct{
	Name: "branch",
}

var EmployeeKey = TypeStruct{
	Name: "employee",
}

var ServiceKey = TypeStruct{
	Name: "service",
}
