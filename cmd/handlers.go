package main

import (
	"net/http"
	"strings"
)

// IdeaInput is the form payload posted to POST /api/ideas.
// Persist these fields when wiring storage.
type IdeaInput struct {
	Title       string
	Description string
	Category    string
	Tags        string
	Contact     string
}

type homeData struct {
	Submitted bool
}

func (app *application) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{
		Submitted: r.URL.Query().Get("ok") == "1",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error("Failed to render homepage", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	idea := IdeaInput{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Category:    strings.TrimSpace(r.FormValue("category")),
		Tags:        strings.TrimSpace(r.FormValue("tags")),
		Contact:     strings.TrimSpace(r.FormValue("contact")),
	}

	if idea.Title == "" || idea.Description == "" || idea.Category == "" {
		http.Error(w, "title, description, and category are required", http.StatusBadRequest)
		return
	}

	// Persist IdeaInput here (DB / API) when ready.
	app.logger.Info("idea submitted",
		"title", idea.Title,
		"description", idea.Description,
		"category", idea.Category,
		"tags", idea.Tags,
		"contact", idea.Contact,
	)

	http.Redirect(w, r, "/?ok=1", http.StatusSeeOther)
}
