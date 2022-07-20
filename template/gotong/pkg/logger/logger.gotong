package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjackv2 "gopkg.in/natefinch/lumberjack.v2"
)

const (
	Console = "console"
	File    = "file"
)

var (
	once  = &sync.Once{}
	Level = map[int8]zapcore.Level{
		-1: zapcore.DebugLevel,
		0:  zapcore.InfoLevel,
		1:  zapcore.WarnLevel,
		2:  zapcore.ErrorLevel,
		3:  zapcore.DPanicLevel,
		4:  zapcore.PanicLevel,
		5:  zapcore.FatalLevel,
	}
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

type Option struct {
	Target     string // option: console | file
	Level      int8   // option: -1 ... 5
	Filename   string // location where the file is stored. Only if Target is `file`
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // CapitalColorLevelEncoder | CapitalLevelEncoder
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // short | full
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Init(opt Option) {
	w := zapcore.AddSync(&lumberjackv2.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize, // megabytes
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge, // days
	})

	var writeSyncer zapcore.WriteSyncer
	if opt.Target == Console {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	if opt.Target == File {
		writeSyncer = zapcore.NewMultiWriteSyncer(w)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writeSyncer,
		Level[opt.Level],
	)

	once.Do(func() {
		Logger = zap.New(core, zap.AddCaller())
		Sugar = Logger.Sugar()
	})
}

func DefaultOption() Option {
	return Option{
		Target:     "console",
		Level:      int8(zapcore.InfoLevel),
		Filename:   "log/tong.log",
		MaxSize:    1024,
		MaxBackups: 10,
		MaxAge:     7,
	}
}
