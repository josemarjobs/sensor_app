package main

import (
	"fmt"

	"github.com/josemarjobs/sensors_app/coordinator"
)

func main() {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
