package logger

type InfoErrorLogger interface {
	Info(interface{})
	Error(err error)
}
