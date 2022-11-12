package netkit

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Meth is a middleware function that ensures that the method you provide
// along with the path and handler is executed only if the request method
// is the same as the one provided.
func Meth(method string, path string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

// Get takes a handler and returns a handler that will only answer GET requests.
func Get(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

// Post takes a handler and returns a handler that will only answer POST requests.
func Post(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

// Put takes a handler and returns a handler that will only answer PUT requests.
func Put(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

// Delete takes a handler and returns a handler that will only answer DELETE requests.
func Delete(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		h(w, r)
	}
}

// HandleStatic is a simple static file handler that takes a prefix
// and filepath to map to. It returns a http.Handler
func HandleStatic(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}

// ServeAFileHandler takes a file pointer and returns a HandlerFunc.
// It attempts to read the file and send it to the client. The caller
// is responsible for opening and closing the file.
func ServeAFileHandler(file *os.File) http.HandlerFunc {
	// init buffer
	buf := make([]byte, 512)
	// read the fist bit
	_, err := file.Read(buf)
	if err != nil {
		panic(err)
	}
	// detect content type
	var contentType string
	contentType = http.DetectContentType(buf)
	// rewind reader back to beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	// now we can actually handle our web request
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, err := io.CopyBuffer(w, file, buf)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusExpectationFailed), http.StatusExpectationFailed)
			return
		}
	}
}

// HandleMetrics is a simple handler that takes a title, and a slice of URL
// paths, and renders a simple web page allowing you to see how your routes
// are being handled.
func HandleMetrics(title string, ss []string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data []string
		data = append(data, fmt.Sprintf("<h3>%s</h3>", title))
		sort.Strings(ss)
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		_, err := fmt.Fprintf(w, strings.Join(data, "<br>"))
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		return
	}
	return http.HandlerFunc(fn)
}

// HandleErrors expects a path of /error/{code} and will display an error page.
func HandleErrors(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) > 1 {
		code, err := strconv.Atoi(p[1])
		if err != nil {
			WriteErrorJSON(w, r, code, err)
			return
		}
		WriteErrorJSON(w, r, code, errors.New(http.StatusText(code)))
		return
	}
}

// HandleWithLogging is a middleware function that takes a *Logger, and a
// http.Handler and returns a http.Handler. It will log everything that
// passes through the http.Handler you provide.
func HandleWithLogging(logger *Logger, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("err: %v, trace: %s\n", err, debug.Stack())
			}
		}()
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: 200,
				size:   0,
			},
		}
		next.ServeHTTP(&lrw, r)
		if 400 <= lrw.data.status && lrw.data.status <= 599 {
			str, args := logStr(lrw.data.status, r)
			logger.Error(str, args...)
			return
		}
		str, args := logStr(lrw.data.status, r)
		logger.Info(str, args...)
		return
	}
	return http.HandlerFunc(fn)
}

type responseData struct {
	status int
	size   int
}

// loggingResponseWriter implements the http.ResponseWriter interface
// with Header(), Write(b []byte) (int, error) and WriteHeader(int)
type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}

func logStr(code int, r *http.Request) (string, []interface{}) {
	return "# %s - - [%s] \"%s %s %s\" %d %d\n", []interface{}{
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.URL.EscapedPath(),
		r.Proto,
		code,
		r.ContentLength,
	}
}
