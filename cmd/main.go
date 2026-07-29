package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger    *slog.Logger
	homeTmpl  *template.Template
	ideasTmpl *template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)

	homeTmpl, ideasTmpl, err := parseTemplates()
	if err != nil {
		logger.Error("Failed to parse templates", "error", err)
		os.Exit(1)
	}

	app := &application{
		logger:    logger,
		homeTmpl:  homeTmpl,
		ideasTmpl: ideasTmpl,
	}

	mux := app.routes()

	logger.Info("Starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, mux)
	logger.Error("Failed to start server", "error", err)
	os.Exit(1)
}

func parseTemplates() (*template.Template, *template.Template, error) {
	base, err := template.ParseFiles("ui/html/base.tmpl")
	if err != nil {
		return nil, nil, err
	}

	homeTmpl, err := base.Clone()
	if err != nil {
		return nil, nil, err
	}
	homeTmpl, err = homeTmpl.ParseFiles("ui/html/home.tmpl")
	if err != nil {
		return nil, nil, err
	}

	ideasTmpl, err := base.Clone()
	if err != nil {
		return nil, nil, err
	}
	ideasTmpl, err = ideasTmpl.ParseFiles("ui/html/ideas.tmpl")
	if err != nil {
		return nil, nil, err
	}

	return homeTmpl, ideasTmpl, nil
}
