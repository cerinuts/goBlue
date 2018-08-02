/*
Copyright (c) 2018 ceriath
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

//AppName is the name of the application
const AppName string = "goBlue/log"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "1"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//LevelPanic program cannot continue. Will actually cause a panic!
const LevelPanic int = 1

//LevelFatal fatal error, program can at least exit gracefully
const LevelFatal int = 2

//LevelError any error occurred
const LevelError int = 3

//LevelInfo just some information
const LevelInfo int = 4

//LevelDebug debug information
const LevelDebug int = 5

//LogfileBehaviourDaily Defines if the logger logs everything to one file for each day (useful for logrotate)
const LogfileBehaviourDaily int = 1

//LogfileBehaviourAll Defines if the logger logs everything to one single file
const LogfileBehaviourAll int = 0

const errorlog string = "error.log"

//CurrentLevel defines the current loglevel
var CurrentLevel = 4

//CurrentLogFileBehaviour defines current logfile behaviour, defaults to LogfileBehaviourAll
var CurrentLogFileBehaviour = LogfileBehaviourAll

//PrintToStderr defines which outputs are used, note that panic always prints to stderr and ignores this setting, while Info and Debug will never print to stderr
var PrintToStderr = true

//PrintToStdout defines which outputs are used, note that everything will be printed to stdout while nothing is printed by default
var PrintToStdout = false

//PrintToFile defines which outputs are used, note that panic, fatal and error always
//prints to error.log, if printToFile is true those are printed to normal log file additionally
var PrintToFile = true

//Path for the default log
var Path = "."

//Logfilename name of the default log
var Logfilename = "log.log"

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
