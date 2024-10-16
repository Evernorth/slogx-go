package main

import (
	"context"
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
	"time"
)

func main() {

	logger, _ := slogx.NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(slogx.FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()

	ctx := context.Background()
	ctx1 := slogx.ContextWithAttrs(ctx, slog.String("goroutine", "1"))
	ctx2 := slogx.ContextWithAttrs(ctx, slog.String("goroutine", "2"))

	startTime := time.Now()

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 1000000; i++ {
			logger.InfoContext(ctx1, "message from goroutine 1", slog.Int("count", i))
		}
		done <- true

	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			logger.InfoContext(ctx2, "message from goroutine 2", slog.Int("count", i))
		}
		done <- true
	}()

	<-done
	<-done

	duration := time.Now().Sub(startTime).Seconds()
	logger.Info("Time taken to log 2 million messages", slog.Float64("duration_seconds", duration))
}
