package fly

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	writer http.ResponseWriter
	reader *http.Request
	st     time.Time
	result *Result
}

func NewContext(writer http.ResponseWriter, reader *http.Request) *Context {
	return &Context{writer: writer, reader: reader, st: time.Now()}
}
func (my *Context) Query(k string) string {
	//todo  implementing ..
	//todo  implementing ..
	return ""
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

func (my *Context) Request() *http.Request {
	return my.reader
}
func (my *Context) Writer() http.ResponseWriter {
	return my.writer
}
