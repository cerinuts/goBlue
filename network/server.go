/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"bufio"
	"net"

	"code.cerinuts.io/libs/goBlue/log"
)

//Server is a simple socketserver. Not doing anything, pretty useless atm. Seriously.
type Server struct {
	Address, Port string
	listener      net.Listener
}

//Create starts a echoserver on address:port. probably.
func (s *Server) Create(address, port string) {
	s.Address = address
	s.Port = port
	ln, err := net.Listen("tcp", address+":"+port)
	s.listener = ln
	if err != nil {
		println("Error creating server", err.Error())
	}
	for {
		println("waiting..")
		conn, err := ln.Accept()
		if err != nil {
			log.E(err)
			continue
		}
		println("accepted")
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			println(err.Error())
		}
		conn.Write([]byte(status))
		println("received", status)
		conn.Close()
		println("closed.")
	}
}
