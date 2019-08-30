package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
	lua "github.com/yuin/gopher-lua"
)

type (
	httpLib struct {
	}
)

func (l httpLib) load(L *lua.LState) int {
	createForMethod := func(method string) func(*lua.LState) int {
		return func(L *lua.LState) int {
			var url, body string
			switch L.GetTop() {
			case 0:
				L.Push(lua.LNil)
				L.Push(lua.LNil)
				L.Push(lua.LString("no body"))
				return 3
			case 1:
				url = L.CheckString(1)
			case 2:
				url = L.CheckString(1)
				body = L.CheckString(2)
			}
			body, status, err := l.httpCall(method, url, body)
			L.Push(lua.LString(body))
			L.Push(lua.LNumber(float64(status)))
			if err != nil {
				L.Push(lua.LString(err.Error()))
			} else {
				L.Push(lua.LNil)
			}
			return 3
		}
	}

	methods := map[string]lua.LGFunction{
		"get":  createForMethod("GET"),
		"post": createForMethod("POST"),
	}
	mod := L.SetFuncs(L.NewTable(), methods)
	L.Push(mod)
	return 1
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

func (l httpLib) json(args []interface{}) (interface{}, error) {
	switch len(args) {
	case 2:
	default:
		println("no args")
		return nil, errors.New("no args")
	}

	var out map[string]interface{}
	println("args[0].(string)", args[0].(string))
	err := json.Unmarshal([]byte(args[0].(string)), &out)
	if err != nil {
		println("errors!")
		return nil, err
	}

	return out, nil
}

func (l httpLib) httpCall(method, url, body string) (string, int, error) {
	s := sling.New().Base(url)
	switch method {
	case "GET":
		return process(s.Get(""))
	case "POST":
		return process(s.Body(bytes.NewBufferString(body)).Post(""))
	}
	return "", 0, errors.New("invalid method")
}

func process(s *sling.Sling) (string, int, error) {
	req, err := s.Request()
	if err != nil {
		return "", 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	return string(buf), resp.StatusCode, nil
}
