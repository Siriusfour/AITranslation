package Config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

// InitLogger

func InitLogger() *zap.SugaredLogger {
	LogMode := zapcore.DebugLevel
	if !viper.GetBool("Mode.develop") {
		LogMode = zapcore.InfoLevel
	}
	Core := zapcore.NewCore(
		getEncoder(),
		getWriteSyncer(),
		LogMode,
	)
	return zap.New(Core).Sugar()
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
	fmt.Println(stLogFilePath)

	lumberJack := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("log.MaxSize"),
		MaxBackups: viper.GetInt("log.MaxBackups"),
		MaxAge:     viper.GetInt("log.MaxAge"),
		Compress:   false,
	}
	return zapcore.AddSync(lumberJack)
}
