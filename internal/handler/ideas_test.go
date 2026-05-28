package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jc8080/bootstrap_template/internal/handler"
)

func TestSubmitIdea(t *testing.T) {
	validBody := `{
		"types": ["Web app (browser)"],
		"typesOther": "",
		"genres": ["Education & learning"],
		"genresOther": "",
		"businessModels": ["Freemium"],
		"businessModelsOther": "",
		"barriers": ["Lack of time"],
		"barriersOther": "",
		"idea": "A tutoring marketplace.",
		"mvp": "Landing page plus waitlist."
	}`

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   string
	}{
		{name: "valid", method: http.MethodPost, body: validBody, wantStatus: http.StatusOK, wantBody: `{"ok":true}` + "\n"},
		{name: "method not allowed", method: http.MethodGet, body: validBody, wantStatus: http.StatusMethodNotAllowed},
		{name: "invalid json", method: http.MethodPost, body: `{`, wantStatus: http.StatusBadRequest},
		{name: "empty types", method: http.MethodPost, body: `{"types":[],"typesOther":"","genres":["Education & learning"],"genresOther":"","businessModels":["Freemium"],"businessModelsOther":"","barriers":["Lack of time"],"barriersOther":"","idea":"x","mvp":"y"}`, wantStatus: http.StatusBadRequest},
		{name: "other without text", method: http.MethodPost, body: strings.ReplaceAll(validBody, `"Web app (browser)"`, `"Other"`), wantStatus: http.StatusBadRequest},
		{name: "invalid type value", method: http.MethodPost, body: strings.ReplaceAll(validBody, `"Web app (browser)"`, `"Not a real type"`), wantStatus: http.StatusBadRequest},
		{name: "idea too long", method: http.MethodPost, body: strings.Replace(validBody, `"A tutoring marketplace."`, `"`+strings.Repeat("x", 6001)+`"`, 1), wantStatus: http.StatusBadRequest},
		{name: "typesOther too long", method: http.MethodPost, body: `{"types":["Other"],"typesOther":"` + strings.Repeat("x", 101) + `","genres":["Education & learning"],"genresOther":"","businessModels":["Freemium"],"businessModelsOther":"","barriers":["Lack of time"],"barriersOther":"","idea":"x","mvp":"y"}`, wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/ideas", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler.SubmitIdea(rr, req)
			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if tt.wantBody != "" && rr.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rr.Body.String(), tt.wantBody)
			}
		})
	}
}
