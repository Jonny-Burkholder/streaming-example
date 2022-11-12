package netkit

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
)

// NotFound is a helper that returns a 404 error status to the client.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// NotImplemented is a helper that returns a 501 error status to the client.
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// BadRequest is a helper that returns a 400 error status to the client.
func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// WriteRaw is a helper that takes a response code and some optional byte data.
// It will attempt to detect the content type of the data (if any is provided)
// and will set the Content-Type headers automatically before writing to the
// response writer.
func WriteRaw(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if data == nil {
		data = []byte(http.StatusText(code))
	}
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		NotImplemented(w, r)
	}
}

// ReadRaw is a helper that simply reads and returns the request body. It
// also returns the content type along with the full request body.
func ReadRaw(w http.ResponseWriter, r *http.Request) ([]byte, string) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, r)
	}
	return data, http.DetectContentType(data)
}

// WriteJSON is a helper that attempts to JSON encode the data provided. It
// sets the Content-Type headers and writes the status code provided before
// writing to the http.ResponseWriter.
func WriteJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		NotImplemented(w, r)
	}
}

// ReadJSON is a helper that assumes the request body contains valid JSON data.
// It attempts to read and decode the JSON data from the request body and will
// decode into the provided ptr.
func ReadJSON(w http.ResponseWriter, r *http.Request, ptr any) {
	err := json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		BadRequest(w, r)
	}
}

// WriteXML is a helper that attempts to XML encode the data provided. It
// sets the Content-Type headers and writes the status code provided before
// writing to the http.ResponseWriter.
func WriteXML(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))
	w.WriteHeader(code)
	err := xml.NewEncoder(w).Encode(data)
	if err != nil {
		NotImplemented(w, r)
	}
}

// ReadXML is a helper that assumes the request body contains valid XML data.
// It attempts to read and decode the XML data from the request body and will
// decode into the provided ptr.
func ReadXML(w http.ResponseWriter, r *http.Request, ptr any) {
	err := xml.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		BadRequest(w, r)
	}
}

// WriteErrorJSON is a helper that takes a response code and an error and
// encodes a JSON error message that will be sent back to the client. The
// provided error is included in the response.
func WriteErrorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(
		struct {
			Code   int    `json:"code"`
			Status string `json:"status"`
			Error  error  `json:"error"`
		}{
			Code:   code,
			Status: http.StatusText(code),
			Error:  err,
		},
	)
	if e != nil {
		NotImplemented(w, r)
	}
}
