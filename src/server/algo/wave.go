package algo

import (
	"SDR-Labo4/src/server"
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	From     int   `json:"from"`
	Active   bool  `json:"active"`
	Topology []int `json:"topology"`
}

func (m Message) String() string {
	return fmt.Sprintf("Data from %d, active: %t, topology: %s", m.From, m.Active, m.Topology)
}

type Wave struct {
	server     server.Server
	neighbours map[int]bool // map of neighbours (true if active)
	topology   []bool
}

func NewWave(server server.Server, nbNodes int) *Wave {
	w := &Wave{
		server:     server,
		neighbours: make(map[int]bool),
		topology:   make([]bool, nbNodes),
	}

	// topology is used to keep data
	w.topology[server.GetId()] = true

	for _, neighbour := range server.GetNeighbours() {
		w.neighbours[neighbour] = true
	}

	return w
}

func (w *Wave) Run() {
	log.Default().Printf("Starting wave algorithm on server %d", w.server.GetId())
	for !w.isTopologyComplete() {

		// Send message to neighbours
		message := Message{
			From:     w.server.GetId(),
			Active:   true,
			Topology: w.getTopology(),
		}
		for neighbour := range w.neighbours {
			w.send(message, neighbour)
		}

		// Receive messages from neighbours
		for neighbour := range w.neighbours {
			message, _ := w.receive()
			w.neighbours[neighbour] = message.Active

			//merge topology
			for _, node := range message.Topology {
				w.topology[node] = true
			}

		}
	}
	message := Message{
		From:     w.server.GetId(),
		Active:   false,
		Topology: w.getTopology(),
	}
	for neighbour, active := range w.neighbours {
		if active {
			w.send(message, neighbour)
		}
	}
	for _, active := range w.neighbours {
		if active {
			_, _ = w.receive()
		}
	}
	log.Default().Printf("Wave algorithm on server %d is done with topology %v", w.server.GetId(), w.getTopology())
}

func (w *Wave) isTopologyComplete() bool {
	for _, active := range w.topology {
		if !active {
			return false
		}
	}
	return true
}

func (w *Wave) getTopology() []int {
	var topology []int
	for i, active := range w.topology {
		if active {
			topology = append(topology, i)
		}
	}
	return topology
}

func (w *Wave) send(message Message, neighbour int) {
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
	return message, nil
}
