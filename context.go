package fly

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	HttpHead    = "HEAD"
	HttpPatch   = "PATCH"
	HttpDelete  = "DELETE"
	HttpOptions = "OPTIONS"
	HttpGet     = "GET"
	HttpPost    = "POST"
	HttpANY     = "ANY"
)

type Context struct {
	writer   http.ResponseWriter
	reader   *http.Request
	attr     map[string]interface{}
	attrLock sync.RWMutex
	route    *Route
	st       time.Time
	result   *Result
}

func NewContext(route *Route, writer http.ResponseWriter, reader *http.Request) *Context {
	return &Context{route: route, writer: writer, reader: reader, attr: make(map[string]interface{}), st: time.Now()}
}

func (my *Context) GetAttr(k string) any {
	my.attrLock.RLock()
	defer my.attrLock.RUnlock()
	return my.attr[k]
}

func (my *Context) SetAttr(k string, v any) (originValue any) {
	my.attrLock.Lock()
	my.attr[k] = v
	my.attrLock.Unlock()
	return v
}

func (my *Context) Query(k string) string {
	request := my.Request()
	if request == nil || request.URL == nil || request.URL.Query() == nil {
		return ""
	}
	return request.URL.Query().Get(k)
}

func (my *Context) Post(k string) string {
	//todo  implementing ..
	//todo  implementing ..
	return ""
}
func (my *Context) Decode(v interface{}) error {
	return json.NewDecoder(my.Request().Body).Decode(v)
}

func (my *Context) Result(v interface{}) error {
	my.result = NewResult(v, http.StatusOK)
	return nil
}

func (my *Context) Json(v interface{}) error {
	b, err := json.Marshal(map[string]interface{}{
		"t":    time.Since(my.st).String(),
		"data": v,
	})
	if err == nil {
		my.result = NewResult(b, http.StatusOK)
	}
	return err
}

func (my *Context) Error(message string) error {
	my.result = NewResult("{\"msg\":"+strconv.Quote(message)+"}", http.StatusUnprocessableEntity)
	return nil
}

func (my *Context) Success(message string) error {
	my.result = NewResult("{\"msg\":"+strconv.Quote(message)+"}", http.StatusOK)
	return nil
}
func (my *Context) Fail(httpCode int) error {
	t := http.StatusText(httpCode)
	if t == "" {
		//UNKNOWN
		my.result = NewResult("{\"msg\":\"UNKNOWN\"}", httpCode)
	} else {
		my.result = NewResult("{\"msg\":"+strconv.Quote(t)+"}", httpCode)
	}
	return nil
}

func (my *Context) RequestIPs() (ips []string) {
	ip := my.RequestIP()
	if ip != "" {
		return []string{ip}
	}
	return
}

func (my *Context) RequestIP() (ip string) {
	req := my.Request()
	cf := strings.TrimSpace(req.Header.Get("CF-Connecting-IP"))
	if cf != "" {
		for _, v := range strings.Split(cf, ",") {
			return v
		}
	}
	forward := strings.TrimSpace(req.Header.Get("X-Forwarded-For"))
	if forward != "" {
		for _, v := range strings.Split(forward, ",") {
			return v
		}
	}
	remote, _, _ := net.SplitHostPort(req.RemoteAddr)
	return remote
}

func (my *Context) Request() *http.Request {
	return my.reader
}
func (my *Context) Writer() http.ResponseWriter {
	return my.writer
}
