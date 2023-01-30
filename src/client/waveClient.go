package client

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/input"
	"SDR-Labo4/src/utils/log"
	"time"
)

func WaveClientStart(servers []config.ServerConfig, word string) {
	for _, server := range servers {
		if client, err := CreateClient(server); err == nil {
			if err := client.Send(Message{Type: "start", Data: word}); err != nil {
				log.Logf(log.Error, "Error sending \"%s\" to %s: %s", word, server, err)
			} else {
				log.Logf(log.Info, "Sent \"%s\" to server %d", word, server.FullAddress())
			}
			client.Close()
		}
	}
}

func AskResultToServer(servers []config.ServerConfig, word string) {
	for {
		// ask for result for a chosen server
		server := input.BasicInput[int]("Choose a server to ask for the result [0-%d]: ", len(servers)-1).AddCheck(func(i int) bool {
			return i >= 0 && i < len(servers)
		}, "Please enter a valid server id").Read()

		if client, err := CreateClient(servers[server]); err == nil {
			if err := client.Send(Message{Type: "result"}); err != nil {
				log.Logf(log.Error, "Error sending \"result\" to %s: %s", servers[server], err)
			} else {
				log.Logf(log.Info, "Sent \"result\" to server %d", servers[server].FullAddress())
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
