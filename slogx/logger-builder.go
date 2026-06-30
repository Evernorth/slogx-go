package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Format int

const (
	FormatText Format = 0
	FormatJSON Format = 1
)

// timeLayoutTokens are the reference substrings that Go's time package recognises as format
// tokens. A layout must contain at least one to produce meaningful timestamp output.
// See https://pkg.go.dev/time#Layout for the full reference time.
var timeLayoutTokens = []string{
	"2006", "06", // year
	"January", "Jan", "01", // month
	"Monday", "Mon", "_2", "02", // day
	"15", "03", // hour (24h / 12h zero-padded)
	"04",                            // minute (zero-padded)
	"05",                            // second (zero-padded)
	".000000000", ".000000", ".000", // sub-second (fixed width)
	".999999999", ".999999", ".999", // sub-second (trailing zeros trimmed)
	"PM", "pm", // AM/PM marker
	"MST", "-0700", "-07:00", "-07", // timezone name / numeric offset
	"Z0700", "Z07:00", // UTC indicator + offset
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
	lvlPtr, err := GetLevelByName(level)
	if err != nil {
		panic(fmt.Sprintf("invalid log level, %s is not a valid log level from the slog package", level))
	}
	lb.level = *lvlPtr
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

// WithTimestampFormat sets the timestamp format of the logs. The format must contain at least
// one recognised Go time layout token (e.g. "2006", "Jan", "15") so that it produces a
// meaningful timestamp. See https://pkg.go.dev/time#Layout for the reference time.
func (lb *defaultLoggerBuilder) WithTimestampFormat(format string) LoggerBuilder {
	for _, token := range timeLayoutTokens {
		if strings.Contains(format, token) {
			lb.timestampFormat = format
			return lb
		}
	}
	panic(fmt.Sprintf("invalid timestamp format: %q contains no recognised Go time layout tokens", format))
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
	}

	// Only install ReplaceAttr when a custom timestamp format is configured. This keeps
	// the per-attribute callback off the hot path entirely for the default format.
	if lb.timestampFormat != "" {
		format := lb.timestampFormat
		handlerOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(slog.TimeKey, a.Value.Time().Format(format))
			}
			return a
		}
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
