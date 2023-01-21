package algo

import (
	"SDR-Labo4/src/server"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
	"fmt"
	"sort"
)

type Message struct {
	From   int             `json:"from"`
	Active bool            `json:"active"`
	Data   map[string]Data `json:"data"`
}

func (m Message) String() string {
	return fmt.Sprintf("Data from %d, active: %t, data: %v", m.From, m.Active, m.Data)
}

type Wave struct {
	server     server.Server
	nbNodes    int
	data       map[string]Data
	neighbours map[int]bool
}

func NewWave(server server.Server, nbNodes int) *Wave {
	w := &Wave{
		server:     server,
		nbNodes:    nbNodes,
		data:       make(map[string]Data),
		neighbours: make(map[int]bool),
	}
	for _, neighbour := range server.GetConfig().Neighbours {
		w.neighbours[neighbour] = true
	}
	return w
}

func (w *Wave) Run() {
	log.Logf(log.Info, "Starting wave algorithm on server %d", w.server.GetId())

	w.waitForClient()

	for !w.isTopologyComplete() {
		// Send message to activeNeigbours
		for _, neighbour := range w.getNeighbours() {
			w.send(true, neighbour)
		}

		for neighbour := range w.neighbours {
			message, _ := w.receive()

			for key, data := range message.Data {
				if _, ok := w.data[key]; !ok {
					w.data[key] = data
				}
			}

			if !message.Active {
				log.Logf(log.Info, "Server %d is removing %d from its neighbours", w.server.GetId(), neighbour)
				delete(w.neighbours, message.From)
			}
		}
	}

	for _, neighbour := range w.getNeighbours() {
		w.send(false, neighbour)
	}

	for _, neighbour := range w.getNeighbours() {
		log.Logf(log.Warn, "Server %d is purging %d", w.server.GetId(), neighbour)
		_, _ = w.receive()
		log.Logf(log.Warn, "Server %d purged %d", w.server.GetId(), neighbour)
	}
	log.Logf(log.Error, "Wave algorithm on server %d is done with data: %v", w.server.GetId(), w.data)

	w.respondToClient()
}

func (w *Wave) getNeighbours() []int {
	keys := make([]int, 0)
	for k, _ := range w.neighbours {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (w *Wave) isTopologyComplete() bool {
	return len(w.data) == w.nbNodes
}

func (w *Wave) send(active bool, neighbour int) {
	message := Message{
		From:   w.server.GetId(),
		Active: active,
		Data:   w.data,
	}
	log.Logf(log.Info, "Server %d is sending %v to %d", w.server.GetId(), message, neighbour)
	if data, err := json.Marshal(message); err == nil {
		w.server.Send(data, neighbour)
	}
}

func (w *Wave) receive() (Message, error) {
	for {
		select {
		case message := <-w.server.GetMessage():
			var m Message
			if err := json.Unmarshal(message.Data, &m); err != nil {
				return m, err
			}
			log.Logf(log.Info, "Server %d received %v", w.server.GetId(), m)
			return m, nil
		}
	}
}

func (w *Wave) waitForClient() {
	// wait for the client to be ready
	for {
		select {
		case message := <-w.server.GetClientMessage():
			word := string(message.Data)
			letter := w.server.GetConfig().Letter
			counter := CountLetter(word, letter)
			log.Logf(log.Info, "Server %d found %d %s in %s", w.server.GetId(), counter, w.server.GetConfig().Letter, word)
			w.data[letter] = counter
			return
		}
	}
}

func (w *Wave) respondToClient() {
	// the process is done, send the result to the client
	for {
		select {
		case message := <-w.server.GetClientMessage():
			if data, err := json.Marshal(w.data); err == nil {
				_ = message.Reply(data)
			}
		}
	}
}
