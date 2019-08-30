package main

import (
	"strconv"
	"strings"
)

func getPathDeep(path string, cur interface{}) interface{} {
	split := strings.Split(path, ".")
	for _, v := range split {
		cur = getPathFrom(v, cur)
	}
	return cur
}

func getPathFrom(name string, cur interface{}) interface{} {
	if cur == nil {
		return cur
	}

	switch cur := cur.(type) {
	case []interface{}:
		return cur[mustInt(name)]
	case map[string]interface{}:
		return cur[name]
	}
	return cur
}

func mustInt(n string) int {
	val, err := strconv.Atoi(n)
	if err != nil {
		panic(err.Error())
	}
	return val
}
