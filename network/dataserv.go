package network

import (
	"encoding/json"
	"fmt"
	"github.com/ceriath/goBlue/log"
	"net/http"
	"net/url"
	"io"
)

type DataServ struct {
	r *http.ServeMux
}

type JSONResponse struct {
	Data []JSONData `json:"data"`
}

type JSONSingleResponse struct {
	Data JSONData `json:"data"`
}

type JSONData struct {
	Id         string         `json:"id"`
	Attributes interface{} `json:"attributes"`
	Type       string      `json:"type"`
}

type JSONErrors struct {
	Errs []JSONError `json:"errors"`
}

type JSONError struct {
	Code int `json:"code"`
	Source string `json:"source"`
	Detail string `json:"detail"`
}

type customError struct {
	err    error
	status int
}

type DataServFunction func() error
type DataServPostFunction func(url.Values) error
type DataServPatchInputFunctionPointer *func()(io.ReadCloser, error)
type DataServPatchFunction func(*io.ReadCloser) error

func NewDataServ() *DataServ {
	ds := new(DataServ)
	ds.r = http.NewServeMux()
	return ds
}

func (ds *DataServ) Start(port string) {
	log.F(http.ListenAndServe(":"+port, ds.r))
}

func (ds *DataServ) Register(route string, jsr JSONResponse) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterSingle(route string, jsr JSONSingleResponse) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterGetFunction(route string, fn DataServFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ds.Headers(w, r)
		err := fn()

		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
		} else {
			handleError(err, w, http.StatusOK)
		}
	})
}

func (ds *DataServ) RegisterWithPost(route string, jsr JSONResponse, fn DataServPostFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterSingleWithPost(route string, jsr JSONSingleResponse, fn DataServPostFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterPostFunction(route string, fn DataServPostFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterWithPatch(route string, jsr JSONResponse, fn DataServPatchFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) RegisterSingleWithPatch(route string, jsr JSONSingleResponse, fn DataServPatchFunction) {
	ds.r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
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

func (ds *DataServ) Headers(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
}

func handleError(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	var es []customError
	es = append(es, customError{err, status})
	c, _ := json.Marshal(getError(es))
	w.Write(c)
}

func getError(ein []customError) JSONErrors {
	var jse JSONErrors
	for _, e := range ein {
		jserr := new(JSONError)
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
