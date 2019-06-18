package blotter

import (
	"log"
	"net/http"
)

// Context blotter上下文对象.
type Context struct {
	blotter  *Blotter
	request  *http.Request
	response *http.Response
}

// GetBlotter 获得当前上下文的Blotter对象
func (ctx *Context) GetBlotter() *Blotter {
	return ctx.blotter
}

// GetRequest 获得当前上下文的的请求
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

// GetResponse 获得当前上下文的响应
func (ctx *Context) GetResponse() *http.Response {
	return ctx.response
}

// GetLogger 获得日志输出器
func (ctx *Context) GetLogger() *log.Logger {
	return ctx.blotter.Logger
}
