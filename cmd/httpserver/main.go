package main

import (
	"github.com/jasonsites/gosk-api/internal/runtime"
)

func main() {
	runconf := &runtime.RunConfig{HTTPServer: true}
	runtime.NewRuntime(nil).Run(runconf)
}
