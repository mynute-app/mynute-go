package namespace

type TypeStruct struct {
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
