package main

import (
	"github.com/ceriath/goBlue/network"
	"time"
	//"github.com/ceriath/goBlue/network/client"
)
	
func main(){
	s := new(network.Server)
	go s.Create("localhost", "12345")
	time.Sleep(2 * time.Second)
	c := new(network.Client)
	err := c.Connect("localhost", "12345")
	if(err != nil){
		println("Error connecting client")
	}
	c.Sendln("hi")
	println("sent.")
	time.Sleep(2 * time.Second)
	
}
