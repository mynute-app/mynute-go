package namespace

type TypeStruct struct {
	Model        string
	Changes      string
	Dto          string
	Associations string
	QueryId      string
	CompanyId    string
	BaseURL      string
}

var GeneralKey = TypeStruct{
	Model:        "model_key",
	Changes:      "changes_key",
	Dto:          "dto_key",
	Associations: "associations_key",
	QueryId:      "id",
	CompanyId:    "companyId",
	BaseURL:      "http://localhost:3000",
}
