package logging

import (
	"fmt"
	"os"
)

type stdIoLogger struct {
}

func MakeStdIoLogger() Logger {
	logger := &stdIoLogger{}
	return logger
}

func (l stdIoLogger) Sync() {
	err := os.Stdout.Sync()
	if err != nil {
		panic(err)
	}

	err = os.Stderr.Sync()
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Fatal(args ...interface{}) {
	_, err := fmt.Fprint(os.Stderr, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Fatalf(template string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, template, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Warnf(template string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stdout, template, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Warn(args ...interface{}) {
	_, err := fmt.Fprint(os.Stdout, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Debugf(template string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stdout, template, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Debug(args ...interface{}) {
	_, err := fmt.Fprint(os.Stdout, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Errorf(template string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, template, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Error(args ...interface{}) {
	_, err := fmt.Fprint(os.Stderr, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Infof(template string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stdout, template, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Info(args ...interface{}) {
	_, err := fmt.Fprint(os.Stdout, args...)
	if err != nil {
		panic(err)
	}
}

func (l stdIoLogger) Panic(args ...interface{}) {
	_, err := fmt.Fprint(os.Stderr, args...)
	if err != nil {
		panic(err)
	}

	panic(args)
}
