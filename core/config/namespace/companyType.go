package namespace

type companyTypeStruct struct {
	InterfaceKey ContextKey
	ChangesKey   ContextKey
}

var CompanyType = companyTypeStruct{
	InterfaceKey: "companyType",
	ChangesKey:   "companyType_changes",
}
