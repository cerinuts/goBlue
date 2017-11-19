/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"encoding/json"
	"fmt"
	"gitlab.ceriath.net/libs/goBlue/log"
	"io"
	"net/http"
	"net/url"
)

//Dataserv offers an api to serve function calls and/or json responses to http.
//Can also wrap a database-like structure to a http-api.
type DataServ struct {
	ServeMux *http.ServeMux
}

//Data array
type DSJSONResponse struct {
	Data []DSJSONData `json:"data"`
}

//single data object
type DSJSONSingleResponse struct {
	Data DSJSONData `json:"data"`
}

//actual data
type DSJSONData struct {
	Id         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
	Type       string      `json:"type"`
}

//error array
type DSJSONErrors struct {
	Errs []DSJSONError `json:"errors"`
}

//single error
type DSJSONError struct {
	Code   int    `json:"code"`
	Source string `json:"source"`
	Detail string `json:"detail"`
}

//a custom error containing status
type customError struct {
	err    error
	status int
}

//a simple function thats called on get request
type DataServFunction func() error
//a function thats called on post request taking input
type DataServPostFunction func(url.Values) error
//a function thats called on patch request
type DataServPatchFunction func(*io.ReadCloser) error

//Creates a new Dataserv
func NewDataServ() *DataServ {
	ds := new(DataServ)
	ds.ServeMux = http.NewServeMux()
	return ds
}

//Start starts the http listener on host:port
func (ds *DataServ) Start(host, port string) {
	log.F(http.ListenAndServe(host+":"+port, ds.ServeMux))
}

//Register adds a route with a static response array on any method
func (ds *DataServ) Register(route string, jsr DSJSONResponse) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		c, err := json.Marshal(jsr)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(c)
		}
	})
}

//RegisterSingle adds a route with a single static response on any method
func (ds *DataServ) RegisterSingle(route string, jsr DSJSONSingleResponse) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		c, err := json.Marshal(jsr)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(c)
		}
	})
}

//RegisterGetFunction adds a route that invokes a function on any method
func (ds *DataServ) RegisterGetFunction(route string, fn DataServFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		err := fn()

		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		} else {
			handleError(err, w, http.StatusOK)
		}
	})
}

//RegisterWithPost adds a route that invokes a function on POST request with post form input or a static response array on any other method
func (ds *DataServ) RegisterWithPost(route string, jsr DSJSONResponse, fn DataServPostFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			}
			err = fn(r.PostForm)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusCreated)
				handleError(err, w, http.StatusOK)
			}
		} else {
			c, err := json.Marshal(jsr)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(c)
			}
		}
	})
}

//RegisterWithPost adds a route that invokes a function on POST request with post form input or a single static response on any other method
func (ds *DataServ) RegisterSingleWithPost(route string, jsr DSJSONSingleResponse, fn DataServPostFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		if r.Method == "POST" {
			err := r.ParseForm()
			println(r.PostForm.Encode())
			println(r.Body)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			}
			err = fn(r.PostForm)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusCreated)
				handleError(err, w, http.StatusOK)
			}
		} else {
			c, err := json.Marshal(jsr)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(c)
			}
		}
	})
}

//RegisterPostFunction adds a route that invokes a function on POST returning an error on any other method
func (ds *DataServ) RegisterPostFunction(route string, fn DataServPostFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		err := r.ParseForm()
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		}

		err = fn(r.PostForm)

		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "{}")
		}

	})
}

//RegisterWithPatch adds a route invoking a function on PATCH or returning a static response array on any other method
func (ds *DataServ) RegisterWithPatch(route string, jsr DSJSONResponse, fn DataServPatchFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		if r.Method == "PATCH" {
			err := fn(&r.Body)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				handleError(err, w, http.StatusOK)
			}
		} else {
			c, err := json.Marshal(jsr)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(c)
			}
		}
	})
}

//RegisterWithPatch adds a route invoking a function on PATCH or returning a single static response on any other method
func (ds *DataServ) RegisterSingleWithPatch(route string, jsr DSJSONSingleResponse, fn DataServPatchFunction) {
	ds.ServeMux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			ds.Headers(w, r)
			err := fn(&r.Body)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				handleError(err, w, http.StatusOK)
			}
		} else {
			ds.Headers(w, r)
			c, err := json.Marshal(jsr)
			if err != nil {
				handleError(err, w, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(c)
			}
		}
	})
}

//Headers writes general headers
func (ds *DataServ) Headers(rw http.ResponseWriter, req *http.Request) {
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

//handleError prepares and sends an error
func handleError(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	var es []customError
	es = append(es, customError{err, status})
	c, _ := json.Marshal(getError(es))
	w.Write(c)
}

//getError translates error into Dataserv json error object
func getError(ein []customError) DSJSONErrors {
	var jse DSJSONErrors
	for _, e := range ein {
		jserr := new(DSJSONError)
		jserr.Code = e.status
		if e.status < 400 {
			jserr.Detail = "ok"
		} else {
			jserr.Detail = e.err.Error()
		}
		jserr.Source = "goBlue/dataserv"
		jse.Errs = append(jse.Errs, *jserr)
	}
	return jse
}
