package network

import (
	"net"
	"bufio"
)

type Server struct{
	Address, Port string
	listener net.Listener
}

func (s *Server) Create(address, port string){
	s.Address = address;
	s.Port = port;
	ln, err := net.Listen("tcp", address + ":" + port)
	s.listener = ln
	if(err != nil){
		println("Error creating server", err.Error())
	}
	for{
		println("waiting..")
		conn, err := ln.Accept()
		if(err != nil){
			//
		}
		println("accepted")
		status, err := bufio.NewReader(conn).ReadString('\n')
		if(err != nil){
			println(err.Error())
		}
		println("received", status)
		conn.Close()
		println("closed.")
	}
}
