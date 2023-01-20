package client

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/network"
	"SDR-Labo4/src/utils/log"
	"net"
)

type Client struct {
	conn     network.Connection
	messages chan []byte
}

func (c *Client) Send(data []byte) error {
	return c.conn.Send(network.CodeClient, data)
}

func (c *Client) GetMessage() chan []byte {
	return c.messages
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) processMessage() {
	for {
		buffer := make([]byte, 1024)
		n, _, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Log(log.Debug, "Client closed")
			return
		}
		c.messages <- buffer[:n]
	}
}

func CreateClient(server config.ServerConfig) (*Client, error) {
	addr, _ := server.Address()
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Logf(log.Error, "Error connecting to server: %s", err)
		return nil, err
	}
	client := &Client{
		conn: network.Connection{UDPConn: conn},
	}
	go client.processMessage()
	return client, nil
}
