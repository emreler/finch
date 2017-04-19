package logger

// InfoErrorLogger is the interface for logger.
type InfoErrorLogger interface {
	Info(interface{})
	Error(err error)
}
