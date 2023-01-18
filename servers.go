package main

import (
	"SDR-Labo4/src/config"
	. "SDR-Labo4/src/server"
	"SDR-Labo4/src/server/algo"
	"fmt"
	"log"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	servers := make([]*Server, len(c.Servers))

	for serverId := range c.Servers {
		server, err := CreateServer(c, serverId)
		if err != nil {
			panic(fmt.Sprintf("Could not create server %d: %s", serverId, err))
		}
		servers[serverId] = server
		server.Start()
	}

	for _, server := range servers {
		wave := algo.NewWave(*server, len(servers))
		go wave.Run()
	}

	select {}
}
