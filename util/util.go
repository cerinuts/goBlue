/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Offers misc utility functions
package util

import (
	"bufio"
	"fmt"
	"github.com/ceriath/goBlue/log"
	"io"
	"os"
)

//Savely copies a file from src to dst
func SaveCopy(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		log.E(err)
		return err
	}
	defer file.Close()
	nf, err2 := os.Create(dst)
	if err2 != nil {
		log.E(err2)
		return err2
	}
	defer nf.Close()
	nf.Close()
	file.Close()
	_, err3 := io.Copy(nf, file)
	if err3 != nil {
		log.E(err2)
		return err2
	}
	return nil
}

//Waits for userinput to continue
func WaitForEnter() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
