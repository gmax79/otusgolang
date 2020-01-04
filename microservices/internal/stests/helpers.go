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

// GetRequest - helper to create get request
func GetRequest(host, path string) (*http.Response, error) {
	fmt.Println("GET", "/"+path)
	resp, err := http.Get(host + "/" + path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetContent - test function to make get and return answer content
func GetContent(host, path string, requiredCode int) ([]byte, error) {
	resp, err := GetRequest(host, path)
	if err != nil {
		return []byte{}, err
	}
	data, err := out(resp, requiredCode)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Get - test function to make get and check returned code and count
func Get(host, path string, requiredCode int, resultCount int) {
	data, err := GetContent(host, path, requiredCode)
	if err != nil {
		fmt.Println(err)
		return
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
	fmt.Println("Bytes count:", len(data))
	if resp.ContentLength != int64(len(data)) {
		return data, fmt.Errorf("Content len not equal readed bytes chunk")
	}
	if len(data) > 2 {
		fmt.Println("Content:", string(data))
		data = data[1 : len(data)-1]
	}
	if requiredCode != resp.StatusCode {
		return data, fmt.Errorf("ERROR, Code must %d", requiredCode)
	}
	return data, nil
}
