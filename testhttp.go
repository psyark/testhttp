package testhttp

import (
	"bytes"
	"fmt"
	"net/http"
)

type Option func(w *responseWriter) error

func TestHTTP(r *http.Request, options ...Option) error {
	w := &responseWriter{body: bytes.NewBuffer(nil)}
	http.DefaultServeMux.ServeHTTP(w, r)
	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}
	return nil
}

func PrintBodyAsString(w *responseWriter) error {
	fmt.Println(w.body.String())
	return nil
}

func MustStatus(status int) Option {
	return func(w *responseWriter) error {
		if w.status != status {
			return fmt.Errorf("status: got=%d, want=%d", w.status, status)
		}
		return nil
	}
}

func MustOK(w *responseWriter) error {
	if w.status != http.StatusOK {
		return fmt.Errorf("status: got=%d, want=%d", w.status, http.StatusOK)
	}
	return nil
}

type responseWriter struct {
	header     http.Header
	headerSent bool
	status     int
	body       *bytes.Buffer
}

func (r *responseWriter) Header() http.Header {
	if r.header == nil {
		r.header = http.Header{}
	}
	return r.header
}

func (r *responseWriter) Write(p0 []byte) (int, error) {
	if !r.headerSent {
		r.WriteHeader(200)
	}
	return r.body.Write(p0)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if !r.headerSent {
		_, _ = fmt.Fprintf(r.body, "HTTP/1.1 %d\n", statusCode)
		for k, v := range r.header {
			for _, vv := range v {
				_, _ = fmt.Fprintf(r.body, "%v: %v\n", k, vv)
			}
		}
		_, _ = fmt.Fprintln(r.body, "")
		r.status = statusCode
		r.headerSent = true
	}
}
