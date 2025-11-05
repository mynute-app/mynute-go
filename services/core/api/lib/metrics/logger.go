package myLogger

type Logger interface {
	Log(message string, labels map[string]string) error
}

// func New(log_type string) Logger {
// 	switch log_type {
// 	case "loki":
// 		return &Loki{}
// 	case "file":
// 		return &File{}
// 	default:
// 		return nil
// 	}
// }

