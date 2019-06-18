package blotter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

// Any 任意类型
type Any interface{}

// AnyFunc 任意函数
type AnyFunc interface{}

// Blotter 服务端数据结构
type Blotter struct {
	address        string
	handle         *Handle
	globalVariable map[string]interface{}
	Logger         *log.Logger
	server         http.Server
}

// NewBlotter 构造一个Blotter对象
func NewBlotter(address string, router map[string]AnyFunc, logger *log.Logger) Blotter {
	blotter := Blotter{
		Logger: logger,
		server: http.Server{
			Addr:    address,
			Handler: nil,
		},
	}
	handle := Handle{
		router:  router,
		Blotter: &blotter,
	}
	blotter.server.Handler = &handle
	return blotter
}

// Start 启动Blotter服务
func (b *Blotter) Start() {
	fmt.Printf("Server start at %s\n", b.server.Addr)
	b.server.ListenAndServe()
}

// Stop 关闭Blotter服务
func (b *Blotter) Stop() {
	fmt.Println("Server stop")
	b.server.Shutdown(nil)
}

// Handle Blotter路由处理
type Handle struct {
	router  map[string]AnyFunc
	Blotter *Blotter
}

// ServeHTTP http路由处理函数
func (handle *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	solveFunc, ok := handle.router[url]

	r.ParseForm()
	args := map[string]string{}
	for key, value := range r.Form {
		args[key] = value[0]
	}
	argsBytes, _ := json.Marshal(args)

	handle.Blotter.Logger.Println(url, string(argsBytes))

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.Write(handle.Blotter.Solve(argsBytes, solveFunc))
	} else {
		// 404
		w.WriteHeader(404)
		w.Write([]byte("404 Not Found"))
	}
}

// Solve 处理函数装饰函数
func (b *Blotter) Solve(data []byte, f interface{}) []byte {
	// 被装饰函数类型获取
	funcType := reflect.TypeOf(f)
	inputType := funcType.In(0)
	//  outputType := funcType.Out(0)

	// 输入参数处理
	inputData := reflect.New(inputType).Interface()
	err := json.Unmarshal(data, inputData)
	if err != nil {
		b.Logger.Println(err)
	}
	args := []reflect.Value{reflect.ValueOf(inputData).Elem()}

	// 函数调用
	outputDataList := reflect.ValueOf(f).Call(args)

	// 输出内容处理
	outputData := outputDataList[0].Interface()
	res, _ := json.Marshal(outputData)

	b.Logger.Printf("%+v %+v %+v %+v \n", funcType, string(data), inputData, args[0])

	return res
}
