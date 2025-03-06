package namespace

type ErrorStruct struct {
	description_en string
	description_br string
	id int
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
