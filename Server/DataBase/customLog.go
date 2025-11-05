package DataBase

import (
	"AITranslatio/Global"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type logger struct {
	gormLog.Writer
	gormLog.Config
	infoStr      string
	warnStr      string
	errStr       string
	traceStr     string
	traceErrStr  string
	traceWarnStr string
}

type Options interface {
	apply(*logger)
}

type OptionFunc func(log *logger)

func (f OptionFunc) apply(log *logger) { f(log) }

// 定义 6 个函数修改内部变量
func SetInfoStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.infoStr = format
	})
}

func SetWarnStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.warnStr = format
	})
}

func SetErrStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.errStr = format
	})
}

func SetTraceStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceStr = format
	})
}
func SetTraceWarnStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceWarnStr = format
	})
}

func SetTraceErrStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceErrStr = format
	})
}

type logOutPut struct{}

func (l logOutPut) Printf(strFormat string, args ...interface{}) {
	logRes := fmt.Sprintf(strFormat, args...)
	logFlag := "gorm 日志:"
	detailFlag := "详情："
	if strings.HasPrefix(strFormat, "[info]") || strings.HasPrefix(strFormat, "[traceStr]") {
		Global.Logger["DB"].Info(logFlag, zap.String(detailFlag, logRes))
	} else if strings.HasPrefix(strFormat, "[error]") || strings.HasPrefix(strFormat, "[traceErr]") {
		Global.Logger["DB"].Error(logFlag, zap.String(detailFlag, logRes))
	} else if strings.HasPrefix(strFormat, "[warn]") || strings.HasPrefix(strFormat, "[traceWarn]") {
		Global.Logger["DB"].Warn(logFlag, zap.String(detailFlag, logRes))
	}
}

func createCustomGormLog(sqlType string, options ...Options) *logger {

	var (
		infoStr      string
		warnStr      string
		errStr       string
		traceStr     string
		traceWarnStr string
		traceErrStr  string
	)

	logConf := gormLog.Config{
		SlowThreshold:        time.Second * Global.Config.GetDuration("SlowThreshold"),
		LogLevel:             gormLog.Warn,
		Colorful:             true,
		ParameterizedQueries: true,
	}

	logger := &logger{
		Writer:       logOutPut{},
		Config:       logConf,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}

	for _, opt := range options {
		opt.apply(logger)
	}

	return logger
}

// LogMode log mode
func (l *logger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLog.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLog.Error && (!errors.Is(err, gormLog.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLog.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormLog.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-1", sql)
		} else {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
