package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Logger *zap.Logger

func InitLogger(dc string, workId string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	//指定时间格式
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("20060102-15:04:05.000"))
	}
	encoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.InfoLevel:
			fallthrough
		case zapcore.WarnLevel:
			encoder.AppendString(strings.ToUpper(level.String()) + " |" + dc + "||" + workId + "|||" + strconv.Itoa(os.Getpid()) + "|" + strconv.Itoa(GoID()))
		default:
			encoder.AppendString(strings.ToUpper(level.String()) + "|" + dc + "||" + workId + "|||" + strconv.Itoa(os.Getpid()) + "|" + strconv.Itoa(GoID()))
		}
	}
	//显示完整文件路径
	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(caller.TrimmedPath())
	}
	encoderConfig.ConsoleSeparator = "|"
	//获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app-NceXdsServer-0300099.log", //日志文件存放目录
		MaxSize:    1,                                     //文件大小限制,单位MB
		MaxBackups: 5,                                     //最大保留日志文件数量
		MaxAge:     30,                                    //日志文件保留天数
		Compress:   true})                                 //是否压缩处理
	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	//AddCaller()为显示文件名和行号
	Logger = zap.New(core, zap.AddCaller())
	Logger.Info("日志配置初始化完成")
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
