package api

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
)

const (apiRoot = "http://kwk.loc/api/v1/")

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

func (a *ApiClient) Decode(key string) *KwkLink {
	k := &KwkLink{}
	r, _, errs := gorequest.New().Get(fmt.Sprintf("%shash/%s", apiRoot, key)).EndStruct(k)
	if errs != nil {
		fmt.Println(errs)
	}
	if r.StatusCode == http.StatusOK { return k }
	if r.StatusCode == http.StatusNotFound { return nil }
	if r.StatusCode != http.StatusBadRequest {
		fmt.Println(k.Message)
		return nil
	} else {
		panic(r.Status)
	}
}

func (a *ApiClient) Create(uri string, path string) *KwkLink {
	k := &KwkLink{}
	message := fmt.Sprintf(`{"url":"%s", "key":"%s"}`, uri, path)
	r, _, _ := gorequest.New().Post(fmt.Sprintf("%shash", apiRoot)).
		//SetDebug(true).
		Set("x-kwk-key", "59364212-aeb2-4100-bf0e-c418ef230529").
		Send(message).
		EndStruct(k)

	if r.StatusCode == http.StatusBadRequest {
		fmt.Println(k.Error)
	} else if r.StatusCode != http.StatusOK {
		fmt.Println(r)
		fmt.Println(r.Status)
	}
	return k
}

