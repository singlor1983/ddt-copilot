package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type LogContext = zerolog.Context

type LogOptions func(LogContext) LogContext

type CustomLogger struct {
	Logger zerolog.Logger

	fw *fileWriter
}

type fileWriter struct {
	logRoot             string
	f                   *os.File
	w                   *bufio.Writer
	mutex               sync.Mutex
	lastUpdateHourInDay int64
}

func (self *fileWriter) Flush() error {
	var err error
	self.mutex.Lock()
	if self.w != nil {
		err = self.w.Flush()
	}
	self.mutex.Unlock()
	return err
}

func (self *fileWriter) Write(d []byte) (int, error) {
	self.mutex.Lock()
	t := time.Now()
	err := self.prepare(&t)
	var n int
	if err == nil {
		n, err = self.w.Write(d)
	}
	self.mutex.Unlock()
	return n, err
}

func (self *fileWriter) prepare(now *time.Time) error {
	hoursInDay := now.Unix() / 3600
	if self.f == nil || hoursInDay != self.lastUpdateHourInDay { // 15:04:05
		fileName := fmt.Sprintf("%s/%s.log", self.logRoot, now.Format("2006010215"))
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600|0644)
		if err != nil {
			return err
		}
		if self.f != nil {
			_ = self.w.Flush()
			self.f.Close()
		}
		self.f = f
		self.w = bufio.NewWriter(f)
		self.lastUpdateHourInDay = hoursInDay
	}
	return nil
}

func newConsoleWriter() *zerolog.ConsoleWriter {
	cw := &zerolog.ConsoleWriter{Out: os.Stdout}
	cw.TimeFormat = "15:04:05" // 跟unity时间格式保持一致
	cw.FormatTimestamp = func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			t = tt
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixMs:
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				case zerolog.TimeFormatUnixMicro:
					nsec = int64(time.Duration(i) * time.Microsecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec).Local()
				t = ts.Format(cw.TimeFormat)
			}
		}
		if cw.NoColor {
			return t
		}
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 90, t)
	}
	return cw
}

func NewLogger(enableConsole bool, root string, opts ...LogOptions) *CustomLogger {
	var w io.Writer
	fw := &fileWriter{logRoot: root}
	ml := &CustomLogger{fw: fw}
	if enableConsole {
		w = zerolog.MultiLevelWriter(fw, newConsoleWriter())
	} else {
		w = fw
	}
	ctx := zerolog.New(w).With().Timestamp()
	for _, v := range opts {
		ctx = v(ctx)
	}
	ml.Init(ctx.Logger())
	return ml
}

func NewConsoleLogger(opts ...LogOptions) *CustomLogger {
	var w io.Writer
	ml := &CustomLogger{}
	w = zerolog.MultiLevelWriter(newConsoleWriter())
	ctx := zerolog.New(w).With().Timestamp()
	for _, v := range opts {
		ctx = v(ctx)
	}
	ml.Init(ctx.Logger())
	return ml
}

// Flush 程序他退出之前记得清空缓冲区，不然可能丢失日志
func Flush() {
	flushAllLog()
}

func (ml *CustomLogger) Init(logger zerolog.Logger) {
	ml.Logger = logger
	addFlusher(ml.fw)
}

func (ml *CustomLogger) NewLogger(opts ...LogOptions) *CustomLogger {
	ctx := ml.Logger.With()
	for _, v := range opts {
		ctx = v(ctx)
	}
	return &CustomLogger{fw: ml.fw, Logger: ctx.Logger()}
}

func (ml *CustomLogger) Trace() *LogEvent {
	return &LogEvent{Event: ml.Logger.Trace()}
}

func (ml *CustomLogger) Debug() *LogEvent {
	return &LogEvent{Event: ml.Logger.Debug()}
}

func (ml *CustomLogger) Info() *LogEvent {
	return &LogEvent{Event: ml.Logger.Info()}
}

func (ml *CustomLogger) Warn() *LogEvent {
	return &LogEvent{Event: ml.Logger.Warn()}
}

func (ml *CustomLogger) Error() *LogEvent {
	return &LogEvent{Event: ml.Logger.Error()}
}

func (ml *CustomLogger) Fatal() *LogEvent {
	return &LogEvent{Event: ml.Logger.Fatal()}
}

func (ml *CustomLogger) Panic() *LogEvent {
	return &LogEvent{Event: ml.Logger.Panic()}
}

func (ml *CustomLogger) Log() *LogEvent {
	return &LogEvent{Event: ml.Logger.Log()}
}
