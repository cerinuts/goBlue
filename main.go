package main

import (
	"github.com/ceriath/goBlue/network"
	//"github.com/ceriath/goBlue/network/client"
)
	
func main(){
	s := network.Server{}
	go s.Create("localhost", "12345")
	
	c := network.Client{}
	err := c.Connect("localhost", "12345")
	if(err != nil){
		println("Error connecting client")
	}
	//c.Send("hi")
	
}
