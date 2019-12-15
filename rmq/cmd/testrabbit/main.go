package main

import (
	"fmt"
	"github.com/gmax79/otusgolang/rmq/internal/calendar"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const host = "http://localhost:8080"

func post(path string, params map[string]string, requiredCode int) {
	all := []string{}
	values := url.Values{}
	for k, v := range params {
		all = append(all, k+"='"+v+"'")
		values[k] = []string{v}
	}
	println("POST", "/"+path, "", strings.Join(all, ", "))
	resp, err := http.PostForm(host+"/"+path, values)
	if err != nil {
		fmt.Println(err)
		return
	}
	out(resp, requiredCode)
}

func out(resp *http.Response, requiredCode int) []byte {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Return Code:", resp.Status)
		fmt.Println("Content-Length:", resp.ContentLength)
		if resp.ContentLength > 0 {
			fmt.Println("Content:", string(data))
			data = data[1 : len(data)-1]
		}
		if requiredCode != resp.StatusCode {
			fmt.Println("ERROR, Code must ", requiredCode)
		}
	}
	return data
}

func main() {
	fmt.Println("Testing rabbit mq pipeline. Create nearby events")
	r1 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 5),
		"event": "RabbitMQ #1",
	}
	post("create_event", r1, http.StatusOK)

	r2 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.1",
	}
	post("delete_event", r2, http.StatusOK)

	r3 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.2",
	}
	post("create_event", r3, http.StatusOK)
}
