package blotter

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"testing"

	"../util/nilwriter"
)

// TestSolve
func TestSolve(t *testing.T) {
	type args struct {
		data []byte
		f    AnyFunc
	}
	var tests = []struct {
		in       args
		expected []byte
	}{
		{
			args{
				data: []byte(`{"a":1,"b":2}`),
				f: func(i struct {
					A int `json:"a"`
					B int `json:"b"`
				}) struct {
					Sum int `json:"sum"`
				} {
					return struct {
						Sum int `json:"sum"`
					}{
						Sum: i.A + i.B,
					}
				},
			},
			[]byte(`{"sum":3}`),
		},
		{
			args{
				data: []byte(`{"a":"a","b":"b"}`),
				f: func(i struct {
					A string `json:"a"`
					B string `json:"b"`
				}) struct {
					C string `json:"c"`
				} {
					return struct {
						C string `json:"c"`
					}{
						C: i.A + i.B,
					}
				},
			},
			[]byte(`{"c":"ab"}`),
		},
	}

	blotter := NewBlotter(":8080", map[string]AnyFunc{}, log.New(nilwriter.NilWriter{}, "", log.LstdFlags))
	for _, test := range tests {
		actual := blotter.Solve(test.in.data, test.in.f)
		if bytes.Compare(actual, test.expected) != 0 {
			t.Errorf("[×] in: %+v out: %+v expected: %+v\n", test.in, actual, test.expected)
		} else {
			t.Logf("[√] in: %+v out: %+v expected: %+v\n", test.in, actual, test.expected)
		}
	}
}

func TestBlotter(t *testing.T) {
	blotter := NewBlotter(":8080", map[string]AnyFunc{}, log.New(os.Stdout, "[blotter]", log.LstdFlags))
	go blotter.Start()
	http.Get()
}
