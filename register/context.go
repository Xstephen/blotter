package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OhYee/blotter/output"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strconv"
	// "strings"
)

var decoder = schema.NewDecoder()

type httpHeader struct {
	key   string
	value string
}

// HandleContext context of a api call
type HandleContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	buf      *bytes.Buffer
	header   []httpHeader
}

// NewHandleContext initial a handle context object
func NewHandleContext(req *http.Request, rep http.ResponseWriter) *HandleContext {
	return &HandleContext{
		Request:  req,
		Response: rep,
		buf:      bytes.NewBuffer([]byte{}),
		header:   make([]httpHeader, 0),
	}
}

// RequestArgs get request args
func (context *HandleContext) RequestArgs(args interface{}) {
	context.Request.ParseForm()
	query := context.Request.Form

	decoder.Decode(args, query)
	output.Debug("args: %+v", args)

	// t := reflect.TypeOf(args).Elem()
	// v := reflect.ValueOf(args).Elem()
	// num := v.NumField()
	// output.Debug("%d", num)
	// for i := 0; i < num; i++ {
	// 	fieldType := t.Field(i)
	// 	fieldValue := v.Field(i)
	// 	value := query.Get(fieldType.Tag.Get("json"))
	// 	// output.Debug("%+v %+v %+v %+v", fieldType,)
	// 	setValue(fieldValue.Addr().Interface(), value)
	// }
}

func strToInt64(str string) (num int64) {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		num = 0
	}
	return
}

func strToUint64(str string) (num uint64) {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		num = 0
	}
	return
}

// ReturnJSON return json data
func (context *HandleContext) ReturnJSON(data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		return
	}

	context.AddHeader("Content-Type", "application/json")
	context.Write(b...)
	return
}

func (context *HandleContext) Write(b ...byte) {
	context.buf.Write(b)
}

// AddHeader add a header in response
func (context *HandleContext) AddHeader(key string, value string) {
	context.header = append(context.header, httpHeader{key, value})
}

func (context *HandleContext) writeHeaderWithCode(code int) {
	for _, header := range context.header {
		context.Response.Header().Add(header.key, header.value)
	}
	context.Response.WriteHeader(code)
}

// Success return 200 success
func (context *HandleContext) Success() {
	context.writeHeaderWithCode(200)
	context.Response.Write(context.buf.Bytes())
}

// PageNotFound return 404 page not found error
func (context *HandleContext) PageNotFound() {
	output.ErrOutput.Printf("404 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(404)
	context.Response.Write([]byte(fmt.Sprintf("Page not found %s", context.Request.RequestURI)))
}

// NotImplemented return 501 Not Implemented
func (context *HandleContext) NotImplemented() {
	output.ErrOutput.Printf("501 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(501)
	context.Response.Write([]byte(fmt.Sprintf("Can not solve request %s", context.Request.RequestURI)))
}

// ServerError return 500 server error
func (context *HandleContext) ServerError(err error) {
	output.ErrOutput.Printf("500 Server Error: %s\n", err.Error())
	context.writeHeaderWithCode(500)
	context.Response.Write([]byte(fmt.Sprintf("Server error %s", err.Error())))
}

func setValue(recv interface{}, value string) {
	v := reflect.ValueOf(recv).Elem()
	switch v.Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(value))
	case reflect.Int:
		v.Set(reflect.ValueOf(int(strToInt64(value))))
	case reflect.Int8:
		v.Set(reflect.ValueOf(int8(strToInt64(value))))
	case reflect.Int16:
		v.Set(reflect.ValueOf(int16(strToInt64(value))))
	case reflect.Int32:
		v.Set(reflect.ValueOf(int32(strToInt64(value))))
	case reflect.Int64:
		v.Set(reflect.ValueOf(int64(strToInt64(value))))
	case reflect.Uint8:
		v.Set(reflect.ValueOf(uint8(strToUint64(value))))
	case reflect.Uint16:
		v.Set(reflect.ValueOf(uint16(strToUint64(value))))
	case reflect.Uint32:
		v.Set(reflect.ValueOf(uint32(strToUint64(value))))
	case reflect.Uint64:
		v.Set(reflect.ValueOf(uint64(strToUint64(value))))
	case reflect.Slice:
		// strings.Split(value., ",")
		v.Set(reflect.ValueOf(uint64(strToUint64(value))))
	default:
		output.Debug("%+v %+v", v.Kind(), value)
	}
}
