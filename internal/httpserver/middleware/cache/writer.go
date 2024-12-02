package cache

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type cachedResponseWriter struct {
	gin.ResponseWriter
	data         []byte
	statusCode   int
	cachedHeader http.Header
	lock         sync.Mutex
}

func newCachedWriter(originalWriter gin.ResponseWriter) *cachedResponseWriter {
	return &cachedResponseWriter{
		ResponseWriter: originalWriter,
		cachedHeader:   http.Header{},
	}
}

func (w *cachedResponseWriter) DoCache() {
	w.lock.Lock()
	defer w.lock.Unlock()
	for k, v := range w.ResponseWriter.Header() {
		for _, vv := range v {
			w.cachedHeader.Add(k, vv)
		}
	}
}

func (w *cachedResponseWriter) Header() http.Header {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.ResponseWriter.Header()
}

func (w *cachedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedResponseWriter) Status() int {
	return w.statusCode
}

func (w *cachedResponseWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *cachedResponseWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		w.data = append(w.data, data...)
	}
	return ret, err
}

func (w *cachedResponseWriter) WriteString(data string) (n int, err error) {
	ret, err := w.ResponseWriter.WriteString(data)
	//cache responses with a status code < 300
	if err == nil {
		w.data = append(w.data, []byte(data)...)
	}
	return ret, err
}

func (w *cachedResponseWriter) FlushTo(writer gin.ResponseWriter) {
	writer.WriteHeader(w.Status())
	for k, values := range w.cachedHeader {
		for _, v := range values {
			writer.Header().Add(k, v)
		}
	}
	_, _ = writer.Write(w.data)
}
