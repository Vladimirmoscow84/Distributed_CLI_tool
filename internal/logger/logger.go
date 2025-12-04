package logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

// NopLogger - заглушка для тестов или на случай непрокидывания логера в функции
type NopLogger struct{}

func (NopLogger) Info(string)  {}
func (NopLogger) Debug(string) {}
func (NopLogger) Error(string) {}
