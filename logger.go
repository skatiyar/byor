package byor

import (
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

func newLogger(prefix string) logger {
	return logger{log.New(os.Stdout, prefix, log.Lmsgprefix)}
}

func (l logger) errorf(format string, v ...any) {
	l.Printf(" [ERROR] "+format+"\n", v...)
}

func (l logger) debugf(format string, v ...any) {
	l.Printf(" [DEBUG] "+format+"\n", v...)
}
