package logger

type MockLogger struct{}

func (m *MockLogger) Info(data interface{}) {}
func (m *MockLogger) Error(error)           {}
