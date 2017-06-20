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
	"time"
	)

type Client struct {
	TargetIP, TargetPort string
	Connection           *net.Conn
	Reader				*bufio.Reader
}

func (c *Client) Connect(ip string, port string) (err error) {
	c.TargetIP = ip
	c.TargetPort = port
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Duration(10 * time.Second))
	c.Connection = &conn
	reader := bufio.NewReader(*(c.Connection))
	c.Reader = reader
	if err != nil {
		println("Error connecting", err.Error())
		return
	}
	return
}

func (c *Client) Recv() (msg string, err error){
	line, _, err := (*c.Reader).ReadLine()
	if err != nil {
		println("Error reading from Socket", err.Error())
		return
	}
	msg = string(line)
	return
}

func (c *Client) Sendln(msg string) {
	fmt.Fprintf(*(c.Connection), msg+"\n")
}

func (c *Client) Close() (err error){
	err = (*(c.Connection)).Close()
	return
}
