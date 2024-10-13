package domain

type Logger interface {
	Infof(format string, args ...interface{})
	Info(args ...any)
	Errorf(format string, args ...interface{})
	Error(args ...any)
	Debugf(format string, args ...interface{})
	Debug(args ...any)
	Warnf(format string, args ...interface{})
	Warn(args ...any)
	Fatalf(format string, args ...interface{})
	Fatal(args ...any)
	Panicf(format string, args ...interface{})
	Panic(args ...any)
}
