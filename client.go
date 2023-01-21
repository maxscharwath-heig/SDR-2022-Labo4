package main

import (
	. "SDR-Labo4/src/client"
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/input"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
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

			for i, server := range c.Servers {
				if client, err := CreateClient(server); err == nil {
					var toSend string
					if root == i {
						toSend = word
					} else {
						toSend = "wait"
					}

					if err := client.Send([]byte(toSend)); err != nil {
						log.Logf(log.Error, "Error sending \"%s\" to %s: %s", word, server, err)
					} else {
						log.Logf(log.Info, "Sent \"%s\" to server %d", word, server.FullAddress())
					}
					client.Close()
				}
			}

		} else {
			log.Logf(log.Error, "Invalid mode %d selected, valid modes are: <1 | 2>", *clientMode)
		}

		for {
			// ask for result for a chosen server
			server := input.BasicInput[int]("Choose a server to ask for the result [0-%d]: ", len(c.Servers)-1).AddCheck(func(i int) bool {
				return i >= 0 && i < len(c.Servers)
			}, "Please enter a valid server id").Read()

			if client, err := CreateClient(c.Servers[server]); err == nil {
				if err := client.Send([]byte("result")); err != nil {
					log.Logf(log.Error, "Error sending \"result\" to %s: %s", c.Servers[server], err)
				} else {
					log.Logf(log.Info, "Sent \"result\" to server %d", c.Servers[server].FullAddress())
				}
				select {
				case msg := <-client.GetMessage():
					PrintResult(msg, word)
				case <-time.After(5 * time.Second):
					log.Log(log.Info, "No result received")
				}
				client.Close()
			}
		}
	}
}

func PrintResult(msg []byte, word string) {
	var objmap map[string]int
	json.Unmarshal(msg, &objmap)
	log.Logf(log.Info, "Result for \"%s\":", word)
	for k, v := range objmap {
		log.Logf(log.Info, "%s: %d", k, v)
	}
}
