package main

import (
	lua "github.com/yuin/gopher-lua"
)

func pushErr(L *lua.LState, err error) {
	if err == nil {
		L.Push(lua.LNil)
		return
	}
	L.Push(lua.LString(err.Error()))
}
