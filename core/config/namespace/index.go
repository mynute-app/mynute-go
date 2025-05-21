package namespace

type GeneralStruct struct {
	Name            string
	Model           string
	ModelArr        string
	Changes         string
	Dto             string
	DtoArr          string
	Associations    string
	DatabaseSession string
	Company         string
	CompanySchema   string
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
	Name:            "name_key",
	Model:           "model_key",
	ModelArr:        "modelArr_key",
	Changes:         "changes_key",
	Dto:             "dto_key",
	DtoArr:          "dtoArr_key",
	Associations:    "associations_key",
	DatabaseSession: "db_session_key",
	CompanySchema:   "company_schema_key",
	Company:         "company_key",
}

type QueryStruct struct {
	Id        string
	CompanyId string
}

type RouteParams struct {
	CompanyID string
}

var RouteParamsKey = RouteParams{
	CompanyID: "company_id",
}

var QueryKey = QueryStruct{
	Id:        "id",
	CompanyId: "companyId",
}

type HeadersStruct struct {
	Company string
	Auth    string
}

var HeadersKey = HeadersStruct{
	Company: "X-Company-ID",
	Auth:    "X-Auth-Token",
}

type TypeStruct struct {
	Name  string
	Model string
}

var CompanyKey = TypeStruct{
	Name:  "company",
	Model: "company_model",
}

var SectorKey = TypeStruct{
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

var Role = struct {
	Owner            string
	GeneralManager   string
	BranchManager    string
	BranchSupervisor string
	Employee         string
}{
	Owner:            "Owner",
	GeneralManager:   "General Manager",
	BranchManager:    "Branch Manager",
	BranchSupervisor: "Branch Supervisor",
	Employee:         "Employee",
}

var (
	CreateActionMethod = "POST"
	ViewActionMethod   = "GET"
	UpdateActionMethod = "PATCH"
	DeleteActionMethod = "DELETE"
)

var UploadsFolder = "my_files"
var StaticServerFolder = "static"
