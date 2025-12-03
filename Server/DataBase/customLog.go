package DataBase

import (
	"AITranslatio/Config/interf"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// DBLogger 实现 gorm.io/gorm/logger.Interface，专门输出到 l.logger["DB"]
type DBLogger struct {
	gormLog.Config
	logger *zap.Logger
	dbType string // 标记数据库类型，比如 mysql / postgres，方便区分
}

// ===== 可选配置 Option =====

type LogOption func(*DBLogger)

// 设置日志等级
func WithLogLevel(level gormLog.LogLevel) LogOption {
	return func(l *DBLogger) {
		l.LogLevel = level
	}
}

// 设置慢 SQL 阈值
func WithSlowThreshold(d time.Duration) LogOption {
	return func(l *DBLogger) {
		l.SlowThreshold = d
	}
}

// 是否忽略 RecordNotFound 错误（不打成 ERROR）
func WithIgnoreRecordNotFound(ignore bool) LogOption {
	return func(l *DBLogger) {
		l.IgnoreRecordNotFoundError = ignore
	}
}

// ===== 构造函数 =====

// NewDBLogger 创建一个工业级的 GORM Logger
func NewDBLogger(cfg interf.ConfigInterface, logger *zap.Logger, sqlType string, opts ...LogOption) gormLog.Interface {
	// 默认配置
	l := &DBLogger{
		Config: gormLog.Config{
			SlowThreshold:        cfg.GetDuration("SlowThreshold"),
			LogLevel:             gormLog.Warn, // 默认只打 Warn 以上的日志
			Colorful:             false,        // 我们用 zap 的格式，不需要 gorm 加颜色
			ParameterizedQueries: true,
			// 默认忽略 record not found，不打成 ERROR（比较接近业务习惯）
			IgnoreRecordNotFoundError: true,
		},
		dbType: sqlType,
		logger: logger,
	}

	// 应用外部配置
	for _, opt := range opts {
		opt(l)
	}

	return l
}

// ===== 实现 gormLog.Interface =====

// LogMode 返回一个设置了新日志等级的 logger 拷贝
func (l *DBLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

// Info 普通信息日志
func (l *DBLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel < gormLog.Info {
		return
	}

	l.logger.Info("gorm-info",
		zap.String("db", l.dbType),
		zap.String("file", utils.FileWithLineNum()),
		zap.String("msg", msg),
		zap.Any("data", data),
	)
}

// Warn 警告日志
func (l *DBLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel < gormLog.Warn {
		return
	}

	l.logger.Warn("gorm-warn",
		zap.String("db", l.dbType),
		zap.String("file", utils.FileWithLineNum()),
		zap.String("msg", msg),
		zap.Any("data", data),
	)
}

// Error 错误日志
func (l *DBLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel < gormLog.Error {
		return
	}

	l.logger.Error("gorm-error",
		zap.String("db", l.dbType),
		zap.String("file", utils.FileWithLineNum()),
		zap.String("msg", msg),
		zap.Any("data", data),
	)
}

// Trace SQL 相关日志（核心）
func (l *DBLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == gormLog.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1e6 // 毫秒
	sql, rows := fc()
	file := utils.FileWithLineNum()

	// 规范一下 rows：-1 表示未知
	if rows == -1 {
		rows = -1
	}

	// 1. 如果有错误（且不是被忽略的 RecordNotFound），打 ERROR
	if err != nil &&
		l.LogLevel >= gormLog.Error &&
		(!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {

		l.logger.Error("gorm-sql-error",
			zap.String("db", l.dbType),
			zap.String("file", file),
			zap.Error(err),
			zap.Float64("elapsed_ms", elapsedMs),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
		return
	}

	// 2. 慢 SQL，打 WARN
	if l.SlowThreshold > 0 &&
		elapsed > l.SlowThreshold &&
		l.LogLevel >= gormLog.Warn {

		l.logger.Warn("gorm-slow-sql",
			zap.String("db", l.dbType),
			zap.String("file", file),
			zap.Float64("elapsed_ms", elapsedMs),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.Duration("slow_threshold", l.SlowThreshold),
		)
		return
	}

	// 3. 普通 SQL Trace，打 INFO
	if l.LogLevel >= gormLog.Info {
		l.logger.Info("gorm-sql",
			zap.String("db", l.dbType),
			zap.String("file", file),
			zap.Float64("elapsed_ms", elapsedMs),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
