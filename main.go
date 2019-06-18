package main

import (
	"log"
	"os"

	"./blotter"
	"./blotter/api"
)

func main() {
	logger := log.New(os.Stdout, "[blotter]", log.LstdFlags)
	server := blotter.NewBlotter(":8080", map[string]blotter.AnyFunc{
		"/index": api.Index,
	}, logger)
	server.Start()
}
