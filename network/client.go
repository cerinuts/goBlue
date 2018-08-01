/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

//Client is a simple network socketclient
type Client struct {
	TargetIP, TargetPort string
	Connection           *net.Conn
	Reader               *bufio.Reader
}

//Connect connects to ip + port
func (c *Client) Connect(ip string, port string) (err error) {
	c.TargetIP = ip
	c.TargetPort = port
	conn, err := net.DialTimeout("tcp", ip+":"+port, 10*time.Second)
	c.Connection = &conn
	reader := bufio.NewReader(*(c.Connection))
	c.Reader = reader
	if err != nil {
		log.E("Error connecting", err.Error())
		return
	}
	return
}

//Recv waits and receives messages. BLOCKING
func (c *Client) Recv() (msg string, err error) {
	line, _, err := (*c.Reader).ReadLine()
	if err != nil {
		return
	}
	msg = string(line)
	return
}

//Sendln sends a message
func (c *Client) Sendln(msg string) {
	fmt.Fprintf(*(c.Connection), msg+"\n")
}

//Close closes the client
func (c *Client) Close() (err error) {
	err = (*(c.Connection)).Close()
	return
}
