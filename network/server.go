package network

import (
	"net"
	"bufio"
)

type Server struct{
	Address, Port string
	listener net.Listener
}

func (s Server) Create(address, port string){
	s.Address = address;
	s.Port = port;
	ln, err := net.Listen("tcp", address + ":" + port)
	s.listener = ln
	if(err != nil){
		println("Error creating server", err.Error())
	}
	for{
		conn, err := ln.Accept()
		if(err != nil){
			//
		}
		status, err := bufio.NewReader(conn).ReadString('\n')
		println("received", status)
		conn.Close()
	}
}

