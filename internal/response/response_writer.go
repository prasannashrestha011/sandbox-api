package response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	err         error // <-- new field
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.status = statusCode
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) WroteHeader() bool {
	return rw.wroteHeader
}

func (rw *ResponseWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

// --- New methods ---
func (rw *ResponseWriter) SetError(err error) {
	rw.err = err
}

func (rw *ResponseWriter) Error() error {
	return rw.err
}
