package logger

import (
	"effectiveMobile/config"
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

type Level int

const (
	DebugLevel Level = -4
	InfoLevel  Level = 0
	WarnLevel  Level = 4
	ErrorLevel Level = 8
	PanicLevel Level = 12
	FatalLevel Level = 16
	NoLevel    Level = 20
	Disabled   Level = 24
)

type Logger interface {
	InitLogger() error
	Debug(msg string)
	Debugf(template string, args ...interface{})
	Info(msg string)
	Infof(template string, args ...interface{})
	Warn(msg string)
	Warnf(template string, args ...interface{})
	Error(error error)
	Errorf(template string, args ...interface{})
	ErrorFull(err error)
	Fatal(msg string)
	Fatalf(template string, args ...interface{})
	Panic(msg string)
	Panicf(template string, args ...interface{})
}

// Logger
type ApiLogger struct {
	cfg    *config.Config
	logger *slog.Logger
}

// App Logger constructor
func NewApiLogger(cfg *config.Config) *ApiLogger {
	return &ApiLogger{cfg: cfg}
}

func (a *ApiLogger) InitLogger() error {
	stdoutHandler := slog.NewTextHandler(os.Stdout, nil)
	a.logger = slog.New(stdoutHandler)
	return nil
}

func (a *ApiLogger) Debug(msg string) {
	go a.logger.Debug(msg)
}

func (a *ApiLogger) Debugf(template string, args ...interface{}) {
	go a.logger.Debug(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) Info(msg string) {
	go a.logger.Info(msg)
}

func (a *ApiLogger) Infof(template string, args ...interface{}) {
	go a.logger.Info(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) Warn(msg string) {
	go a.logger.Warn(msg)
}

func (a *ApiLogger) Warnf(template string, args ...interface{}) {
	go a.logger.Warn(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) Error(err error) {
	go a.logger.Error(err.Error())
}

func (a *ApiLogger) Errorf(template string, args ...interface{}) {
	go a.logger.Error(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) Panic(msg string) {
	go a.logger.Error(msg)
}

func (a *ApiLogger) Panicf(template string, args ...interface{}) {
	go a.logger.Error(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) Fatal(msg string) {
	go a.logger.Error(msg)
}

func (a *ApiLogger) Fatalf(template string, args ...interface{}) {
	go a.logger.Error(fmt.Sprintf(template, args...))
}

func (a *ApiLogger) ErrorFull(err error) {
	pc, _, line, _ := runtime.Caller(1)
	det := runtime.FuncForPC(pc)
	msg := fmt.Sprintf("ERROR:\n%s :: %d :: %s", det.Name(), line, err.Error())
	go a.logger.Error(msg)
}
