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
	"math/rand"
	"os"
	"time"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/util", "1", "0", "d"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GetRandomAlphanumericString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
