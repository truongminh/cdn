package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	Server struct {
		Port int
		Addr string
	}
	Services []struct {
		Host   string
		Name   string
		Domain string
	} `toml:"service"`
}

func readConfig(ctx context.Context) config {
	configfile := flag.String("conf", "config.toml", "config file")
	flag.Parse()
	c := config{}
	log.Printf("config file %s", *configfile)
	buf, err := ioutil.ReadFile(*configfile)
	if err != nil {
		panic(err)
	}
	_, err = toml.Decode(string(buf), &c)
	if err != nil {
		panic(err)
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8081
	}
	c.Server.Addr = fmt.Sprintf(":%d", c.Server.Port)
	return c
}

func prerun(ctx context.Context) {
	install := flag.Bool("install", false, "install")
	unit := flag.String("unit", "", "systemd unit")
	dir := flag.String("dir", "", "run dir")
	username := flag.String("user", "", "user for install")
	flag.Parse()
	if *install {
		if *username == "" {
			log.Fatal("install requires user")
		}
		log.Printf("install on user %s", *username)
		ctrl := control{
			unit:     *unit,
			username: *username,
			runDir:   *dir,
		}
		err := ctrl.Install(ctx)
		log.Fatal(err)
	}
}
