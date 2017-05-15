package network

import(
	"net"
	"fmt"
	"bufio"
	"bytes"
)

type Client struct{
	TargetIP, TargetPort string
	Connection *net.Conn
}

func (c *Client) Connect(ip string, port string) (err error){
	c.TargetIP = ip
	c.TargetPort = port
	conn, err := net.Dial("tcp", ip + ":" + port)
	c.Connection = &conn
	if(err != nil){
		println("Error connecting", err.Error())
		return
	}
	return
}

func (c *Client) Recv() (msg string, err error){
	buff := bytes.NewBufferString(msg)
	continueLine := true
	for(continueLine){
		line, isPref, erro := bufio.NewReader(*(c.Connection)).ReadLine();
		continueLine = isPref
		if( erro != nil){
			println("Error reading from Socket", erro.Error())
			err = erro
			return
		}
		buff.Write(line)
	}
	msg = buff.String()
	return
}

func (c *Client) Sendln(msg string){
	fmt.Fprintf(*(c.Connection), msg + "\n")
}

func (c *Client) Close(){
	(*(c.Connection)).Close()
}