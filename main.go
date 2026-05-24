package main

import (
	"echoserver/config"
	"echoserver/server"
	"flag"
	"log"
)

var cfg config.Config

func setupFlags() {
	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "Host to Listen on")
	flag.IntVar(&cfg.Port, "port", 7379, "Port to Listen on")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Echo Server is Live...")
	server.RunSynchronousTCPServer(cfg)
}
