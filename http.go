package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type server struct {
	addr     string
	services []*Service
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	for _, service := range s.services {
		if strings.HasSuffix(host, service.Domain) {
			service.ServeHTTP(w, r)
			return
		}
	}
	for _, service := range s.services {
		if service.Name == "default" {
			service.ServeHTTP(w, r)
			return
		}
	}
	http.Error(w, "host "+r.Host+" not found", http.StatusNotFound)
}

func (s *server) listen(ctx context.Context) {
	log.Printf("listen on %s", s.addr)
	err := http.ListenAndServe(s.addr, s)
	if err != nil {
		log.Fatal(err)
	}
}
