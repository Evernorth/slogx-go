package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Format int

const (
	FormatText Format = 0
	FormatJSON Format = 1
)

// lb.timestampFormat must be one of the standard time constants
// https://pkg.go.dev/time#pkg-constants
var validTimestampFormats = map[string]struct{}{
	time.Layout:      {},
	time.ANSIC:       {},
	time.UnixDate:    {},
	time.RubyDate:    {},
	time.RFC822:      {},
	time.RFC822Z:     {},
	time.RFC850:      {},
	time.RFC1123:     {},
	time.RFC1123Z:    {},
	time.RFC3339:     {},
	time.RFC3339Nano: {},
	time.Kitchen:     {},
	time.Stamp:       {},
	time.StampMilli:  {},
	time.StampMicro:  {},
	time.StampNano:   {},
	time.DateTime:    {},
	time.DateOnly:    {},
	time.TimeOnly:    {},
}

type LoggerBuilder interface {
	WithContextHandler() LoggerBuilder
	WithFormat(format Format) LoggerBuilder
	WithWriter(writer io.Writer) LoggerBuilder
	WithLevel(level slog.Level) LoggerBuilder
	WithLevelString(level string) LoggerBuilder
	WithLevelEnvVar(key string) LoggerBuilder
	WithLevelFunc(key string, levelFunc LevelFunc) LoggerBuilder
	WithTimestampFormat(format string) LoggerBuilder
	Build() (*slog.Logger, *slog.LevelVar)
}

type defaultLoggerBuilder struct {
	writer            io.Writer
	format            Format
	level             slog.Level
	useContextHandler bool
	levelKey          string
	levelFunc         LevelFunc
	timestampFormat   string
}

// NewLoggerBuilder creates a new LoggerBuilder with default values.  The default values are:  LevelInfo, FormatText,
// useContextHandler=false, levelKey="", levelFunc=nil, writer=os.Stderr. and the slog default for timestamp format
func NewLoggerBuilder() LoggerBuilder {
	return &defaultLoggerBuilder{
		level:             slog.LevelInfo,
		format:            FormatText,
		useContextHandler: false,
		levelKey:          "",
		levelFunc:         nil,
		writer:            os.Stderr,
		timestampFormat:   "",
	}
}

// WithContextHandler enables the ContextHandler for the logger.
func (lb *defaultLoggerBuilder) WithContextHandler() LoggerBuilder {
	lb.useContextHandler = true
	return lb
}

// WithFormat sets the Format for the logger.
func (lb *defaultLoggerBuilder) WithFormat(format Format) LoggerBuilder {
	lb.format = format
	return lb
}

// WithWriter sets the io.Writer for the logger.
func (lb *defaultLoggerBuilder) WithWriter(writer io.Writer) LoggerBuilder {
	lb.writer = writer
	return lb
}

// WithLevel sets the slog.Level for the logger.
func (lb *defaultLoggerBuilder) WithLevel(level slog.Level) LoggerBuilder {
	lb.level = level
	return lb
}

// WithLevelString sets the slog.Level for the logger with a string. Defaults to INFO if the string is invalid
func (lb *defaultLoggerBuilder) WithLevelString(level string) LoggerBuilder {
	levelString := strings.ToUpper(strings.TrimSpace(level))
	var err error
	levelPtr, err := GetLevelByName(levelString)
	if err != nil {
		panic(fmt.Sprintf("invalid log level, %s is not a valid log level from the slog package", level))
	}
	lb.level = *levelPtr
	return lb
}

// WithLevelEnvVar sets the environment variable key to use to get the level.
func (lb *defaultLoggerBuilder) WithLevelEnvVar(key string) LoggerBuilder {
	lb.levelKey = key
	return lb
}

// WithLevelFunc sets the level variable key and function to use to get the level name.
func (lb *defaultLoggerBuilder) WithLevelFunc(key string, levelFunc LevelFunc) LoggerBuilder {
	lb.levelKey = key
	lb.levelFunc = levelFunc
	return lb
}

// WithTimestampFormat sets the timestamp format of the logs.
func (lb *defaultLoggerBuilder) WithTimestampFormat(format string) LoggerBuilder {

	// If format isnt a standard format set to default
	if _, ok := validTimestampFormats[format]; !ok {
		panic(fmt.Sprintf("invalid timestamp format: %q must be a valid constant from time package", format))
	}
	lb.timestampFormat = format
	return lb
}

// Build creates a new slog.Logger with the provided configuration. A slog.LevelVar to control the
// logger level is also returned.
func (lb *defaultLoggerBuilder) Build() (*slog.Logger, *slog.LevelVar) {

	// Set the default level
	levelVar := new(slog.LevelVar)
	levelVar.Set(lb.level)

	// If a level variable key is provided and a level function is provided, set the level from the level function
	// Otherwise, try to set the level from the environment
	if lb.levelKey != "" {
		if lb.levelFunc != nil {
			levelVar.Set(GetLevelFromFunc(lb.levelKey, lb.levelFunc, lb.level))
		} else {
			levelVar.Set(GetLevelFromFunc(lb.levelKey, getEnvLevelFunc(), lb.level))
		}
	}

	// Create the handler
	handlerOpts := &slog.HandlerOptions{
		Level: levelVar,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && lb.timestampFormat != "" {
				t := a.Value.Time()
				return slog.String(
					slog.TimeKey,
					t.Format(lb.timestampFormat),
				)
			}
			return a
		},
	}
	var handler slog.Handler
	if lb.format == FormatJSON {
		handler = slog.NewJSONHandler(lb.writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(lb.writer, handlerOpts)
	}

	// If the context handler is enabled, wrap the handler with a ContextHandler
	if lb.useContextHandler {
		handler = NewContextHandler(handler)
	}

	// Create the logger
	logger := slog.New(handler)

	slog.Default()

	return logger, levelVar
}
