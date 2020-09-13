package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	prerun(ctx)
	log.Printf("initing ...")
	conf := readConfig(ctx)
	s := &server{addr: conf.Server.Addr}
	for _, cs := range conf.Services {
		service := &Service{
			Host:   cs.Host,
			Name:   cs.Name,
			Domain: cs.Domain,
		}
		service.init()
		s.services = append(s.services, service)
	}
	s.listen(ctx)
}
