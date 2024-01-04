package main

import (
	"github.com/gidoichi/ical-converter/di"
)

func main() {
	server := di.DI()
	server.Run()
}
