package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const host = "http://localhost:8080"

func post(path string, params map[string]string, requiredCode int) {
	all := []string{}
	values := url.Values{}
	for k, v := range params {
		all = append(all, k+"='"+v+"'")
		values[k] = []string{v}
	}
	println("POST", path, "", strings.Join(all, ", "))
	resp, err := http.PostForm(host+"/"+path, values)
	if err != nil {
		println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
	} else {
		println("Return Code:", resp.Status)
		println("Content-Length:", resp.ContentLength)
		if resp.ContentLength > 0 {
			println("Content:", string(data))
		}
		if requiredCode != resp.StatusCode {
			println("ERROR. Code must ", requiredCode)
		}
	}
}

func main() {
	println("Testing calendar app")

	r1 := map[string]string{
		"time":  "2020-10-22 18:00:00",
		"event": "Maks birthday",
	}
	post("create_event", r1, http.StatusOK)

	r2 := map[string]string{
		"time":  "2020-10-22 18:00:00",
		"event": "Maks birthday",
	}
	post("delete_event", r2, http.StatusOK)
}
