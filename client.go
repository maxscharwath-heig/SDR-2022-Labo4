// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Run a client with selected mode

package main

import (
	. "SDR-Labo4/src/client"
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/input"
	"SDR-Labo4/src/utils/log"
	"flag"
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
			WaveClientStart(c.Servers, word)
			AskResultToServer(c.Servers, word)
		} else if *clientMode == 2 {
			root := input.BasicInput[int]("Enter the root server: ").AddCheck(func(s int) bool {
				return s >= 0 && s < len(c.Servers)
			}, "Please enter a valid server").Read()

			ProbesAndEchoesClientStart(c.Servers, word, root)
		} else {
			log.Logf(log.Error, "Invalid mode %d selected, valid modes are: <1 | 2>", *clientMode)
		}
	}
}
