package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.getHomepageHandler)
	mux.HandleFunc("POST /api/ideas", app.createIdeaHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))

	return mux
}
