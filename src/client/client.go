package client

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/utils/log"
	"net"
)

type Client struct {
	conn     *net.UDPConn
	messages chan []byte
}

func (c *Client) Send(data []byte) (err error) {
	_, err = c.conn.Write(data)
	return
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
		conn: conn,
	}
	go client.processMessage()
	return client, nil
}
