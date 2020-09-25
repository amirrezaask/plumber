package plumber

type Logger interface {
	Debug(format string, args ...string)
	Warn(format string, args ...string)
	Error(format string, args ...string)
	Fatal(format string, args ...string)
}

func NewLoggerFromZap() Logger {
	return nil
}
