package fly

import (
	"log"
	"net/http"
	"sync"
)

var routes = make(map[string]map[string]*Route)
var routeRegisterMutex sync.Mutex

func ANY(path string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpANY, path, httpFun)
}

func GET(path string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpGet, path, httpFun)
}

func POST(pattern string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpPost, pattern, httpFun)
}

func OPTIONS(pattern string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpOptions, pattern, httpFun)
}

func HEAD(pattern string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpHead, pattern, httpFun)
}

func DELETE(pattern string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpDelete, pattern, httpFun)
}

func PATCH(pattern string, httpFun func(c *Context) error) *Route {
	return registerHttpMethod(HttpPatch, pattern, httpFun)
}

func registerHttpMethod(method string, pattern string, httpFun func(c *Context) error) *Route {

	routeRegisterMutex.Lock()
	defer routeRegisterMutex.Unlock()

	route := NewHttpRoute(method, httpFun)
	if m, ok := routes[pattern]; ok {
		if _, ok := m[method]; ok {
			log.Panic(method + ":" + pattern + " already exists")
			return nil
		}
		m[method] = route
	} else {
		routes[pattern] = map[string]*Route{method: route}
		http.HandleFunc(pattern, routeSwitch(pattern))
	}

	return route
}

func routeSwitch(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := routes[pattern]
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		route, ok := m[r.Method]
		if !ok {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		c := NewContext(route, w, r)
		//execute middleware
		httpFunc(c, globalMiddlewareLink.execute(c, globalMiddlewareLink.next))
	}
}
