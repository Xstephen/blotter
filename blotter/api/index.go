package api

import "fmt"

type indexInputType struct {
	Name string `json:"name"`
}
type indexOutputType struct {
	Data string `json:"data"`
}

func Index(in indexInputType) indexOutputType {
	return indexOutputType{
		Data: fmt.Sprintf("Hi, %s", in.Name),
	}
}
