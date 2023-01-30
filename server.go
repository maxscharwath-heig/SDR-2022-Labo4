// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Run a server with selected mode

package main

import (
	"SDR-Labo4/src/config"
	. "SDR-Labo4/src/server"
	"SDR-Labo4/src/server/algo"
	"SDR-Labo4/src/utils/log"
	"flag"
	"os"
)

func main() {
	serverIdInput := flag.Int("id", 0, "server id")
	serverMode := flag.Int("mode", 1, "mode")
	flag.Parse()

	log.Logf(log.Info, "Selected mode: %d", *serverMode)

	c, _ := config.LoadConfig()

	serverId := *serverIdInput

	server, _ := CreateServer(c, serverId)
	server.Start()

	if *serverMode == 1 {
		wave := algo.NewWave(*server, len(c.Servers))
		go wave.Run()
	} else if *serverMode == 2 {
		probes := algo.NewProbesAndEchoes(*server)
		go probes.Run()
	} else {
		log.Logf(log.Error, "Invalid mode %d selected, valid modes are: <1 | 2>", *serverMode)
		os.Exit(1)
	}
	select {}
}
