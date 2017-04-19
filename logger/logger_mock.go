package logger

// MockLogger is the mock that satisfies InfoErrorLogger interface.
type MockLogger struct{}

// Info .
func (m *MockLogger) Info(data interface{}) {}

// Error .
func (m *MockLogger) Error(error) {}
