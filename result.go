package fly

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func httpFunc(c *Context, err error) {
	if c.writer.Header().Get("Content-Type") == "" {
		//set default Content-Type
		c.writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}
	//c.writer.Header().Set("Connection", "keep-alive")
	//gzip
	if !strings.Contains(strings.ToLower(c.reader.Header.Get("Accept-Encoding")), "gzip") {
		//err := fn(c)
		writeResult(c.writer, c.result, err)
		return
	}
	c.writer.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(c.writer)
	defer gz.Close()
	gzr := &GzipResponseWriter{Writer: gz, ResponseWriter: c.writer}
	//err := fn(c)
	writeResult(gzr, c.result, err)
}

type Result struct {
	res      interface{}
	httpCode int
}

func NewResult(res interface{}, httpCode int) *Result {
	return &Result{res, httpCode}
}

func writeResult(w http.ResponseWriter, result *Result, err error) {
	if err != nil {
		log.Println(err)
		if result == nil {
			result = NewResult("{\"msg\":\""+http.StatusText(http.StatusInternalServerError)+"\"}", http.StatusInternalServerError)
		}
	}
	if result == nil {
		return
	}
	w.WriteHeader(result.httpCode)
	if result.res == nil {
		return
	}
	switch v := result.res.(type) {
	case string:
		io.WriteString(w, v)
	case []byte:
		w.Write(v)
	default:
		j, err := json.Marshal(v)
		if err != nil {
			log.Println("failed to marsha1 of result", err)
		}
		w.Write(j)
	}
}
