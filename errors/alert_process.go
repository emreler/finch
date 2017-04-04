package errors

// RetryProcessError is an error type which implies that process should be retried
type RetryProcessError struct {
	Msg string
}

func (e RetryProcessError) Error() string {
	return e.Msg
}

// InvalidAlertIDError is an error which is used when alertID is invalid
type InvalidAlertIDError struct {
	Msg string
}

func (e *InvalidAlertIDError) Error() string {
	return e.Msg
}
