package main

import (
	"log"
	"net/http"

	"github.com/jc8080/bootstrap_template/internal/config"
	"github.com/jc8080/bootstrap_template/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", srv.Addr())
	log.Fatal(http.ListenAndServe(srv.Addr(), srv.Handler()))
}
