/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package log offers various logging options
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/log", "0", "1", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

const LevelPanic, LevelFatal, LevelError, LevelInfo, LevelDebug int = 1, 2, 3, 4, 5

//Defines if the logger logs everything to one file or creates one for each day (useful for logrotate)
const LogfileBehaviourDaily, LogfileBehaviourAll int = 1, 0
const errorlog string = "error.log"

//CurrentLevel defines the current loglevel
var CurrentLevel = 4

//CurrentLogFileBehaviour defines current logfile behaviour, defaults to LogfileBehaviourAll
var CurrentLogFileBehaviour = LogfileBehaviourAll

//Define which outputs are used, note that panic always prints to stderr and error file, while fatal and error always
//print to error file, if printToFile is true those additionally printed to normal log
var PrintToStderr, PrintToStdout, PrintToFile = true, false, true

//Logpath and filename for the normal log
var Path, Logfilename = ".", "log.log"

//P logs panic, NOTE that this will throw a panic at the end!
func P(a ...interface{}) {
	if CurrentLevel < LevelPanic {
		return
	}
	_, file, no, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	appended := fmt.Sprintf("[PANIC] - [%s#%d] - [%s] - [%s] - %s", file, no, time.Now().Format("2006-01-02 15:04:05"), a[0], a[1:])
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	fmt.Fprintln(os.Stderr, appended)

	printToFile(appended, true)
	panic(appended)
	//	os.Exit(1)
}

//F logs fatal
func F(a ...interface{}) {
	if CurrentLevel < LevelFatal {
		return
	}
	_, file, no, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	appended := fmt.Sprintf("[FATAL] - [%s#%d] - [%s] - [%s] - %s", file, no, time.Now().Format("2006-01-02 15:04:05"), a[0], a[1:])
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	if PrintToStderr {
		fmt.Fprintln(os.Stderr, appended)
	}
	printToFile(appended, true)
}

//E logs error
func E(a ...interface{}) {
	if CurrentLevel < LevelError {
		return
	}
	_, file, no, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	appended := fmt.Sprintf("[ERROR] - [%s#%d] - [%s] - [%s] - %s", file, no, time.Now().Format("2006-01-02 15:04:05"), a[0], a[1:])
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	if PrintToStderr {
		fmt.Fprintln(os.Stderr, appended)
	}
	printToFile(appended, true)
}

//I logs info
func I(a ...interface{}) {
	if CurrentLevel < LevelInfo {
		return
	}
	_, file, no, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	appended := fmt.Sprintf("[INFO] - [%s#%d] - [%s] - [%s] - %s", file, no, time.Now().Format("2006-01-02 15:04:05"), a[0], a[1:])
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	printToFile(appended, false)
}

//D logs debug
func D(a ...interface{}) {
	if CurrentLevel < LevelDebug {
		return
	}
	_, file, no, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	appended := fmt.Sprintf("[DEBUG] - [%s#%d] - [%s] - [%s] - %s", file, no, time.Now().Format("2006-01-02 15:04:05"), a[0], a[1:])
	if PrintToStdout {
		fmt.Fprintln(os.Stdout, appended)
	}
	printToFile(appended, false)
}

//print to logfile and errorlog if applicable
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
