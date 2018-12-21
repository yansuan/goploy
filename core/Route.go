package core

import (
	"net/http"
)

// 路由定义
type route struct {
	pattern     string                                         // 正则表达式
	callback    func(w http.ResponseWriter, r *http.Request)   //Controller函数
	middlewares []func(w http.ResponseWriter, r *http.Request) //中间件
}

// Router is route slice and golbal middlewares
type Router struct {
	Routes      []route
	middlewares []func(w http.ResponseWriter, r *http.Request) error //中间件
}

// Start a router
func (rt *Router) Start() {
	http.HandleFunc("/", rt.router)
}

// Add router
// pattern path
// callback  where path should be handle
func (rt *Router) Add(pattern string, callback func(w http.ResponseWriter, r *http.Request), middleware ...func(w http.ResponseWriter, r *http.Request)) {
	r := route{pattern: pattern, callback: callback}
	for _, m := range middleware {
		r.middlewares = append(r.middlewares, m)
	}
	rt.Routes = append(rt.Routes, r)
}

// Middleware golbal Middleware handle function
// Example handle praseToken
func (rt *Router) Middleware(middleware func(w http.ResponseWriter, r *http.Request) error) {
	rt.middlewares = append(rt.middlewares, middleware)
}

func (rt *Router) router(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range rt.middlewares {
		err := middleware(w, r)
		if err != nil {
			return
		}
	}
	for _, route := range rt.Routes {
		if route.pattern == r.URL.Path {
			for _, middleware := range route.middlewares {
				middleware(w, r)
			}
			route.callback(w, r)
		}
	}
}
