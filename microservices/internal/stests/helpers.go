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

// PostRequest - helper to create post request with some parameters
func PostRequest(host, path string, params map[string]string) (*http.Response, error) {
	all := []string{}
	values := url.Values{}
	if params != nil {
		for k, v := range params {
			all = append(all, k+"='"+v+"'")
			values[k] = []string{v}
		}
	}
	fmt.Println("POST", "/"+path, "", strings.Join(all, ", "))
	return http.PostForm(host+"/"+path, values)

}

// Post - test function to make post and check returned code
func Post(host, path string, params map[string]string, requiredCode int) {
	resp, err := PostRequest(host, path, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = out(resp, requiredCode)
	if err != nil {
		fmt.Println(err)
	}
}

// Get - test function to make get and check returned code
func Get(host, path string, requiredCode int, resultCount int) {
	fmt.Println("GET", "/"+path)
	resp, err := http.Get(host + "/" + path)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := out(resp, requiredCode)
	if err != nil {
		fmt.Println(err)
	}
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

func out(resp *http.Response, requiredCode int) ([]byte, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	fmt.Println("Return Code:", resp.Status)
	fmt.Println("Content-Length:", resp.ContentLength)
	if resp.ContentLength > 0 {
		fmt.Println("Content:", string(data))
		data = data[1 : len(data)-1]
	}
	if requiredCode != resp.StatusCode {
		return data, fmt.Errorf("ERROR, Code must %d", requiredCode)
	}
	return data, nil
}
