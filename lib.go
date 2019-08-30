package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type (
	lib struct {
		canRunCmd bool
	}
)

func (l lib) load(L *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"fatal": func(L *lua.LState) int {
			ec := L.CheckInt(1)
			msg := L.CheckString(2)

			log.Print(-1, msg)
			os.Exit(ec)
			return 2
		},
	}
	if l.canRunCmd {
		exports["run"] = func(L *lua.LState) int {
			top := L.GetTop()
			switch top {
			case 0:
				L.Push(lua.LNil)
				L.Push(lua.LString("missing binary"))
				return 2
			}
			var args []string
			bin := L.CheckString(1)
			for i := 2; i <= top; i++ {
				args = append(args, lua.LVAsString(L.Get(i)))
			}
			output, err := l.runCmd(bin, args...)
			L.Push(lua.LString(output))
			pushErr(L, err)
			return 2
		}
	}
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

func (l *lib) fatal(code int, reason interface{}) {
	log.Print(-1, fmt.Sprintf("%v", reason))
	os.Exit(code)
}

func (l *lib) runCmd(binary string, args ...string) (string, error) {
	cmd := exec.Command(binary, args...)
	cmd.Stderr = os.Stderr
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	err := cmd.Run()
	return strings.TrimSpace(stdoutBuf.String()), err
}

func toStringArray(args []interface{}) []string {
	ret := make([]string, 0, len(args))
	for _, v := range args {
		ret = append(ret, fmt.Sprintf("%v", v))
	}
	return ret
}
