package main

import (
	"context"
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"net/http"
	"os"
)

var (
	// This gets us a slog.Logger with context support that logs in JSON format to stdout.
	logger, levelVar = slogx.NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(slogx.FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()
)

func handler(response http.ResponseWriter, request *http.Request) {

	// Get the ctx from the request
	ctx := request.Context()

	// Extract the id from the query parameters
	id := request.URL.Query().Get("id")

	// Log the request
	logger.InfoContext(ctx, "Received request", slog.String("id", id))

	// Write the response
	response.WriteHeader(http.StatusOK)
	_, err := response.Write([]byte("Id: " + id))
	if err != nil {
		return
	}

}
func init() {
	// Set the default level manager
	// Enroll the levelVar to be managed from an environment variable.
	err := slogx.GetLevelManager().ManageLevelFromEnv(levelVar, "APP_LOG_LEVEL")
	if err != nil {
		panic(err)
	}

	// Trigger the LevelManager to update the levels from the environment.
	// This should be done again if the environment is updated, so that the logger level
	// can be updated without restarting the application.
	slogx.GetLevelManager().UpdateLevels()

	slog.InfoContext(context.Background(), "Logger initialized", slog.String("level", levelVar.Level().String()))
}

func main() {

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
