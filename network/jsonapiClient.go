package network

import (
	"bytes"
	"encoding/json"
	"github.com/ceriath/goBlue/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

type JsonApiClient struct {
}

type JsonError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type JsonError2 struct {
	Status struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"status"`
}

func (jso *JsonError) String() string {
	return strconv.Itoa(jso.Status) + "-" + jso.Error + "-" + jso.Message
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

	if res.StatusCode == 200 {
		json.Unmarshal(body, &response)
		return nil, nil
	} else {
		//try if its an error
		jsoErr := new(JsonError)
		marshErr := json.Unmarshal(body, &jsoErr)
		if marshErr == nil {
			return jsoErr, nil
		}
		//try if its an error2
		jsoErr2 := new(JsonError2)
		marshErr = json.Unmarshal(body, &jsoErr2)
		if marshErr == nil {
			jsoErr.Status = jsoErr2.Status.StatusCode
			jsoErr.Message = jsoErr2.Status.Message
			jsoErr.Error = jsoErr2.Status.Message
			return jsoErr, nil
		}
		//otherwise some error
		log.E(marshErr)
		return nil, marshErr
	}
}
