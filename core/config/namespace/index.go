package namespace

type GeneralStruct struct {
	Name         string
	Model        string
	ModelArr     string
	Changes      string
	Dto          string
	DtoArr       string
	Associations string
}

type RequestStruct struct {
	Body_Byte   string
	Path        string
	Auth_Token  string
	Auth_Claims string
}

var RequestKey = RequestStruct{
	Body_Byte:   "req_body_byte",
	Path:        "req_path",
	Auth_Token:  "req_auth_token",
	Auth_Claims: "req_auth_claims",
}

var GeneralKey = GeneralStruct{
	Name:         "name_key",
	Model:        "model_key",
	ModelArr:     "modelArr_key",
	Changes:      "changes_key",
	Dto:          "dto_key",
	DtoArr:       "dtoArr_key",
	Associations: "associations_key",
}

type QueryStruct struct {
	Id        string
	CompanyId string
}

var QueryKey = QueryStruct{
	Id:        "id",
	CompanyId: "companyId",
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

var ClientKey = TypeStruct{
	Name:  "client",
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
