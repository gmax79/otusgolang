package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const host = "http://localhost:8080"

func post(path string, params map[string]string) {
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
		println("Data:", string(data))
	}
}

func main() {
	println("Testing calendar app")

	r1 := map[string]string{
		"time":  "2020.10.22 18:00:00",
		"event": "M birthday",
	}
	post("create_event", r1)
}
