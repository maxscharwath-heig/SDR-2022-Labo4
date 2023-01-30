package algo

import (
	"SDR-Labo4/src/client"
	"SDR-Labo4/src/server"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
)

type MessageType string

const (
	Probe MessageType = "Probe"
	Echo  MessageType = "Echo"
)

type peMessage struct {
	MsgType MessageType     `json:"type"`
	From    int             `json:"from"` // Utilisé si Probe
	Data    map[string]Data `json:"data"` // Utilisé si Echo
	Word    string          `json:"word"` // Utilisé si Probe
}

type ProbeAndEchoes struct {
	server        server.Server
	data          map[string]Data
	neighbours    []int
	parent        int
	pending       chan client.Message
	initierClient server.Message
}

func NewProbesAndEchoes(server server.Server) *ProbeAndEchoes {
	s := &ProbeAndEchoes{
		server:     server,
		data:       make(map[string]Data),
		neighbours: server.GetNeighbours(),
		parent:     -1,
		pending:    make(chan client.Message),
	}
	go s.parseClientMessage()
	return s
}

func (pe *ProbeAndEchoes) Run() {
	for {
		select {
		case message := <-pe.pending:
			log.Logf(log.Info, "Server %d received start message: %v", pe.server.GetId(), message)
			switch message.Type {
			case "start":
				pe.startAsRoot(message.Data)
				if data, err := json.Marshal(pe.data); err == nil {
					_ = pe.initierClient.Reply(data)
				}

			case "probe":
				pe.startAsNode()
			}
		}
	}
}

func (pe *ProbeAndEchoes) startAsRoot(word string) {
	log.Logf(log.Info, "P&E algorithm started on server %d as the root", pe.server.GetId())
	pe.data = make(map[string]Data)

	letter := pe.server.GetConfig().Letter
	counter := CountLetter(word, letter)
	log.Logf(log.Info, "Server %d found %d %s in %s", pe.server.GetId(), counter, letter, word)
	pe.data[letter] = counter

	for _, neighbour := range pe.neighbours {
		pe.send(Probe, word, neighbour)
		log.Logf(log.Info, "Server %d sent %s to %d", pe.server.GetId(), Probe, neighbour)
	}

	log.Logf(log.Info, "Root server sent all its probes to neighbours")

	for range pe.neighbours {
		pe.receive()
		log.Logf(log.Info, "Server %d received %s", pe.server.GetId(), Echo)
	}

	log.Logf(log.Info, "Server %d (root) as received all echoes", pe.server.GetId())
}

func (pe *ProbeAndEchoes) startAsNode() {
	log.Logf(log.Info, "P&E algorithm started on server %d as a node", pe.server.GetId())
	pe.data = make(map[string]Data)
	message, _ := pe.receive() // Wait for a probe
	word := message.Word

	// Count the letters
	letter := pe.server.GetConfig().Letter
	counter := CountLetter(word, letter)
	log.Logf(log.Info, "Server %d found %d %s in %s", pe.server.GetId(), counter, letter, word)
	pe.data[letter] = counter

	pe.parent = message.From

	for _, neighbour := range pe.neighbours {
		if neighbour != pe.parent {
			log.Logf(log.Info, "Server %d sent %s to %d, content: %s", pe.server.GetId(), Probe, neighbour, word)
			pe.send(Probe, word, neighbour)
		}
	}

	for i := 0; i < len(pe.neighbours)-1; i++ {
		pe.receive()
	}

	pe.send(Echo, word, pe.parent)
}

func (pe *ProbeAndEchoes) send(msgType MessageType, word string, neighbour int) {
	message := peMessage{
		MsgType: msgType,
		From:    pe.server.GetId(),
		Data:    pe.data,
		Word:    word,
	}
	if data, err := json.Marshal(message); err == nil {
		pe.server.Send(data, neighbour)
	}
}

func (pe *ProbeAndEchoes) receive() (peMessage, error) {
	data := (<-pe.server.GetMessage()).Data
	var message peMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return message, err
	}

	// Update data
	for key, data := range message.Data {
		if _, ok := pe.data[key]; !ok {
			pe.data[key] = data
		}
	}

	log.Logf(log.Info, "Server %d got &s", pe.server.GetId(), message)
	return message, nil
}

func (pe *ProbeAndEchoes) parseClientMessage() {
	for {
		select {
		case message := <-pe.server.GetClientMessage():
			var m client.Message
			if err := json.Unmarshal(message.Data, &m); err != nil {
				log.Logf(log.Error, "Error while parsing client message: %v", err)
				continue
			}
			switch m.Type {
			case "start":
				pe.initierClient = message
				go func() { pe.pending <- m }()
			case "probe":
				go func() { pe.pending <- m }()
			}
		}
	}
}
