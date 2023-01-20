package server

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/network"
	"SDR-Labo4/src/utils/log"
	"net"
)

type Message struct {
	server *Server
	From   *net.UDPAddr
	Data   []byte
}

func (m *Message) Reply(data []byte) (err error) {
	return m.server.conn.SendTo(network.CodeServer, data, m.From)
}

type Server struct {
	config         config.Config
	id             int
	conn           network.Connection
	address        *net.UDPAddr
	messages       chan Message
	clientMessages chan Message
}

func CreateServer(config *config.Config, id int) (*Server, error) {
	address, err := config.Servers[id].Address()
	if err != nil {
		return nil, err
	}

	return &Server{
		config:         *config,
		id:             id,
		address:        address,
		messages:       make(chan Message),
		clientMessages: make(chan Message),
	}, nil
}

func (s *Server) Start() error {
	log.Logf(log.Info, "Starting server %d with address %s", s.id, s.address)
	conn, err := net.ListenUDP("udp", s.address)
	if err != nil {
		log.Logf(log.Error, "Error starting server %d: %s", s.id, err)
		return err
	}
	s.conn = network.Connection{UDPConn: conn}
	go s.processMessage()
	return nil
}

func (s *Server) GetMessage() chan Message {
	return s.messages
}

func (s *Server) GetClientMessage() chan Message {
	return s.clientMessages
}

func (s *Server) Send(data []byte, serverId int) error {
	address, err := s.config.Servers[serverId].Address()
	if err != nil {
		log.Logf(log.Error, "Error sending to server %d: %s", s.id, err)
		return err
	}
	return s.conn.SendTo(network.CodeServer, data, address)
}

func (s *Server) GetNeighbours() []int {
	return s.config.Servers[s.id].Neighbours
}

func (s *Server) Stop() {
	s.conn.Close()
	log.Logf(log.Info, "Server %d stopped", s.id)
}

func (s *Server) GetId() int {
	return s.id
}

func (s *Server) GetConfig() config.ServerConfig {
	return s.config.Servers[s.id]
}

func (s *Server) processMessage() {
	for {
		code, data, from, err := s.conn.Receive()
		if err != nil {
			log.Log(log.Debug, "Server closed")
			return
		}

		log.Logf(log.Debug, "Received message[%d] from %s: %s", code, from, string(data))

		message := Message{
			server: s,
			From:   from,
			Data:   data,
		}

		switch code {
		case network.CodeServer:
			s.messages <- message
		case network.CodeClient:
			s.clientMessages <- message
		}
	}
}
