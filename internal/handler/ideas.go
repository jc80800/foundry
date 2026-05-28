package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const maxText = 6000
const maxOther = 100

var (
	allowedTypes = []string{
		"Mobile app (iOS/Android)",
		"Desktop app (Windows/macOS/Linux)",
		"Web app (browser)",
		"Browser extension",
		"API / developer library",
		"SaaS platform",
		"Marketplace / platform",
		"Hardware + software",
		"Game",
		"Plugin / integration",
		"CLI tool",
		"Content / media product",
		"Community / forum",
		"Newsletter / publication",
		"Physical product + digital layer",
		"Other",
	}
	allowedGenres = []string{
		"Education & learning",
		"Cooking & food",
		"Finance & investing",
		"Health & fitness",
		"Mental health & wellness",
		"Productivity & work",
		"Parenting & family",
		"Dating & relationships",
		"Travel & local",
		"Real estate",
		"Legal & compliance",
		"HR & hiring",
		"Sales & CRM",
		"Marketing & growth",
		"Design & creative tools",
		"Developer tools",
		"Gaming & entertainment",
		"Music & audio",
		"Sports",
		"Fashion & beauty",
		"Pets",
		"Sustainability & climate",
		"Nonprofit & social impact",
		"B2B operations",
		"E-commerce & retail",
		"Agriculture",
		"Construction & trades",
		"Automotive",
		"Science & research",
		"Other",
	}
	allowedBusinessModels = []string{
		"Free (no direct revenue)",
		"Freemium",
		"Monthly subscription",
		"Annual subscription",
		"One-time purchase",
		"Pay-per-use / usage-based",
		"Tiered plans",
		"Marketplace commission",
		"Advertising",
		"Affiliate / referrals",
		"Licensing (B2B)",
		"Enterprise / custom contracts",
		"Donations / tips",
		"Crowdfunding / preorders",
		"Data / insights (ethical B2B)",
		"Services + product bundle",
		"Other",
	}
	allowedBarriers = []string{
		"Lack of technical skills",
		"Lack of time",
		"Lack of funding / budget",
		"No co-founder / team",
		"Unclear MVP scope",
		"Don't know where to start",
		"Legal / regulatory uncertainty",
		"Need domain expertise",
		"Distribution / audience access",
		"Competing priorities",
		"Fear of failure / perfectionism",
		"Past failed attempts",
		"Tooling / stack overwhelm",
		"Design / UX skills gap",
		"Marketing / sales skills gap",
		"Other",
	}
)

type ideaSubmission struct {
	Types               []string `json:"types"`
	TypesOther          string   `json:"typesOther"`
	Genres              []string `json:"genres"`
	GenresOther         string   `json:"genresOther"`
	BusinessModels      []string `json:"businessModels"`
	BusinessModelsOther string   `json:"businessModelsOther"`
	Barriers            []string `json:"barriers"`
	BarriersOther       string   `json:"barriersOther"`
	Idea                string   `json:"idea"`
	MVP                 string   `json:"mvp"`
}

func SubmitIdea(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var sub ideaSubmission
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := validateSubmission(sub); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func validateSubmission(sub ideaSubmission) error {
	if err := validateGroup("types", sub.Types, allowedTypes, sub.TypesOther); err != nil {
		return err
	}
	if err := validateGroup("genres", sub.Genres, allowedGenres, sub.GenresOther); err != nil {
		return err
	}
	if err := validateGroup("businessModels", sub.BusinessModels, allowedBusinessModels, sub.BusinessModelsOther); err != nil {
		return err
	}
	if err := validateGroup("barriers", sub.Barriers, allowedBarriers, sub.BarriersOther); err != nil {
		return err
	}
	if strings.TrimSpace(sub.Idea) == "" || len(sub.Idea) > maxText {
		return errf("idea must be 1–%d characters", maxText)
	}
	if strings.TrimSpace(sub.MVP) == "" || len(sub.MVP) > maxText {
		return errf("mvp must be 1–%d characters", maxText)
	}
	return nil
}

func validateGroup(name string, values, allowed []string, other string) error {
	if len(values) == 0 {
		return errf("%s: select at least one option", name)
	}
	hasOther := false
	for _, v := range values {
		if !contains(allowed, v) {
			return errf("%s: invalid value %q", name, v)
		}
		if v == "Other" {
			hasOther = true
		}
	}
	other = strings.TrimSpace(other)
	if len(other) > maxOther {
		return errf("%sOther: max %d characters", name, maxOther)
	}
	if hasOther {
		if other == "" {
			return errf("%sOther: required when Other is selected", name)
		}
		return nil
	}
	if other != "" {
		return errf("%sOther: must be empty unless Other is selected", name)
	}
	return nil
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

type validationError string

func (e validationError) Error() string { return string(e) }

func errf(format string, args ...any) error {
	return validationError(fmt.Sprintf(format, args...))
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
