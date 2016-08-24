package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

const (
	//apiRoot = "http://kwk.loc/api/v1/"
	apiRoot = "http://localhost:8080/api/v1/"
)

type ApiClient struct {
}

type KwkLink struct {
	Id      int64 `json:"id"`
	Key     string `json:"key"`
	Root    string `json:"root"`
	Uri     string `json:"url"`
	AfToken string `json:"afToken"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (k *KwkLink) Err() string {
	return k.Error
}

func (a *ApiClient) Decode(key string) *KwkLink {
	k := &KwkLink{}
	Request("GET", fmt.Sprintf("hash/%s", key), "", k)
	return k
}

func (a *ApiClient) Create(uri string, path string) *KwkLink {
	body := fmt.Sprintf(`{"url":"%s", "key":"%s"}`, uri, path)
	k := &KwkLink{}
	Request("POST", "hash", body, k)
	return k
}

func Request(method string, path string, body string, response interface{}) {
	url := fmt.Sprintf("%s%s", apiRoot, path)
	var req *http.Request
	if body != "" {
		b := []byte(body)
		buffer := bytes.NewBuffer(b)
		req, _ = http.NewRequest(method, url, buffer)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("x-kwk-key", "d50ce4ec-97dc-46f2-a247-5d2a834caedf")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	r, e := client.Do(req)
	if e != nil {
		fmt.Print("kwk server is unavailable, please try again later or tweet us @kwklinks.")
		return
	}
	defer r.Body.Close()
	responseBytes, _ := ioutil.ReadAll(r.Body)
	if e := json.Unmarshal(responseBytes, response); e != nil {
		handleResponse(response, r)
		return
	}
	handleResponse(response, r)
}

func handleResponse(i interface{}, r *http.Response) {
	switch {
	case r.StatusCode == http.StatusBadRequest :
		fmt.Println(i.(ErrorResponse).Err())
	case r.StatusCode == http.StatusForbidden :
		fmt.Println("Sign in please: 'kwk signin <username> <password>'")
	case r.StatusCode != http.StatusOK :
		fmt.Println(r)
		fmt.Println(r.StatusCode)
	}
}

type ErrorResponse interface{
	Err() string
}

