package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)

	app := &application{
		logger: logger,
	}

	mux := app.routes()

	logger.Info("Starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)
	logger.Error("Failed to start server", "error", err)
	os.Exit(1)
}
