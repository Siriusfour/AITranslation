package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/Hooks"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitLogger

func InitLogger() {

	Global.Logger = CreateZapFactory(Hooks.ZapLogHandler)
}

// 表明日志的格式
func getEncoder() zapcore.Encoder {
	//这个函数返回一个生产环境下推荐的编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	//这里设置了日志中时间字段的键为 "time"，，表示日志中时间的键名为 "time"
	encoderConfig.TimeKey = "time"
	//配置将日志级别以大写形式输出，例如 "INFO"、"ERROR"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//对日志中时间字段进行编码：自定义了一个函数，接收一个时间参数 t 和一个编码器 encoder，将格式化后的时间以字符串形式追加到编码器中。
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(time.DateTime))
	}
	//基于所配置的编码器配置，创建了一个 JSON 格式的编码器实例，该实例将应用于日志记录器
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	stRootDir, _ := os.Getwd()
	stSeparator := string(filepath.Separator)
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format("2006-01-02") + ".log"

	lumberJack := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("log.MaxSize"),
		MaxBackups: viper.GetInt("log.MaxBackups"),
		MaxAge:     viper.GetInt("log.MaxAge"),
		Compress:   false,
	}
	return zapcore.AddSync(lumberJack)
}

func CreateZapFactory(entry func(zapcore.Entry) error) *zap.Logger {

	// 获取程序所处的模式：  开发调试 、 生产
	//variable.ConfigYml := yml_config.CreateYamlFactory()
	appDebug := Global.Config.GetBool("Mode.Develop")

	// 判断程序当前所处的模式，调试模式直接返回一个便捷的zap日志管理器地址，所有的日志打印到控制台即可
	if appDebug == true {
		if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
			return logger
		} else {
			log.Fatal("创建zap日志包失败，详情：" + err.Error())
		}
	}

	// 以下才是 非调试（生产）模式所需要的代码
	encoderConfig := zap.NewProductionEncoderConfig()

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
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at" // 生成json格式日志的时间键字段，默认为 ts,修改以后方便日志导入到 ELK 服务器

	var encoder zapcore.Encoder
	switch Global.Config.GetString("Logs.TextFormat") {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	}

	//写入器
	fileName := Consts.BasePath + Global.Config.GetString("Logs.GoSkeletonLogName")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,                                //日志文件的位置
		MaxSize:    Global.Config.GetInt("Logs.MaxSize"),    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: Global.Config.GetInt("Logs.MaxBackups"), //保留旧文件的最大个数
		MaxAge:     Global.Config.GetInt("Logs.MaxAge"),     //保留旧文件的最大天数
		Compress:   Global.Config.GetBool("Logs.Compress"),  //是否压缩/归档旧文件
	}
	writer := zapcore.AddSync(lumberJackLogger)
	// 开始初始化zap日志核心参数，
	//参数一：编码器
	//参数二：写入器
	//参数三：参数级别，debug级别支持后续调用的所有函数写日志，如果是 fatal 高级别，则级别>=fatal 才可以写日志
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry), zap.AddStacktrace(zap.WarnLevel))
}
