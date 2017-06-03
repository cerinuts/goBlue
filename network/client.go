/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

package network

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	TargetIP, TargetPort string
	Connection           *net.Conn
	Reader				*bufio.Reader
}

func (c *Client) Connect(ip string, port string) (err error) {
	c.TargetIP = ip
	c.TargetPort = port
	conn, err := net.Dial("tcp", ip+":"+port)
	c.Connection = &conn
	reader := bufio.NewReader(*(c.Connection))
	c.Reader = reader
	if err != nil {
		println("Error connecting", err.Error())
		return
	}
	return
}

func (c *Client) Recv() (msg string, err error) {
	line, _, erro := (*c.Reader).ReadLine()
	err = erro
	if err != nil {
		println("Error reading from Socket", erro.Error())
		return
	}
	msg = string(line)
	return
}

func (c *Client) Sendln(msg string) {
	fmt.Fprintf(*(c.Connection), msg+"\n")
}

func (c *Client) Close() {
	(*(c.Connection)).Close()
}
