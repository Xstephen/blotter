package main

import (
	"./blotter"
	"./blotter/api"
)

func main() {
	server := blotter.NewBlotter(":8080", map[string]blotter.AnyFunc{
		"/index": api.Index,
	})
	server.Start()
}
