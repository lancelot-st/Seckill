package common

import (
	"net/http"
	"strings"
)

//声明一个新的数据类型（函数类型）
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	//用来拦截需要拦截的URL
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

type WebHandle func(rw http.ResponseWriter, req *http.Request)

//执行拦截器返回拦截信息
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) {
				if path == r.RequestURI {
					//执行拦截业务逻辑

					err := handle(rw, r)
					if err != nil {
						rw.Write([]byte(err.Error()))
						return
					}
					break
				}
			}
		}
		//执行正常函数
		webHandle(rw, r)
	}
}
