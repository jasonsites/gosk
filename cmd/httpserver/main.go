package main

import (
	"github.com/jasonsites/gosk/internal/runtime"
)

func main() {
	runconf := &runtime.RunConfig{Entry: "http"}
	runtime.NewRuntime(nil).Run(runconf)
}
