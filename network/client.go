package network

import(
	"net"
	"fmt"
)

type Client struct{
	TargetIP, TargetPort string
	connection net.Conn
}

func (c Client) Connect(ip string, port string) (err error){
	c.TargetIP = ip
	c.TargetPort = port
	conn, err := net.Dial("tcp", ip + ":" + port)
	c.connection = conn
	if(err != nil){
		println("Error connecting", err.Error())
		return
	}
	return
}

func (c Client) Send(msg string){
	fmt.Fprintf(c.connection, msg)
}