/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

const parseCodeErr = -1
const parseCodeOk = 0
const parseCodeDispatch = 1

//EventsourceClient is a simple Client for Eventsource streams
type EventsourceClient struct {
	Stream *EventStream
}

//Event contains information about a single event
type Event struct {
	ID      string
	Name    string
	Payload string
}

//EventStream contains a channel with events
type EventStream struct {
	client     *http.Client
	req        *http.Request
	timeout    time.Duration
	closed     bool
	EventQueue chan Event
	lastEvent  string
}

//Subscribe subscribes to the given Eventsource-Server URL
func (ec *EventsourceClient) Subscribe(url string) (*EventStream, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	cli := &http.Client{
		//Timeout: time.Second * 10,
	}

	stream := &EventStream{
		client:     cli,
		req:        req,
		timeout:    5 * time.Second,
		EventQueue: make(chan Event),
	}

	ec.Stream = stream

	reader, err := stream.connect()
	if err != nil {
		return nil, err
	}
	go stream.start(reader)
	return stream, nil
}

//Close closes the stream
func (ec *EventsourceClient) Close() {
	if ec.Stream.closed {
		return
	}

	ec.Stream.closed = true
	close(ec.Stream.EventQueue)
}

func (stream *EventStream) connect() (r io.ReadCloser, err error) {
	var resp *http.Response
	stream.req.Header.Set("Accept", "text/event-stream")

	resp, err = stream.client.Do(stream.req)
	if err != nil {
		log.E(err)
		return
	}
	if resp.StatusCode != 200 {
		code, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(code))
	}
	r = resp.Body
	return
}

func (stream *EventStream) reconnect() {
	for {
		time.Sleep(stream.timeout)
		if stream.closed {
			return
		}

		reader, err := stream.connect()
		if err == nil {
			go stream.start(reader)
			return
		}
	}
}

func (stream *EventStream) start(reader io.Reader) {
	stream.recv(reader)

	stream.reconnect()
}

func (stream *EventStream) recv(reader io.Reader) {

	buffer := bufio.NewReader(reader)
	ev := new(Event)

	for {
		line, err := buffer.ReadString('\n')
		if stream.closed {
			return
		}
		if err != nil {
			log.E(err)
			return
		}

		switch parse(ev, stream, line) {
		case parseCodeOk:
			{
				continue
			}
		case parseCodeDispatch:
			{
				go dispatch(stream, ev)
				ev = new(Event)
			}
		case parseCodeErr:
			{
				continue
			}
		}
	}
}

func parse(ev *Event, evS *EventStream, line string) (code int) {

	line = strings.TrimSpace(line)

	if len(line) < 1 {
		return parseCodeDispatch
	}
	if strings.HasPrefix(line, ":") {
		return parseCodeOk
	}

	var field string
	var value string

	if strings.Contains(line, ":") {
		split := strings.SplitN(line, ":", 2)
		field = split[0]
		value = strings.TrimPrefix(strings.TrimSuffix(split[1], "\n"), " ")
	} else {
		field = strings.TrimSuffix(line, "\n")
	}

	switch field {
	case "event":
		{
			ev.Name = value
			return parseCodeOk
		}
	case "data":
		{
			if len(ev.Payload) > 0 {
				ev.Payload += "\n"
			}
			ev.Payload += value
			return parseCodeOk
		}
	case "id":
		{
			evS.lastEvent = value
			ev.ID = value
			return parseCodeOk
		}
	case "retry":
		{
			to, err := strconv.Atoi(value)
			if err != nil {
				evS.timeout = time.Millisecond * time.Duration(to)
			}
			return parseCodeOk
		}
	}
	return parseCodeOk

}

func dispatch(evS *EventStream, ev *Event) {
	if len(ev.Payload) < 1 {
		return
	}

	evS.EventQueue <- *ev
}
