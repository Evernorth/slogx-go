# log/slogx
A collection of `slog` extensions. 
* `ContextHandler` allows you to add `slog` attributes (`slog.Attr` instances) to a `context.Context`.  These attributes are added to log records when the `*Context` function variants (`InfoContext`, `ErrorContext`, etc) on the logger are used. 
* `LoggerBuilder` provides a simple way to build `slog.Logger` instances.
* `LevelManager` provides a way to manage `slog.LevelVar` instances from environment variables.

We are planning to open source this package in the future.

## Usage Examples
### Context-aware logging
Create a context-aware `slog.Logger`.
```go
import (
    "github.sys.cigna.com/cigna/common-go/log/slogx"
    "log/slog"
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

    // Enroll the levelVar to be managed from an environment variable.
    must.OK(slogx.GetLevelManager().ManageLevelFromEnv(levelVar, "APP_LOG_LEVEL"))

	// Trigger the LevelManager to update the levels from the environment.
	// This should be done again if the environment is updated, so that the logger level
	// can be updated without restarting the application.
    slogx.GetLevelManager().UpdateLevels()
)
```
Use the logger, making sure to use a `*Context` variant function.
```go
import (
    "github.com/go-chi/chi/v5"
    "log/slog"
    "net/http"
)

func getEmployee(httpRespWriter http.ResponseWriter, httpReq *http.Request) {

    ctx := httpReq.Context()
    resourceId := chi.URLParam(httpReq, "resourceId")
    
    // Using the *Context log functions will cause any logging context attributes added by middleware components
    // to be included on each log record.
    logger.InfoContext(ctx, "Fetching Employee...",
        slog.String("resourceId", resourceId))

    ...
}
```