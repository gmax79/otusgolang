package main

import (
	"fmt"
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
	println("POST", "/"+path, "", strings.Join(all, ", "))
	resp, err := http.PostForm(host+"/"+path, values)
	if err != nil {
		fmt.Println(err)
		return
	}
	out(resp, requiredCode)
}

func get(path string, requiredCode int) {
	println("GET", "/"+path)
	resp, err := http.Get(host + "/" + path)
	if err != nil {
		fmt.Println(err)
		return
	}
	out(resp, requiredCode)
}

func out(resp *http.Response, requiredCode int) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Return Code:", resp.Status)
		fmt.Println("Content-Length:", resp.ContentLength)
		if resp.ContentLength > 0 {
			fmt.Println("Content:", string(data))
		}
		if requiredCode != resp.StatusCode {
			fmt.Println("ERROR. Code must ", requiredCode)
		}
	}
}

func main() {
	fmt.Println("Testing calendar app")
	post("", map[string]string{}, http.StatusNotFound)
	post("a", map[string]string{}, http.StatusNotFound)
	post("b", map[string]string{}, http.StatusNotFound)

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

	r3 := map[string]string{
		"time":  "2020-01-07 12:00:00",
		"event": "Party",
	}
	post("create_event", r3, http.StatusOK)
	r3move := map[string]string{
		"time":    "2020-01-07 12:00:00",
		"event":   "Party",
		"newtime": "2020-01-10 18:00:00",
	}
	post("move_event", r3move, http.StatusOK)

	r5 := map[string]string{
		"time":  "2020-01-08 16:00:00",
		"event": "Party in club",
	}
	post("create_event", r5, http.StatusOK)

	r6 := map[string]string{
		"time":  "2020-01-15 12:00:00",
		"event": "Exam",
	}
	post("create_event", r6, http.StatusOK)

	get("events_for_day?day=2020-01-07", http.StatusOK)
	get("events_for_day?day=2020-01-10", http.StatusOK)
	get("events_for_week?week=2020-02", http.StatusOK)
	get("events_for_month?month=2020-01", http.StatusOK)

}
