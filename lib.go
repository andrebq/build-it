package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rumlang/rum/runtime"
)

type (
	lib struct {
		canRunCmd bool
	}
)

func (l *lib) LoadLib(ctx *runtime.Context) {
	if l.canRunCmd {
		ctx.SetFn("buildit.run", l.runCmd)
	}
	ctx.SetFn("buildit.fatal", l.fatal)
	ctx.SetFn("concat", l.concat)
	ctx.SetFn("concat-with", l.concatWith)
}

func (l *lib) concat(args []interface{}) string {
	return strings.Join(toStringArray(args), "")
}

func (l *lib) concatWith(sep string, args []interface{}) string {
	return strings.Join(toStringArray(args), sep)
}

func (l *lib) fatal(code int, reason interface{}) {
	log.Print(-1, fmt.Sprintf("%v", reason))
	os.Exit(code)
}

func (l *lib) runCmd(binary string, args ...interface{}) (string, error) {
	cmd := exec.Command(binary, toStringArray(args)...)
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
