package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/rumlang/rum/runtime"
)

type (
	httpLib struct {
	}
)

func (l httpLib) LoadLibrary(ctx *runtime.Context) {
	ctx.SetFn("http.get", l.httpFn("GET"))
	ctx.SetFn("http.post", l.httpFn("POST"))
	ctx.SetFn("http.status", l.status)
	ctx.SetFn("http.body", l.body)
	// ctx.SetFn("http.head", l.httpFn("HEAD"))
	// ctx.SetFn("http.put", l.httpFn("PUT"))
	// ctx.SetFn("http.delete", l.httpFn("DELETE"))
}

func (l httpLib) status(args []interface{}) (int, error) {
	switch len(args) {
	case 2:
	default:
		return 0, errors.New("no args")
	}
	return args[1].(int), nil
}

func (l httpLib) body(args []interface{}) (string, error) {
	switch len(args) {
	case 2:
	default:
		return "", errors.New("no args")
	}
	return args[0].(string), nil
}

func (l httpLib) httpFn(method string) func(args ...string) ([]interface{}, error) {
	return func(args ...string) ([]interface{}, error) {
		var url, body string

		switch len(args) {
		case 0:
			return nil, errors.New("missing url")
		case 1:
			url = args[0]
		case 2:
			url, body = args[0], args[1]
		}

		s := sling.New().Base(url)
		switch method {
		case "GET":
			return process(s.Get(""))
		case "POST":
			return process(s.Body(bytes.NewBufferString(body)).Post(""))
		}
		return nil, errors.New("invalid method")
	}
}

func process(s *sling.Sling) ([]interface{}, error) {
	req, err := s.Request()
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return []interface{}{string(buf), resp.StatusCode}, nil
}
