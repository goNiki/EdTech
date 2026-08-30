package main

import (
	"log/slog"
	"os"

	"edtech/internal/app"
)

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}

	slog.Info("application stopped gracefully")
}
