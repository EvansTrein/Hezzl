package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Mode     string
	LogPath  string
	LogLevel string
}

var logger *slog.Logger

type CustomHandler struct {
	handler slog.Handler
	output  io.Writer
	attrs   []slog.Attr
	mu      *sync.Mutex
}

func NewCustomHandler(output io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		handler: slog.NewTextHandler(output, opts),
		output:  output,
		mu:      &sync.Mutex{},
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var buf bytes.Buffer

	buf.Write([]byte("\n"))

	buf.Write([]byte("Apps log: "))

	buf.Write([]byte(r.Time.Format(time.Stamp) + "\n"))

	level := r.Level.String()
	switch r.Level {
	case slog.LevelInfo:
		level = "\033[32m" + level + "\033[0m" // green color
	case slog.LevelError:
		level = "\033[31m" + level + "\033[0m" // red color
	case slog.LevelDebug:
		level = "\033[34m" + level + "\033[0m" // blue color
	case slog.LevelWarn:
		level = "\033[33m" + level + "\033[0m" // yellow color
	}
	buf.Write([]byte("level--> " + level + "\n"))

	buf.Write([]byte("\033[4m" + "message--> " + r.Message + "\033[0m" + "\n")) // underlined text

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		source := "file--> " + f.File +
			"\ncode_line--> " + "\033[38;5;208m" + strconv.Itoa(f.Line) + "\033[0m" + "\n" // orange color
		buf.Write([]byte(source))
	}

	for _, attr := range h.attrs {
		if attr.Key == "operation" {
			buf.Write([]byte("\033[38;5;90m" + attr.Key + "--> " + attr.Value.String() + "\033[0m" + "\n")) // purple color
		} else {
			buf.Write([]byte(attr.Key + "--> " + attr.Value.String() + "\n"))
		}
	}

	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "error" || attr.Key == "err" {
			buf.Write([]byte("\033[31m" + attr.Key + "--> " + attr.Value.String() + "\033[0m" + "\n")) // red color
		} else {
			buf.Write([]byte(attr.Key + "--> " + attr.Value.String() + "\n"))
		}
		return true
	})

	buf.Write([]byte("\n"))

	if _, err := h.output.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
		output:  h.output,
		attrs:   attrs,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
		output:  h.output,
	}
}

func getLevel(level string) slog.Level {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	default:
		log.Println("Unknown log level using default level")
		logLevel = slog.LevelInfo
	}

	return logLevel
}

func getLogFile(logPath string) io.Writer {
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30, // Days
		Compress:   true,
	}
}

func InitLog(c LogConfig) {
	logFile := getLogFile(c.LogPath)
	level := getLevel(c.LogLevel)

	switch c.Mode {
	case "local":
		logger = slog.New(NewCustomHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level, // DEBUG
			AddSource: true,
		}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level, 
			AddSource: true,
		}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level:     level, // INFO
			AddSource: false,
		}))
	default:
		panic("unknown environment: " + c.Mode)
	}
}

func GetLogger() *slog.Logger {
	return logger
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

func EslogStr(err error) slog.Attr {
	return slog.String("error", fmt.Sprintf("%v", err))
}
