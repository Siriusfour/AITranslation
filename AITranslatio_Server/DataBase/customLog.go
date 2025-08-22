package DataBase

import (
	gormLog "gorm.io/gorm/logger"
)

type logger struct {
	gormLog.Writer
	gormLog.Config

	infoStr string
	warnStr string
	errStr  string

	traceStr     string
	traceErrStr  string
	traceWarnStr string
}

type option interface {
	apply(*logger)
}
type logWriter struct{}

func (lw logWriter) Write(s string, any ...interface{}) {}

type optionFunc func(*logger)

func (f optionFunc) apply(log *logger) { f(log) }

func createCustomGormLog(SQL_Type string, options ...option) {

	var (
		infoStr      string
		warnStr      string
		errStr       string
		traceStr     string
		traceWarnStr string
		traceErrStr  string
	)

	var logger logger

	for _, opt := range options {
	}

}
