package cache

import (
	"bufio"
	"net"
	"net/http"
	"sync"
	"testing"
)

type MockWriter struct {
	header http.Header
}

func (w MockWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (w MockWriter) WriteHeader(statusCode int) {
}

func (w MockWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (w MockWriter) Flush() {
}

func (w MockWriter) CloseNotify() <-chan bool {
	return nil
}

func (w MockWriter) Status() int {
	return 0
}

func (w MockWriter) Size() int {
	return 0
}

func (w MockWriter) WriteString(s string) (int, error) {
	return 0, nil
}

func (w MockWriter) Written() bool {
	return false
}

func (w MockWriter) WriteHeaderNow() {
}

func (w MockWriter) Pusher() http.Pusher {
	return nil
}

func (w MockWriter) Header() http.Header {
	return w.header
}

func TestCachedResponseWriter_FlushTo(t *testing.T) {
	cachedWriter := newCachedWriter(MockWriter{header: map[string][]string{}})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			writer := MockWriter{header: map[string][]string{}}
			cachedWriter.FlushTo(writer)
			wg.Done()
		}()
	}
	wg.Add(1)
	go func() {
		cachedWriter.Header().Set("", "")
		wg.Done()
	}()
	wg.Wait()
}
