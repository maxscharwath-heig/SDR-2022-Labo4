package main

import (
	. "SDR-Labo4/src/client"
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/input"
	"SDR-Labo4/src/utils/log"
	"flag"
	"time"
)

func main() {
	c, _ := config.LoadConfig()
	clientMode := flag.Int("mode", 1, "mode")
	flag.Parse()
	log.Logf(log.Info, "Selected mode: %d", *clientMode)

	for {
		word := input.StringInput("Enter a word: ").AddCheck(func(s string) bool {
			return len(s) > 0
		}, "Please enter a non-empty string").Read()

		if *clientMode == 1 {
			// Send word to all servers

			for _, server := range c.Servers {
				if client, err := CreateClient(server); err == nil {
					if err := client.Send([]byte(word)); err != nil {
						log.Logf(log.Error, "Error sending \"%s\" to %s: %s", word, server, err)
					} else {
						log.Logf(log.Info, "Sent \"%s\" to server %d", word, server.FullAddress())
					}
					client.Close()
				}
			}
		} else if *clientMode == 2 {
			root := input.BasicInput[int]("Enter the root server: ").AddCheck(func(s int) bool {
				return s >= 0 && s < len(c.Servers)
			}, "Please enter a valid server").Read()

			// TODO: Send the word to a server
			log.Logf(log.Info, "Sent \"%s\" to server %d", word, c.Servers[root].FullAddress())

		} else {
			log.Logf(log.Error, "Invalid mode %d selected, valid modes are: <1 | 2>", *clientMode)
		}

		// ask for result for a chosen server
		server := input.BasicInput[int]("Choose a server to ask for the result [0-%d]: ", len(c.Servers)-1).AddCheck(func(i int) bool {
			return i >= 0 && i < len(c.Servers)
		}, "Please enter a valid server id").Read()

		if client, err := CreateClient(c.Servers[server]); err == nil {
			client.Send([]byte("result"))
			select {
			case msg := <-client.GetMessage():
				log.Logf(log.Info, "Result: %s", string(msg))
				continue
			case <-time.After(5 * time.Second):
				log.Log(log.Info, "No result received")
				continue
			}
			client.Close()
		}
	}

}
