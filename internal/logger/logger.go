package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var Logger *slog.Logger

// implementing it as a method of strung way too much work
func Fatal(msg string) {
	Logger.Error(msg)
	os.Exit(1)
}

func SetupLogger(filename string) {
	env := os.Getenv("ENV")
	if !(env == "prod" || env == "dev") {
		panic("ENV variable value is wrong. Use either `dev` or `prod`")
	}

	if env == "prod" {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(fmt.Errorf("failed to open log file: %w", err))
		}

		handler := slog.NewJSONHandler(file, nil)
		Logger = slog.New(handler)

		return
	} else {
		handler := slog.NewTextHandler(os.Stdout, nil)
		Logger = slog.New(handler)
	}
}
