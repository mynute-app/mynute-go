package namespace

type ContextKey string

type TypeStruct struct {
	Model        ContextKey
	Changes      ContextKey
	Dto          ContextKey
	Associations ContextKey
	QueryId      ContextKey
}

var GeneralKey = TypeStruct{
	Model:        "model_key",
	Changes:      "changes_key",
	Dto:          "dto_key",
	Associations: "associations_key",
	QueryId:      "id",
}
