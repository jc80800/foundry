package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/jc8080/bootstrap_template/internal/config"
	"github.com/jc8080/bootstrap_template/internal/handler"
	"github.com/jc8080/bootstrap_template/web"
)

type Server struct {
	cfg config.Config
	mux *http.ServeMux
}

func New(cfg config.Config) (*Server, error) {
	staticFS, err := fs.Sub(web.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("load embedded web assets: %w", err)
	}

	s := &Server{cfg: cfg, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /api/health", handler.Health)
	s.mux.HandleFunc("POST /api/ideas", handler.SubmitIdea)
	s.mux.Handle("/", http.FileServer(http.FS(staticFS)))
	return s, nil
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.cfg.Port)
}
