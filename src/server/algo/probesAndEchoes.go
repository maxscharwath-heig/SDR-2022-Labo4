package algo

import (
	"SDR-Labo4/src/server"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
)

type MessageType string

const (
	Probe MessageType = "Probe"
	Echo  MessageType = "Echo"
)

type PEMessage struct {
	MsgType MessageType  `json:"type"`
	From    int          `json:"from"` // Utilisé si Probe
	Data    map[int]Data `json:"data"` // Utilisé si Probe
}

type ProbeAndEchoes struct {
	server     server.Server
	data       map[int]Data
	neighbours []int
	parent     int
}

func NewProbesAndEchoes(server server.Server) *ProbeAndEchoes {
	s := &ProbeAndEchoes{
		server:     server,
		data:       make(map[int]Data),
		neighbours: server.GetNeighbours(),
		parent:     -1,
	}
	s.data[server.GetId()] = server.GetConfig().Letter
	return s
}

func (pe *ProbeAndEchoes) StartAsRoot() {
	log.Logf(log.Info, "P&E algorithm started on server %d as the root", pe.server.GetId())

	for neighbour := range pe.neighbours {
		pe.send(Probe, neighbour)
		log.Logf(log.Info, "Server %d sent %s", pe.server.GetId(), Probe)
	}

	for range pe.neighbours {
		pe.receive()
		log.Logf(log.Info, "Server %d received %s", pe.server.GetId(), Echo)
	}

	log.Logf(log.Info, "Server %d (root) as received all echoes", pe.server.GetId())
}

func (pe *ProbeAndEchoes) StartAsNode() {
	log.Logf(log.Info, "P&E algorithm started on server %d as a node", pe.server.GetId())
	message, _ := pe.receive()

	text := message.Data
	pe.parent = message.From

	for i, neighbour := range pe.neighbours {
		if i != pe.parent {
			log.Logf(log.Info, "Server %d sent %s, content: %s", pe.server.GetId(), Probe, text)
			pe.send(Probe, neighbour) // TODO: send text
		}
	}

	for i := 0; i < len(pe.neighbours)-1; i++ {
		pe.receive()
	}

	pe.send(Echo, pe.parent)
}

func (pe *ProbeAndEchoes) send(msgType MessageType, neighbour int) {
	message := PEMessage{
		MsgType: msgType,
		From:    pe.server.GetId(),
		Data:    pe.data,
	}
	if data, err := json.Marshal(message); err == nil {
		pe.server.Send(data, neighbour)
	}
}

func (pe *ProbeAndEchoes) receive() (PEMessage, error) {
	data := (<-pe.server.GetMessage()).Data
	var message PEMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return message, err
	}
	log.Logf(log.Info, "Server %d got &s", pe.server.GetId(), message)
	return message, nil
}
