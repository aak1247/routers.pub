package infra

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"routers.pub/utils"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	env "routers.pub/env"
)

// Log 全局日志变量
var Log *zap.SugaredLogger

// InitLogger
// 初始化日志
// filename 日志文件路径
// level 日志级别
// maxSize 每个日志文件保存的最大尺寸 单位：M
// maxBackups 日志文件最多保存多少个备份
// maxAge 文件最多保存多少天
// compress 是否压缩
// serviceName 服务名
// 由于zap不具备日志切割功能, 这里使用lumberjack配合
func InitLogger() {
	now := time.Now()
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", env.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", env.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	var coreArr []zapcore.Core

	// 获取编码器
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 指定时间格式
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // ，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	////encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        // 显示完整文件路径
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder, //zapcore.CapitalColorLevelEncoder, // 有颜色输出到控制台
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		//EncodeTime: zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		//EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendInt64(int64(d) / 1000000)
		//},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeCaller: zapcore.FullCallerEncoder,
		//EncodeName:       nil,
		//ConsoleSeparator: "",
	}

	// encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 设置为没有颜色 输出到文件
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	noColorEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	// 当yml配置中的等级大于Error时，lowPriority级别日志停止记录
	if env.Conf.Logs.Level >= 2 {
		lowPriority = func(level zapcore.Level) bool {
			return false
		}
	}

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,          //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    env.Conf.Logs.MaxSize,    //文件大小限制,单位MB
		MaxAge:     env.Conf.Logs.MaxAge,     //日志文件保留天数
		MaxBackups: env.Conf.Logs.MaxBackups, //最大保留日志文件数量
		LocalTime:  false,
		Compress:   env.Conf.Logs.Compress, //是否压缩处理
	})

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,         //日志文件存放目录
		MaxSize:    env.Conf.Logs.MaxSize,    //文件大小限制,单位MB
		MaxAge:     env.Conf.Logs.MaxAge,     //日志文件保留天数
		MaxBackups: env.Conf.Logs.MaxBackups, //最大保留日志文件数量
		LocalTime:  false,
		Compress:   env.Conf.Logs.Compress, //是否压缩处理
	})
	// 文件输出
	infoFileCore := zapcore.NewCore(noColorEncoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer), lowPriority)
	errorFileCore := zapcore.NewCore(noColorEncoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer), highPriority)

	// 控制台输出
	infoOutCore := zapcore.NewCore(noColorEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), lowPriority)
	errorOutCore := zapcore.NewCore(noColorEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore, errorFileCore, errorOutCore, infoOutCore)

	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	Log = logger.Sugar()
	Log.Info("初始化zap日志完成!")
}

func AlertError(err error) {
	Log.Errorf("[ERROR-ALERT] api panic: %v\n Stack:%s", err, utils.CallStack(20, 1))
}

func AlertMessage(msg string) {
	Log.Errorf("[ERROR-ALERT] api panic: %v\n Stack:%s", msg, utils.CallStack(20, 1))
}
