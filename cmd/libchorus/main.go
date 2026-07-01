// Package main provides a C-shared library wrapper around the Chorus build
// engine, built with `go build -buildmode=c-shared`.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"sync"
	"unsafe"

	"z3nnix/chorus/pkg/chorus"
)

var (
	mu        sync.Mutex
	lastError string
)

func setLastError(err error) {
	mu.Lock()
	defer mu.Unlock()
	if err != nil {
		lastError = err.Error()
	} else {
		lastError = ""
	}
}

//export ChorusBuild
func ChorusBuild(configPath *C.char, targets *C.char) C.int {
	path := C.GoString(configPath)
	if path == "" {
		path = "chorus.build"
	}

	config, err := chorus.LoadConfig(path)
	if err != nil {
		setLastError(err)
		return -1
	}

	targetList := strings.Fields(C.GoString(targets))
	if err := config.Build(targetList...); err != nil {
		setLastError(err)
		return -1
	}

	setLastError(nil)
	return 0
}

//export ChorusLastError
func ChorusLastError() *C.char {
	mu.Lock()
	defer mu.Unlock()
	return C.CString(lastError)
}

//export ChorusFreeString
func ChorusFreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {}
