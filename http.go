package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

type ServerOptions struct {
	Http    string
	Https   string
	Domains []string
}

type server struct {
	opt      ServerOptions
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
	if s.opt.Http != "" {
		go func() {
			log.Printf("listen on http %s", s.opt.Http)
			err := http.ListenAndServe(s.opt.Http, s)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	if s.opt.Https != "" {
		go func() {
			log.Printf("listen on https %s", s.opt.Https)
			log.Printf("tls domains %s", s.opt.Domains)
			m := &autocert.Manager{
				Cache:      autocert.DirCache("tmp"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.opt.Domains...),
			}
			ss := &http.Server{
				Addr:      s.opt.Https,
				Handler:   s,
				TLSConfig: m.TLSConfig(),
			}
			err := ss.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
