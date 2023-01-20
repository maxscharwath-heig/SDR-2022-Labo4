package main

import (
	"SDR-Labo4/src/config"
	. "SDR-Labo4/src/server"
	"SDR-Labo4/src/server/algo"
	"SDR-Labo4/src/utils/log"
	"flag"
)

func main() {
	serverMode := flag.Int("mode", 1, "mode")
	flag.Parse()
	log.Logf(log.Info, "Selected mode: %d", *serverMode)

	c, err := config.LoadConfig()
	if err != nil {
		log.Logf(log.Fatal, "Error loading config: %s", err)
	}

	servers := make([]*Server, len(c.Servers))

	for serverId := range c.Servers {
		server, err := CreateServer(c, serverId)
		if err != nil {
			log.Logf(log.Fatal, "Could not create server %d: %s", serverId, err)
		}
		servers[serverId] = server
		server.Start()
	}

	// Wave mode
	if *serverMode == 1 {
		for _, server := range servers {
			wave := algo.NewWave(*server, len(servers))
			go wave.Run()
		}
	} else if *serverMode == 2 {
		// Probes & Echoes mode
		for i, server := range servers {
			probes := algo.NewProbesAndEchoes(*server)
			if i != 0 {
				go probes.StartAsNode()
			}
		}
		// Start the first server as the "root"
		root := algo.NewProbesAndEchoes(*servers[0])
		go root.StartAsRoot()
	}

	select {}
}
