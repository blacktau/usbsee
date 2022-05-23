package logging

type Logger interface {
	Sync()
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Warn(args ...interface{})
	Debugf(template string, args ...interface{})
	Debug(args ...interface{})
	Errorf(template string, args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
}
