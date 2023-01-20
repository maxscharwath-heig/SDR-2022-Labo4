package network

import "net"

type Code byte

const (
	CodeClient = Code(2)
	CodeServer = Code(4)
)

type Connection struct {
	*net.UDPConn
}

func (c *Connection) SendTo(code Code, data []byte, address *net.UDPAddr) (err error) {
	_, err = c.WriteToUDP(append([]byte{byte(code)}, data...), address)
	return
}

func (c *Connection) Send(code Code, data []byte) (err error) {
	_, err = c.Write(append([]byte{byte(code)}, data...))
	return
}

func (c *Connection) Receive() (code Code, data []byte, address *net.UDPAddr, err error) {
	buffer := make([]byte, 1024)
	n, address, err := c.ReadFromUDP(buffer)
	if err != nil {
		return
	}
	code = Code(buffer[0])
	data = buffer[1:n]
	return
}
