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
	ActiveNav string
}

type ideasData struct {
	ActiveNav string
	Ideas     []IdeaInput
}

func (app *application) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{
		Submitted: r.URL.Query().Get("ok") == "1",
		ActiveNav: "submit",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.homeTmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error("Failed to render homepage", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getIdeasHandler renders the ideas exhibit.
// TODO: load ideas from the database instead of mock data.
func (app *application) getIdeasHandler(w http.ResponseWriter, r *http.Request) {
	data := ideasData{
		ActiveNav: "ideas",
		Ideas:     mockIdeas(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.ideasTmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error("Failed to render ideas page", "error", err)
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

// mockIdeas returns placeholder exhibit entries until persistence is wired.
// TODO: replace with a database query.
func mockIdeas() []IdeaInput {
	return []IdeaInput{
		{
			Title:       "Workshop queue board",
			Description: "A shared board where makers post what they are building this week and claim spare bench time without a group chat thread. The board should show machine availability, material conflicts, and a quiet way to offer help on someone else's build without turning the shop into another messaging app. Ideal version syncs with a simple wall display near the entrance.",
			Category:    "Product",
			Tags:        "collaboration, scheduling",
			Contact:     "maya@foundry.dev",
		},
		{
			Title:       "Material scrap matcher",
			Description: "Scan leftover stock and match it to open project briefs so offcuts get used before they hit the bin.",
			Category:    "Tool",
			Tags:        "inventory, sustainability",
			Contact:     "",
		},
		{
			Title:       "Acoustic sketch diary",
			Description: "Record short ambient clips from the studio and pin them to sketches so future you remembers how a piece felt in the room. Over a season, the diary becomes a map of which corners ring, which benches hum, and which finishes change the sound of a space. Export should be a simple timeline, not a social feed.",
			Category:    "Research",
			Tags:        "audio, journaling",
			Contact:     "leo",
		},
		{
			Title:       "One-button fixture jig",
			Description: "A parametric jig generator that outputs cut lists for common clamp setups from a single measurement.",
			Category:    "Experiment",
			Tags:        "cnc, fixtures",
			Contact:     "",
		},
		{
			Title:       "Visitor idea lottery",
			Description: "Drop a slip when you tour the space; once a month a drawn idea gets a half-day build slot with a mentor. Slips stay anonymous unless the visitor wants credit. The draw is public, the build notes go back on the wall, and rejected slips stay in a quiet archive so good ideas are not lost to the moment.",
			Category:    "Other",
			Tags:        "community",
			Contact:     "front desk",
		},
		{
			Title:       "Finish recipe cards",
			Description: "Printable cards for stain and oil mixes, with batch notes and a QR that deep-links to the last successful run.",
			Category:    "Product",
			Tags:        "finishing, docs",
			Contact:     "sam@workshop.local",
		},
	}
}
