package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Host   string
	Name   string
	Domain string
	proxy  http.Handler
}

func (s *Service) init() {
	if s.Domain == "" {
		err := fmt.Errorf("service %s missing domain", s.Name)
		log.Fatal(err)
	}
	log.Printf("add service [%s -> %s]", s.Domain, s.Host)
	host, err := url.Parse(s.Host)
	if err != nil {
		err = fmt.Errorf("host %s: %w", s.Host, err)
		log.Fatal(err)
	}
	s.proxy = httputil.NewSingleHostReverseProxy(host)
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
	return
}
