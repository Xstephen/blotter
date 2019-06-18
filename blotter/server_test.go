package blotter

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

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

type Client struct {
	*http.Client
}

func (c *Client) GetRetry(times int, second time.Duration, path string) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < times; i++ {
		resp, err = c.Get(path)
		if err == nil {
			break
		}
		time.Sleep(second * time.Second)
	}
	return resp, err
}

func (c *Client) PostRetry(times int, second time.Duration, path string, form url.Values) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < times; i++ {
		resp, err = c.PostForm(path, form)
		if err == nil {
			break
		}
		time.Sleep(second * time.Second)
	}
	return resp, err
}

func TestBlotter(t *testing.T) {
	checkError := func(err error) {
		if err != nil && err.Error() != "EOF" {
			t.Errorf("Error: %+v\n", err)
			t.FailNow()
		}
	}
	getResponseString := func(r *http.Response) string {
		length, err := strconv.Atoi(r.Header.Get("content-length"))
		checkError(err)
		data := make([]byte, length)
		length, err = r.Body.Read(data)
		checkError(err)
		return string(data)
	}

	type inputType struct {
		Name string `json:"name"`
	}
	type outputType struct {
		Data string `json:"data"`
	}
	router := map[string]AnyFunc{
		"/test": func(in inputType) outputType {
			return outputType{
				Data: "Hi, " + in.Name + "!",
			}
		},
	}
	logger := log.New(os.Stdout, "[blotter]", log.LstdFlags)
	blotter := NewBlotter(":8080", router, logger)

	go blotter.Start()

	client := Client{
		&http.Client{},
	}

	resp, err := client.GetRetry(5, 1, "http://127.0.0.1:8080/test?name=Test")
	checkError(err)
	getResult := getResponseString(resp)
	if getResult == `{"data":"Hi, Test!"}` {
		t.Log("Get method ok")
	} else {
		t.Errorf("Get method error. Result: %s\n", getResult)
	}

	resp, err = client.PostRetry(5, 1, "http://127.0.0.1:8080/test", url.Values{"name": []string{"Test"}})
	checkError(err)
	postResult := getResponseString(resp)
	if postResult == `{"data":"Hi, Test!"}` {
		t.Log("Post method ok")
	} else {
		t.Errorf("Post method error. Result: %s\n", postResult)
	}

	blotter.Stop()
}
