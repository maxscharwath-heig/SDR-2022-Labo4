package main

import (
	. "SDR-Labo4/src/client"
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/input"
	"SDR-Labo4/src/utils/log"
)

func main() {
	c, _ := config.LoadConfig()

	word := input.StringInput("Enter a word: ").AddCheck(func(s string) bool {
		return len(s) > 0
	}, "Please enter a non-empty string").Read()

	// send word to server

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
}
