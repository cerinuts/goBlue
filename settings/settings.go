/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//settings is a simple jsonfile-to-struct api to save settings
package settings

import (
	"encoding/json"
	"github.com/ceriath/goBlue/log"
	"os"
	"path/filepath"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/settings", "0", "1", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//Reads a json-config file to any struct
func ReadJsonConfig(filename string, config interface{}) error {
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

//writes json-config from any struct
func WriteJsonConfig(filename string, config interface{}) error {
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
