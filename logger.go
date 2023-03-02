package fly

import (
	"io"
	"log"
	"strings"
)

var flyDefaultErrorLogger = log.New(&httpErrLogWriter{writer: log.Default().Writer()}, log.Default().Prefix(), log.Default().Flags())

type httpErrLogWriter struct {
	writer io.Writer
}

func (my *httpErrLogWriter) Write(p []byte) (n int, err error) {

	if strings.Contains(string(p), "golang.org/issue/25192") {
		//server.go:3215: http: URL query contains semicolon, which is no longer a supported separator; parts of the query may be stripped when parsed; see golang.org/issue/25192
		return 0, nil
	}

	return my.writer.Write(p)
}
