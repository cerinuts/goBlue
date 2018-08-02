/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package util offers misc utility functions
package util

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

//AppName is the name of the application
const AppName string = "goBlue/util"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "1"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//SaveCopy savely copies a file from src to dst
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

//WaitForEnter for userinput to continue
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

//GetRandomAlphanumericString returns a random alphanumeric string of length n
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

//NormalizeDurationStringHMS formats a time d into a string containing "x hours x minutes x seconds"
//while it removes 0 valued parts
func NormalizeDurationStringHMS(d time.Duration) string {
	seconds := int(d.Seconds())
	minutes := seconds / 60
	seconds = seconds - minutes*60
	hours := minutes / 60
	minutes = minutes - hours*60

	str := fmt.Sprintf("%d hours %d minutes %d seconds", hours, minutes, seconds)

	if strings.HasSuffix(str, " 0 seconds") {
		str = str[:len(str)-9]
	}
	str = strings.TrimPrefix(str, "0 hours")
	str = strings.Replace(str, " 0 minutes", " ", 1)
	return str
}

//NormalizeDurationStringYMDHMS formats the output of TimeDifferenceYMDHMS into a string
func NormalizeDurationStringYMDHMS(year, month, day, hour, min, sec int) string {
	var res string
	if year > 0 {
		res = res + strconv.Itoa(year) + " years "
	}
	if month > 0 {
		res = res + strconv.Itoa(month) + " months "
	}
	if day > 0 {
		res = res + strconv.Itoa(day) + " days "
	}
	if hour > 0 {
		res = res + strconv.Itoa(hour) + " hours "
	}
	if min > 0 {
		res = res + strconv.Itoa(min) + " minutes "
	}
	if sec > 0 {
		res = res + strconv.Itoa(sec) + " seconds "
	}
	return strings.TrimRight(res, " ")
}

//TimeDifferenceYMDHMS calculates the difference between times a and b and outputs a the difference as integers
func TimeDifferenceYMDHMS(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1
	hour = h2 - h1
	min = m2 - m1
	sec = s2 - s1

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

//OpenWebsiteInDefaultBrowser uses the system api to open a website in the default browser
func OpenWebsiteInDefaultBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"rundll32", "url.dll,FileProtocolHandler"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

//RemoveIndexFromStringSlice removes a string on index i from a slice
func RemoveIndexFromStringSlice(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

//RemoveFromStringSlice removes the first matching string from a slice
func RemoveFromStringSlice(slice []string, el string) []string {
	var idx int
	for i, e := range slice {
		if e == el {
			idx = i
			break
		}
	}
	return RemoveIndexFromStringSlice(slice, idx)
}
