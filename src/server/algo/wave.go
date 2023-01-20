package algo

import (
	"SDR-Labo4/src/server"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
	"fmt"
)

type Data = string

type Message struct {
	From   int          `json:"from"`
	Active bool         `json:"active"`
	Data   map[int]Data `json:"data"`
}

func (m Message) String() string {
	return fmt.Sprintf("Data from %d, active: %t, data: %v", m.From, m.Active, m.Data)
}

type Wave struct {
	server     server.Server
	nbNodes    int
	data       map[int]Data
	neighbours map[int]bool
}

func NewWave(server server.Server, nbNodes int) *Wave {
	w := &Wave{
		server:     server,
		nbNodes:    nbNodes,
		data:       make(map[int]Data),
		neighbours: make(map[int]bool),
	}
	w.data[server.GetId()] = server.GetConfig().Letter
	for _, neighbour := range server.GetNeighbours() {
		w.neighbours[neighbour] = true
	}
	return w
}

func (w *Wave) Run() {
	log.Logf(log.Info, "Starting wave algorithm on server %d", w.server.GetId())

	w.waitForClient()

	for !w.isTopologyComplete() {

		// Send message to neighbours
		for neighbour := range w.neighbours {
			w.send(true, neighbour)
		}

		// Receive messages from neighbours
		for _ = range w.neighbours {
			message, _ := w.receive()
			w.neighbours[message.From] = message.Active
			// Update data
			for id, data := range message.Data {
				if _, ok := w.data[id]; !ok {
					w.data[id] = data
				}
			}

		}
	}
	for neighbour, active := range w.neighbours {
		if active {
			w.send(false, neighbour)
		}
	}
	for _, active := range w.neighbours {
		if active {
			log.Logf(log.Warn, "Server %d is purging", w.server.GetId())
			_, _ = w.receive()
		}
	}
	log.Logf(log.Info, "Wave algorithm on server %d is done with data: %v", w.server.GetId(), w.data)
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
	log.Logf(log.Trace, "Server %d sending message to %d: %s", w.server.GetId(), neighbour, message)
	if data, err := json.Marshal(message); err == nil {
		w.server.Send(data, neighbour)
	}
}

func (w *Wave) receive() (Message, error) {
	data := <-w.server.GetMessage()
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return message, err
	}
	log.Logf(log.Trace, "Server %d received message from %d: %s", w.server.GetId(), message.From, message)
	return message, nil
}

func (w *Wave) waitForClient() {
	log.Log(log.Info, "Waiting for client to send data")
	word := string(<-w.server.GetMessage())
	log.Logf(log.Info, "Server %d received word: %s", w.server.GetId(), word)
}
