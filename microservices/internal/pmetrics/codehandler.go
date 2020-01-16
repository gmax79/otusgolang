package pmetrics

import (
	"fmt"
	"net/http"
	"strconv"
)

type urlcode struct {
	url  string
	code int
}

type urlcodevalue struct {
	count  int
	setter GaugeFunc
}

// handler to count http return codes
type returnCodesMetricsHandler struct {
	values map[urlcode]urlcodevalue
	agent  *Agent
}

// function to return wrapper handler from main handler
func (mh *returnCodesMetricsHandler) Attach(labels map[string]string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var k urlcode
		k.url = r.URL.Path
		var rw responseWriterProxy
		rw.statusCode = http.StatusOK
		rw.ResponseWriter = w
		h.ServeHTTP(&rw, r)
		k.code = rw.statusCode
		if v, ok := mh.values[k]; ok {
			v.count = v.count + 1
			mh.values[k] = v
			v.setter(float64(v.count))

			fmt.Println(v)
		} else {
			labels["handler"] = k.url
			labels["code"] = strconv.Itoa(k.code)
			f, err := mh.agent.RegisterGaugeMetric("http_return_codes_count", "counts return codes by url", labels)
			if err != nil {
				return
			}
			var v urlcodevalue
			v.count = 1
			v.setter = f
			v.setter(float64(v.count))
			mh.values[k] = v

			fmt.Println(v)
		}

	})
}

func createReturnCodesMetricsHandler(a *Agent) MetricsHandler {
	var h returnCodesMetricsHandler
	h.values = make(map[urlcode]urlcodevalue)
	h.agent = a
	return &h
}
