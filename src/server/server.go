package server

import (
	"SDR-Labo4/src/config"
	"log"
	"net"
)

type Server struct {
	config   config.Config
	id       int
	conn     *net.UDPConn
	address  *net.UDPAddr
	messages chan []byte
}

func CreateServer(config *config.Config, id int) (*Server, error) {
	address, err := config.Servers[id].Address()
	if err != nil {
		return nil, err
	}

	return &Server{
		config:   *config,
		id:       id,
		address:  address,
		messages: make(chan []byte),
	}, nil
}

func (s *Server) Start() error {
	log.Default().Printf("Starting server %d with address %s", s.id, s.address)
	conn, err := net.ListenUDP("udp", s.address)
	if err != nil {
		log.Default().Printf("Error starting server %d: %s", s.id, err)
		return err
	}
	s.conn = conn
	log.Default().Printf("Server %d started", s.id)
	go s.processMessage()
	return nil
}

func (s *Server) GetMessage() chan []byte {
	return s.messages
}

func (s *Server) Send(data []byte, serverId int) error {
	address, err := s.config.Servers[serverId].Address()
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(data, address)
	return err
}

func (s *Server) GetNeighbours() []int {
	return s.config.Servers[s.id].Neighbours
}

func (s *Server) Stop() {
	log.Default().Printf("Stopping server %d", s.id)
	s.conn.Close()
	log.Default().Printf("Server %d stopped", s.id)
}

func (s *Server) GetId() int {
	return s.id
}

func (s *Server) GetConfig() config.ServerConfig {
	return s.config.Servers[s.id]
}

func (s *Server) processMessage() {
	for {
		buffer := make([]byte, 1024)
		n, _, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Default().Printf("Error reading from UDP: %s", err)
			return
		}
		s.messages <- buffer[:n]
	}
}
