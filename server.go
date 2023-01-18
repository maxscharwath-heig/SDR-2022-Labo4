package main

import (
	"SDR-Labo4/src/config"
	. "SDR-Labo4/src/server"
	"flag"
)

func main() {
	serverIdInput := flag.Int("id", 0, "server id")
	flag.Parse()

	c, _ := config.LoadConfig()

	serverId := *serverIdInput

	server, _ := CreateServer(c, serverId)
	server.Start()
}
