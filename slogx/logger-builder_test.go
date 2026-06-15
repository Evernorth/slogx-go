package slogx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HELPERS

func expectPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but none occurred")
		}
	}()
	fn()
}

// TESTS

func TestNewLoggerBuilder(t *testing.T) {
	builder := NewLoggerBuilder().(*defaultLoggerBuilder)
	assert.Equal(t, slog.LevelInfo, builder.level)
	assert.Equal(t, FormatText, builder.format)
	assert.False(t, builder.useContextHandler)
	assert.Equal(t, "", builder.levelKey)
	nameFuncPtr := reflect.ValueOf(builder.levelFunc).Pointer()
	namePtrName1 := runtime.FuncForPC(nameFuncPtr).Name()
	assert.Equal(t, "", namePtrName1)
	assert.Equal(t, os.Stderr, builder.writer)
}

func TestWithContextHandler(t *testing.T) {
	builder := NewLoggerBuilder().WithContextHandler().(*defaultLoggerBuilder)
	assert.True(t, builder.useContextHandler)
}

func TestWithFormat(t *testing.T) {
	builder := NewLoggerBuilder().WithFormat(FormatJSON).(*defaultLoggerBuilder)
	assert.Equal(t, FormatJSON, builder.format)
}

func TestWithWriter(t *testing.T) {
	writer := os.Stdout
	builder := NewLoggerBuilder().WithWriter(writer).(*defaultLoggerBuilder)
	assert.Equal(t, writer, builder.writer)
}

func TestWithLevel(t *testing.T) {
	builder := NewLoggerBuilder().WithLevel(slog.LevelDebug).(*defaultLoggerBuilder)
	assert.Equal(t, slog.LevelDebug, builder.level)
}

func TestWithLevelString(t *testing.T) {
	builder := NewLoggerBuilder().WithLevelString("debug").(*defaultLoggerBuilder)
	assert.Equal(t, slog.LevelDebug, builder.level)
}

func TestWithLevelStringInvalid(t *testing.T) {
	expectPanic(t, func() {
		NewLoggerBuilder().WithLevelString("invalid")
	})
}

func TestWithLevelEnvVar(t *testing.T) {
	envVar := "LOG_LEVEL"
	builder := NewLoggerBuilder().WithLevelEnvVar(envVar).(*defaultLoggerBuilder)
	assert.Equal(t, envVar, builder.levelKey)
}

func TestWithLevelFunc(t *testing.T) {
	envVar := "LOG_LEVEL"
	levelFunc := getEnvLevelFunc()
	builder := NewLoggerBuilder().WithLevelFunc(envVar, levelFunc).(*defaultLoggerBuilder)
	assert.Equal(t, envVar, builder.levelKey)

	nameFuncPtr := reflect.ValueOf(builder.levelFunc).Pointer()
	namePtrName1 := runtime.FuncForPC(nameFuncPtr).Name()

	nameFuncPtr = reflect.ValueOf(levelFunc).Pointer()
	namePtrName2 := runtime.FuncForPC(nameFuncPtr).Name()

	assert.Equal(t, namePtrName2, namePtrName1)
}

func TestWithTimestampFormat(t *testing.T) {
	builder := NewLoggerBuilder().WithTimestampFormat(time.RFC1123).(*defaultLoggerBuilder)
	assert.Equal(t, time.RFC1123, builder.timestampFormat)
}

func TestWithTimestampFormatUnset(t *testing.T) {
	builder := NewLoggerBuilder().(*defaultLoggerBuilder)
	assert.Equal(t, "", builder.timestampFormat)
}

func TestWithTimestampFormatInvalid(t *testing.T) {
	expectPanic(t, func() {
		NewLoggerBuilder().WithTimestampFormat("invalid")
	})
}

func TestBuild_WithEnvVar(t *testing.T) {
	builder := NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelDebug).
		WithContextHandler().
		WithLevelEnvVar("LOG_LEVEL")

	require.NoError(t, os.Setenv("LOG_LEVEL", "INFO"))

	logger, levelVar := builder.Build()
	assert.NotNil(t, logger)
	assert.NotNil(t, levelVar)
	assert.Equal(t, slog.LevelInfo, levelVar.Level())
}

func TestBuild_WithLevelFunc(t *testing.T) {
	builder := NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelDebug).
		WithContextHandler().
		WithLevelFunc("LOG_LEVEL", getEnvLevelFunc())

	require.NoError(t, os.Setenv("LOG_LEVEL", "INFO"))

	logger, levelVar := builder.Build()
	assert.NotNil(t, logger)
	assert.NotNil(t, levelVar)
	assert.Equal(t, slog.LevelInfo, levelVar.Level())
}

func TestBuild_WithTimestampFormat(t *testing.T) {
	tests := []struct {
		name            string
		timestampFormat string
		shouldFormat    bool // true if timestampFormat was explicitly set
		expectedParseFn func(string) error
	}{
		{
			name:            "valid format RFC3339 is preserved",
			timestampFormat: time.RFC3339,
			shouldFormat:    true,
			expectedParseFn: func(s string) error {
				_, err := time.Parse(time.RFC3339, s)
				return err
			},
		},
		{
			name:            "valid format Kitchen is preserved",
			timestampFormat: time.Kitchen,
			shouldFormat:    true,
			expectedParseFn: func(s string) error {
				_, err := time.Parse(time.Kitchen, s)
				return err
			},
		},
		{
			name:            "empty string does not touch logger options",
			timestampFormat: "",
			shouldFormat:    false,
			expectedParseFn: func(s string) error {
				// When not formatting, slog uses RFC3339Nano by default
				_, err := time.Parse(time.RFC3339Nano, s)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			lb := &defaultLoggerBuilder{
				level:           slog.LevelInfo,
				format:          FormatJSON,
				writer:          &buf,
				timestampFormat: tt.timestampFormat,
			}

			logger, _ := lb.Build()
			logger.Info("test message")

			// Parse the JSON output and extract the time field
			var entry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
				t.Fatalf("failed to parse log output as JSON: %v", err)
			}

			rawTime, ok := entry["time"].(string)
			if !ok {
				t.Fatal("time field missing or not a string")
			}

			// Verify the timestamp parses correctly
			if err := tt.expectedParseFn(rawTime); err != nil {
				t.Errorf("timestamp %q does not parse with expected format: %v", rawTime, err)
			}
		})
	}
}
