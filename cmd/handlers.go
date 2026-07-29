package main

import "net/http"

func (app *application) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("GET / homepage handler", "request", r)
	w.Write([]byte("This is the homepage"))
}
