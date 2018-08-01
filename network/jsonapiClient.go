/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

//JSONAPIClient offers a simple json api client interface
type JSONAPIClient struct {
}

//JSONError contains a simple JSON error message
type JSONError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//JSONError2 contains another common JSON error message
type jsonError2 struct {
	Status struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"status"`
}

//String converts a json error to pretty loggable/printable string
func (jso *JSONError) String() string {
	return strconv.Itoa(jso.Status) + "-" + jso.Error + "-" + jso.Message
}

//Request calls url with GET and sets header. Tries to parse repsonse into any struct, returns jsonerror if request returned one
//or error on internal errors
func (jac *JSONAPIClient) Request(url string, header map[string]string, response interface{}) (*JSONError, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)
}

//Delete calls url with DELETE and sets header. Tries to parse repsonse into any struct, returns jsonerror if request returned one
//or error on internal errors
func (jac *JSONAPIClient) Delete(url string, header map[string]string, response interface{}) (*JSONError, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)
}

//Put calls url with PUT with given data and sets header. Tries to parse repsonse into any struct, returns jsonerror if request returned one
//or error on internal errors
func (jac *JSONAPIClient) Put(url string, header map[string]string, data interface{}, response interface{}) (*JSONError, error) {

	body, marshErr := json.Marshal(data)
	if marshErr != nil {
		return nil, marshErr
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)

}

//Post calls url with POST with given data and sets header. Tries to parse repsonse into any struct, returns jsonerror if request returned one
//or error on internal errors
func (jac *JSONAPIClient) Post(url string, header map[string]string, data interface{}, response interface{}) (*JSONError, error) {

	body, marshErr := json.Marshal(data)
	if marshErr != nil {
		return nil, marshErr
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)

}

//runRequest actually runs the request prepared by functions above
func (jac *JSONAPIClient) runRequest(req *http.Request, header map[string]string, response interface{}) (*JSONError, error) {
	cli := &http.Client{
		Timeout: time.Second * 10,
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, getErr := cli.Do(req)
	if getErr != nil {
		log.E(getErr)
		return nil, getErr
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil, nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.E(readErr)
		return nil, readErr
	}

	if res.StatusCode == 200 {
		return nil, json.Unmarshal(body, &response)
	}
	//try if its an error
	jsoErr := new(JSONError)
	marshErr := json.Unmarshal(body, &jsoErr)
	if marshErr == nil {
		return jsoErr, nil
	}
	//try if its an error2
	jsoErr2 := new(jsonError2)
	marshErr = json.Unmarshal(body, &jsoErr2)
	if marshErr == nil {
		jsoErr.Status = jsoErr2.Status.StatusCode
		jsoErr.Message = jsoErr2.Status.Message
		jsoErr.Error = jsoErr2.Status.Message
		return jsoErr, nil
	}
	//otherwise some error
	log.E(marshErr, res.StatusCode, string(body))
	return nil, marshErr

}
