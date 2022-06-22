package fly

import (
	"log"
	"net/http"
	"sync"
)

var routers = make(map[string]map[string]*Route)
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

func registerHttpMethod(method string, pattern string, httpFun func(c *Context) error) *Route {

	routeRegisterMutex.Lock()
	defer routeRegisterMutex.Unlock()

	route := NewHttpRoute(method, httpFun)
	if m, ok := routers[pattern]; ok {
		if _, ok := m[method]; ok {
			log.Panic(method + ":" + pattern + " already exists")
			return nil
		}
		m[method] = route
	} else {
		routers[pattern] = map[string]*Route{method: route}
		http.HandleFunc(pattern, routeSwitch(pattern))
	}

	return route
}

func routeSwitch(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := routers[pattern]
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		router, ok := m[r.Method]
		if !ok {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		c := NewContext(w, r)

		//global middleware
		if globalMiddlewareLink.next != nil {
			if err := globalMiddlewareLink.execute(c, globalMiddlewareLink.next); err != nil {
				httpFunc(c, err)
				return
			}
		}

		// route
		httpFunc(c, router.middlewareLink.execute(c, router.middlewareLink.next))
	}
}
