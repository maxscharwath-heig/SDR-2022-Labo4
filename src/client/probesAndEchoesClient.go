// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Client to interact with probes & echoes servers

package client

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/log"
)

func ProbesAndEchoesClientStart(servers []config.ServerConfig, word string, root int) {
	for i, server := range servers {
		if i == root {
			continue
		}
		// Start the nodes
		if client, err := CreateClient(server); err == nil {
			message := Message{Type: "probe"}

			if err := client.Send(message); err != nil {
				log.Logf(log.Error, "Error sending \"%s\" to %s: %s", word, server, err)
			} else {
				log.Logf(log.Info, "Sent \"%s\" to server %d", message, server.FullAddress())
			}
			client.Close()
		}
	}

	// Start the root
	rootSrv := servers[root]
	if client, err := CreateClient(rootSrv); err == nil {
		message := Message{Type: "start", Data: word}

		if err := client.Send(message); err != nil {
			log.Logf(log.Error, "Error sending \"%s\" to %s: %s", word, rootSrv, err)
		} else {
			log.Logf(log.Info, "Sent \"%s\" to server %d", message, rootSrv.FullAddress())
		}
		select {
		case msg := <-client.GetMessage():
			PrintResult(msg, word)
		}
		client.Close()
	}
}
