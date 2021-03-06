package stests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type result struct {
	Count int `json:"result"`
}

func Post(host, path string, params map[string]string, requiredCode int) {
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

func Get(host, path string, requiredCode int, resultCount int) {
	println("GET", "/"+path)
	resp, err := http.Get(host + "/" + path)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := out(resp, requiredCode)
	s := string(data)
	s = strings.ReplaceAll(s, "\\", "")
	var r result
	err = json.Unmarshal([]byte(s), &r)
	if err != nil {
		fmt.Println(err)
	}
	if r.Count != resultCount {
		fmt.Println("ERROR, Count must ", resultCount)
	}
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
