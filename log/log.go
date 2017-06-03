/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const LevelPanic, LevelFatal, LevelError, LevelInfo, LevelDebug int = 1, 2, 3, 4, 5
const LogfileBehaviourDaily, LogfileBehaviourAll int = 1, 0
const errorlog string = "error.log"

var CurrentLevel = 3
var CurrentLogFileBehaviour = 0
var PrintToStderr, PrintToStdout, PrintToFile = true, false, true
var Path, Logfilename = ".", "log.log"

//panic
func P(a ...interface{}) {
	if CurrentLevel < LevelPanic {
		return
	}
	appended := fmt.Sprintf("[PANIC] - [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), a)
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	if PrintToStderr {
		fmt.Fprintln(os.Stderr, appended)
	}
	printToFile(appended, true)
	panic(appended)
//	os.Exit(1)
}

//fatal
func F(a ...interface{}) {
	if CurrentLevel < LevelFatal {
		return
	}
	appended := fmt.Sprintf("[FATAL] - [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), a)
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	if PrintToStderr {
		fmt.Fprintln(os.Stderr, appended)
	}
	printToFile(appended, true)
}

//error
func E(a ...interface{}) {
	if CurrentLevel < LevelError {
		return
	}
	appended := fmt.Sprintf("[ERROR] - [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), a)
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	if PrintToStderr {
		fmt.Fprintln(os.Stderr, appended)
	}
	printToFile(appended, true)
}

//info
func I(a ...interface{}) {
	if CurrentLevel < LevelInfo {
		return
	}
	appended := fmt.Sprintf("[INFO] - [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), a)
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	printToFile(appended, false)
}

//debug
func D(a ...interface{}) {
	if CurrentLevel < LevelDebug {
		return
	}
	appended := fmt.Sprintf("[DEBUG] - [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), a)
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	printToFile(appended, false)
}

func printToFile(line string, isError bool) {
	filename := Logfilename
	if CurrentLogFileBehaviour == LogfileBehaviourDaily {
		filename = time.Now().Format("2006_01_02.log")
	}
	if PrintToFile {
		f, err := os.OpenFile(filepath.Join(Path, filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(line + "\n"); err != nil {
			panic(err)
		}
	}
	if isError {
		f, err := os.OpenFile(filepath.Join(Path, errorlog), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(line + "\n"); err != nil {
			panic(err)
		}
	}
}
