package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/imperfectgo/zap-syslog"
	"github.com/imperfectgo/zap-syslog/syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(path string) {
	// file logger
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	fileWriteSyncer := zapcore.AddSync(file)
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	// console logger
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
    developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	// system logger
	syslogCfg := zapsyslog.SyslogEncoderConfig{
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		Facility: syslog.LOG_LOCAL0,
		Hostname: "localhost",
		PID:	  os.Getpid(),
		App: 	  "syslog-test",
	}
	syslogEncoder := zapsyslog.NewSyslogEncoder(syslogCfg)
	sync, err := zapsyslog.NewConnSyncer("tcp", "ed.de:514")
	if err != nil {
		fmt.Println(err)
	}
	// set core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, fileWriteSyncer, zapcore.DebugLevel),
		zapcore.NewCore(syslogEncoder, zapcore.Lock(sync), zapcore.DebugLevel),
	)
	Log = zap.New(core)
}

func Info(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Warn(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

func main() {
	InitLogger("test2.log")
	Info("ahahahahaha")
	Error("test")
}
