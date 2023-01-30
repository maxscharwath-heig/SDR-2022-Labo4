// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Client to interacts with servers

package client

import (
	"SDR-Labo4/src/config"
	"SDR-Labo4/src/network"
	"SDR-Labo4/src/utils/log"
	"encoding/json"
	"fmt"
	"net"
	"sort"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Client struct {
	conn     network.Connection
	messages chan []byte
}

func (c *Client) Send(message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
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
		_, data, _, err := c.conn.Receive()
		if err != nil {
			return
		}
		go func() { c.messages <- data }()
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
		conn:     network.Connection{UDPConn: conn},
		messages: make(chan []byte),
	}
	go client.processMessage()
	return client, nil
}

func PrintResult(msg []byte, word string) {
	var results map[string]int
	json.Unmarshal(msg, &results)

	// Sort the map before printing
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Display results
	fmt.Printf("\n\n Result for word: \"%s\" \n", word)
	for _, k := range keys {
		fmt.Printf("%s: %d \n", k, results[k])
	}
}
