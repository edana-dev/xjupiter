package log

import (
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap/zapcore"
	"time"
)

func init() {
	logConfig := xlog.Config{
		Debug: true,
		Async: false,
		EncoderConfig: &zapcore.EncoderConfig{
			TimeKey:       "ts",
			LevelKey:      "lv",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stack",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	xlog.DefaultLogger = logConfig.Build()
	xlog.JupiterLogger = logConfig.Build()
}
