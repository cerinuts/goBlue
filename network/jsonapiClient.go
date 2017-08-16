package network

import (
	"bytes"
	"encoding/json"
	"github.com/ceriath/goBlue/log"
	"io/ioutil"
	"net/http"
)

type JsonApiClient struct {
}

type JsonError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (jso *JsonError) String() string {
	return string(jso.Status) + "-" + jso.Error + "-" + jso.Message
}

func (jac *JsonApiClient) Request(url string, header map[string]string, response interface{}) (*JsonError, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)
}

func (jac *JsonApiClient) Delete(url string, header map[string]string, response interface{}) (*JsonError, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.E(err)
		return nil, err
	}

	return jac.runRequest(req, header, response)
}

func (jac *JsonApiClient) Put(url string, header map[string]string, data interface{}, response interface{}) (*JsonError, error) {

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

func (jac *JsonApiClient) Post(url string, header map[string]string, data interface{}, response interface{}) (*JsonError, error) {

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

func (jac *JsonApiClient) runRequest(req *http.Request, header map[string]string, response interface{}) (*JsonError, error) {
	cli := new(http.Client)

	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, getErr := cli.Do(req)
	if getErr != nil {
		log.E(getErr)
		return nil, getErr
	}

	if res.StatusCode == 204 {
		return nil, nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.E(readErr)
		return nil, readErr
	}

	//try to unmarshal
	jsonErr := json.Unmarshal(body, &response)

	if jsonErr != nil {
		log.I(jsonErr)
		//try if its an error
		jsoErr := new(JsonError)
		marshErr := json.Unmarshal(body, &jsoErr)
		if marshErr == nil {
			return jsoErr, nil
		}
		log.E(jsonErr)
		return nil, jsonErr
	}

	return nil, nil
}
