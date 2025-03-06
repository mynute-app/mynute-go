package namespace

type ErrorStruct struct {
	description string
	http_status int
}

type ErrorTypes struct {
	InterfaceDataNotFound ErrorStruct
	InvalidLogin          ErrorStruct
}

var MyErrors = ErrorTypes {
	InterfaceDataNotFound: ErrorStruct{

	}
}
