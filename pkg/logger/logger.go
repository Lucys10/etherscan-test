package logger

import "github.com/sirupsen/logrus"

type Log struct {
	*logrus.Logger
}

func NewLogger(logLvl logrus.Level) *Log {
	l := &Log{logrus.New()}
	l.SetLevel(logLvl)
	return l
}
