package support

import (
	"encoding/json"
	"net/http"
)

// HTTPResponse - send http answer to client
func HTTPResponse(w http.ResponseWriter, v interface{}) {
	switch v.(type) {
	case nil:
		w.WriteHeader(http.StatusNotFound)
	case int:
		w.WriteHeader(v.(int))
	case error:
		err := v.(error)
		text := "{ \"error\" : \"" + err.Error() + "\"  }\n"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(text))
	default:
		answer, err := json.Marshal(v)
		if err != nil {
			HTTPResponse(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(answer)
		}
	}
}

// HTTPPostRequest - post reader
type HTTPPostRequest struct {
	request *http.Request
}

// ReadPostRequest - check what is post, and parse parameters
func ReadPostRequest(r *http.Request, w http.ResponseWriter) *HTTPPostRequest {
	if r.Method != http.MethodPost {
		HTTPResponse(w, http.StatusInternalServerError)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		HTTPResponse(w, err)
		return nil
	}
	return &HTTPPostRequest{request: r}
}

// Get - simple read post parameter
func (r *HTTPPostRequest) Get(id string) string {
	request := r.request
	return request.Form.Get(id)
}

// HTTPGetRequest - get method reader
type HTTPGetRequest struct {
	request *http.Request
}

// ReadGetRequest - check what is get
func ReadGetRequest(r *http.Request, w http.ResponseWriter) *HTTPGetRequest {
	if r.Method != http.MethodGet {
		HTTPResponse(w, http.StatusInternalServerError)
		return nil
	}
	return &HTTPGetRequest{request: r}
}

// Get - reader get query parameter
func (r *HTTPGetRequest) Get(id string) string {
	p := r.request.URL.Query()
	values := p[id]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
