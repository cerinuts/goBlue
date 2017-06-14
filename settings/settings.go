package settings

import (
	"encoding/json"
	"github.com/ceriath/goBlue/log"
	"os"
)

//Reads a json-config file to any struct
func ReadJsonConfig(filename string, config interface{}) error {
	file, _ := os.Open(filename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		log.E(err)
		return err
	}
	return nil
}

//writes json-config from any struct
func WriteJsonConfig(filename string, config interface{}) error {
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
