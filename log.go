package main

import (
	logp "log"
)

type (
	logT int
)

var (
	log = logT(0)
)

func (l logT) SetPrefix(prefix string) {
	logp.SetPrefix(prefix)
}

func (l logT) Print(level int, args ...interface{}) {
	if int(l) < level {
		return
	}
	logp.Print(args)
}

func (l logT) Fatal(args ...interface{}) {
	logp.Fatal(args)
}
