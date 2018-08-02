/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package settings is a simple jsonfile-to-struct api to save settings
package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"code.cerinuts.io/libs/goBlue/log"
)

//AppName is the name of the application
const AppName string = "goBlue/settings"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibilty
const VersionMinor string = "1"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//ReadJSONConfig Reads a json-config file to any struct
func ReadJSONConfig(filename string, config interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.E(err)
		return err
	}
	return nil
}

//WriteJSONConfig writes json-config from any struct
func WriteJSONConfig(filename string, config interface{}) error {
	dir, _ := filepath.Split(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.E(err)
		return err
	}

	file, err := os.Create(filename + ".tmp")
	if err != nil {
		log.E(err)
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err1 := encoder.Encode(&config)
	if err1 != nil {
		log.E(err1)
		return err
	}
	file.Close()
	err2 := os.Rename(filename+".tmp", filename)
	return err2
}
