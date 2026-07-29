package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	tmpl   *template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)

	tmpl, err := template.ParseFiles("ui/html/base.tmpl", "ui/html/home.tmpl")
	if err != nil {
		logger.Error("Failed to parse templates", "error", err)
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		tmpl:   tmpl,
	}

	mux := app.routes()

	logger.Info("Starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, mux)
	logger.Error("Failed to start server", "error", err)
	os.Exit(1)
}
