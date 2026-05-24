package main

import (
	"echoserver/config"
	"echoserver/server"
	"flag"
	"log"
)

var cfg_sync config.Config
var cfg_async config.Config

func setupFlags() {
	flag.StringVar(&cfg_sync.Host, "sync-host", "0.0.0.0", "Host for Synchronous TCP Server to Listen on")
	flag.IntVar(&cfg_sync.Port, "sync-port", 7379, "Port for Synchronous TCP Server to Listen on")
	flag.StringVar(&cfg_async.Host, "async-host", "0.0.0.0", "Host for Asynchronous TCP Server to Listen on")
	flag.IntVar(&cfg_async.Port, "async-port", 7380, "Port for Asynchronous TCP Server to Listen on")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Echo Server is Live...")
	go server.RunSynchronousTCPServer(cfg_sync)
	go server.RunAsynchronousTCPServer(cfg_async)

	select {}
}
