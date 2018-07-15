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

	"gitlab.ceriath.net/libs/goBlue/log"
)

const PARSE_CODE_ERR = -1
const PARSE_CODE_OK = 0
const PARSE_CODE_DISPATCH = 1

type EventsourceClient struct {
	Stream *EventStream
}

type Event struct {
	Id      string
	Name    string
	Payload string
}

type EventStream struct {
	client     *http.Client
	req        *http.Request
	timeout    time.Duration
	closed     bool
	EventQueue chan Event
	lastEvent  string
}

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
		errors.New(string(code))
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

func (stream *EventStream) start(reader io.ReadCloser) {
	defer reader.Close()

	stream.recv(reader)

	stream.reconnect()
}

func (stream *EventStream) recv(reader io.ReadCloser) {

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
		case PARSE_CODE_OK:
			{
				continue
			}
		case PARSE_CODE_DISPATCH:
			{
				go dispatch(stream, ev)
				ev = new(Event)
			}
		case PARSE_CODE_ERR:
			{
				continue
			}
		}
	}
}

func parse(ev *Event, evS *EventStream, line string) (code int) {

	line = strings.TrimSpace(line)

	if len(line) < 1 {
		return PARSE_CODE_DISPATCH
	}
	if strings.HasPrefix(line, ":") {
		return PARSE_CODE_OK
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
			return PARSE_CODE_OK
		}
	case "data":
		{
			if len(ev.Payload) > 0 {
				ev.Payload += "\n"
			}
			ev.Payload += value
			return PARSE_CODE_OK
		}
	case "id":
		{
			evS.lastEvent = value
			ev.Id = value
			return PARSE_CODE_OK
		}
	case "retry":
		{
			to, err := strconv.Atoi(value)
			if err != nil {
				evS.timeout = time.Millisecond * time.Duration(to)
			}
			return PARSE_CODE_OK
		}
	}
	return PARSE_CODE_OK

}

func dispatch(evS *EventStream, ev *Event) {
	if len(ev.Payload) < 1 {
		return
	}

	evS.EventQueue <- *ev
}
