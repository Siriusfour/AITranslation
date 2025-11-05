package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/Hooks"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

// InitLogger

func InitLogger() {
	Global.Logger = CreateZapFactory(Hooks.ZapLogHandler)
}

func CreateZapFactory(entry func(zapcore.Entry) error) map[string]*zap.Logger {

	// 获取程序所处的模式：开发调试 或 生产
	appDebug := Global.Config.GetBool("Mode.Develop")

	// 返回不同组件的日志记录器
	loggers := make(map[string]*zap.Logger)

	// 判断程序当前所处的模式，调试模式直接返回一个便捷的zap日志管理器地址
	if appDebug {
		// 如果是开发模式，直接创建一个开发模式的logger
		logger, err := zap.NewDevelopment(zap.Hooks(entry))
		if err != nil {
			log.Fatal("创建zap日志包失败，详情：" + err.Error())
		}
		loggers["default"] = logger
		return loggers
	}

	// 以下是生产模式的代码
	encoderConfig := zap.NewProductionEncoderConfig()

	// 获取日志时间精度配置
	timePrecision := Global.Config.GetString("Logs.TimePrecision")
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"
	}

	// 设置时间格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at"

	// 根据配置选择输出格式
	var encoder zapcore.Encoder
	switch Global.Config.GetString("Logs.TextFormat") {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 默认控制台模式
	}

	// 准备日志文件路径和切割策略

	logPaths := map[string]string{
		"business": Consts.BasePath + Global.Config.GetString("Logs.BusinessPath"),
		"db":       Consts.BasePath + Global.Config.GetString("Logs.DbPath"),
		"mq":       Consts.BasePath + Global.Config.GetString("Logs.MQPath"),
	}

	// 定义不同模块的日志文件
	modules := []string{"business", "db", "mq"}

	// 为每个模块创建不同的日志文件和 logger
	for _, module := range modules {
		// 创建每个模块的日志文件路径
		fileName := fmt.Sprintf("%s_%s.log", logPaths[module], module)
		lumberJackLogger := &lumberjack.Logger{
			Filename:   fileName,                                // 日志文件位置
			MaxSize:    Global.Config.GetInt("Logs.MaxSize"),    // 最大大小（MB）
			MaxBackups: Global.Config.GetInt("Logs.MaxBackups"), // 保留的旧文件最大个数
			MaxAge:     Global.Config.GetInt("Logs.MaxAge"),     // 旧文件保留天数
			Compress:   Global.Config.GetBool("Logs.Compress"),  // 是否压缩旧文件
		}
		writer := zapcore.AddSync(lumberJackLogger)

		// 设置每个模块的日志级别
		var level zapcore.Level
		switch module {
		case "business":
			level = zap.InfoLevel // 业务日志可以是 Info 级别
		case "db":
			level = zap.WarnLevel // 数据库日志可能需要 Warn 级别
		case "mq":
			level = zap.DebugLevel // MQ 日志可能更偏向调试级别
		default:
			level = zap.InfoLevel
		}

		// 创建每个模块的 Core
		zapCore := zapcore.NewCore(encoder, writer, level)

		// 使用不同的 module 名称生成 Logger
		loggers[module] = zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry), zap.AddStacktrace(zap.WarnLevel))
	}

	// 返回多个日志记录器
	return loggers
}
