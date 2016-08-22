package api

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
)

var apiRoot ="http://kwk.loc/api/v1/"

type ApiClient struct {

}

type decoded struct {
	Uri string `json:"url"`
	Message string `json:"message" `
}

func (a *ApiClient) Decode(kwklink string) string {
	var d decoded
	r, _, errs := gorequest.New().Get(fmt.Sprintf("%shash/%s", apiRoot, kwklink)).EndStruct(&d)
	if r.StatusCode != http.StatusOK {
		panic(r.Status)
	}
	if errs != nil {
		fmt.Println(errs)
	}
	return d.Uri
}

type kwklink struct {
	Id int64 `json:"id"`
	Key string `json:"key"`
	Root string `json:"root"`
	Url string `json:"url"`
	AfToken string `json:"afToken"`
	Error string `json:"error"`
}

func (a *ApiClient) Create(uri string, path string) kwklink {
	var k kwklink
	message := fmt.Sprintf(`{"url":"%s", "key":"%s"}`, uri, path)
	r, _, _ := gorequest.New().Post(fmt.Sprintf("%shash", apiRoot)).
		//SetDebug(true).
		Set("x-kwk-key", "59364212-aeb2-4100-bf0e-c418ef230529").
		Send(message).
		EndStruct(&k)

	if r.StatusCode == http.StatusBadRequest {
		fmt.Println(k.Error)
	} else if r.StatusCode != http.StatusOK {
		fmt.Println(r)
		fmt.Println(r.Status)
	}
	return k
}

