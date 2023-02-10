package fly

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type Link struct {
	Filter  func(c *Context) error
	execute func(c *Context, next *Link) error
	next    *Link
}

func (my *Link) addAfter(fn func(c *Context, next *Link) error) *Link {
	my.next = &Link{}
	my.execute = fn
	my.Filter = func(c *Context) error {
		return my.execute(c, my.next)
	}
	return my.next
}

func (my *Link) lastLink() *Link {
	last := my
	for last.next != nil {
		last = last.next
	}
	return last
}

var globalMiddlewareLink = &Link{}

func Middleware(middleware ...func(c *Context, next *Link) error) *Link {
	last := globalMiddlewareLink.lastLink()
	for _, m := range middleware {
		last = last.addAfter(m)
	}
	return last
}

type Route struct {
	method         string                 //http method
	middlewareLink *Link                  //link of middleware , will initial by first middleware
	httpFun        func(c *Context) error //last middleware
}

func NewHttpRoute(method string, fn func(c *Context) error) *Route {
	first := &Link{}
	first.addAfter(func(c *Context, next *Link) error {
		return next.Filter(c)
	})
	return &Route{
		method:         method,
		middlewareLink: first, // initial first middleware
		httpFun:        fn,    // last middleware
	}
}

func (my *Route) Use(middleware ...func(c *Context, next *Link) error) *Route {
	if len(middleware) == 0 {
		return my
	}
	last := my.middlewareLink.lastLink()
	for _, m := range middleware {
		last = last.addAfter(m)
	}
	return my
}

func ListenAndServe(addr string) error {
	//add router middleware
	routeRegisterMutex.Lock()
	for _, m := range routes {
		for _, route := range m {
			cur := route
			route.middlewareLink.lastLink().addAfter(func(c *Context, _ *Link) error {
				return cur.httpFun(c)
			})
		}
	}
	routeRegisterMutex.Unlock()

	//add last global middleware
	globalMiddlewareLink.lastLink().addAfter(func(c *Context, _ *Link) error {
		//execute router middleware
		return c.route.middlewareLink.execute(c, c.route.middlewareLink.next)
	})

	addr = strings.TrimSpace(addr)
	if addr == "" {
		return errors.New("HTTP listen addr is empty")
	}
	if addr[0] != ':' {
		addr = ":" + addr
	}

	return http.ListenAndServe(addr, nil)
}
