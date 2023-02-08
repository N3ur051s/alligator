package log

import (
	"alligator/pkg/config"
	"fmt"
	"io"
	"os"

	"github.com/ehazlett/simplelog"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	l = logrus.New()
)

func Init(c *cli.Context, cfg config.Options) {
	switch c.String("log-format") {
	case "simple":
		l.SetFormatter(&simplelog.StandardFormatter{})
	case "text":
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	}
	l.SetOutput(os.Stdout)
	if cfg.Debug {
		l.SetLevel(logrus.DebugLevel)
		l.Debugf("Loglevel set to [%v]", logrus.DebugLevel)
	}

	file, err := os.OpenFile(cfg.AuditLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		l.SetOutput(fileAndStdoutWriter)
	} else {
		l.Errorf("failed to log to file: [%s]", err)
	}
}

func Debug(i ...interface{}) {
	l.Debug(i)
}

func Debugf(format string, values ...interface{}) {
	l.Debugf(format, values...)
}

func Info(i ...interface{}) {
	l.Info(i)
}

func Infof(format string, values ...interface{}) {
	l.Infof(format, values...)
}

func Warn(i ...interface{}) {
	l.Warn(i)
}

func Warnf(format string, values ...interface{}) {
	l.Warnf(format, values...)
}

func Error(i ...interface{}) {
	l.Error(i)
}

func Errorf(format string, values ...interface{}) {
	l.Errorf(format, values...)
}

func Fatal(i ...interface{}) {
	l.Fatal(i)
}

func Fatalf(format string, values ...interface{}) {
	l.Fatalf(format, values...)
}

func Panic(i ...interface{}) {
	l.Panic(i)
}

func Panicf(format string, args ...interface{}) {
	l.Panicf(format, args)
}

type OrmWriter struct {
	ormlog *logrus.Logger
}

func (m *OrmWriter) Printf(format string, v ...interface{}) {
	logstr := fmt.Sprintf(format, v...)
	m.ormlog.Error(logstr)
}

func NewOrmWriter() *OrmWriter {
	return &OrmWriter{ormlog: l}
}
