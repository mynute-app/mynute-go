package namespace

type ContextKey string

type TypeStruct struct {
	InterfaceKey    ContextKey
	ChangesKey      ContextKey
	DtoKey          ContextKey
	AssociationsKey ContextKey
}
