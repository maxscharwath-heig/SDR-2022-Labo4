package server

import (
	"SDR-Labo4/src/config"
	"net"
)

type Server struct {
	id      int
	address *net.UDPAddr
}

func CreateServer(config *config.Config, id int) (*Server, error) {
	address, err := config.Servers[id].Address()
	if err != nil {
		return nil, err
	}
	return &Server{
		id,
		address,
	}, nil
}

func (s *Server) Start() {
	return
}
