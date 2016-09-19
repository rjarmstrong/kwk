package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"bufio"
	"os"
	"strconv"
	"strings"
	"net/url"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/models"
	"github.com/kwk-links/kwk-cli/libs/settings"
)

const (
	//apiRoot = "http://kwk.loc/api/v1/"
	apiRoot = "http://localhost:8080/api/v1/"
	userDbKey = "user"
)

type ApiClient struct {
	Settings settings.ISettings
}

func New(s settings.ISettings) IApi {
	return &ApiClient{Settings:s}
}

type Alias struct {
	Id        int64 `json:"id"`

	FullKey   string `json:"fullKey"`
	Username  string `json:"username"`
	Key       string `json:"key"`
	Extension string `json:"extension"`

	Uri       string `json:"uri"`
	Version   int `json:"version"`
	Runtime   string `json:"runtime"`
	Media     string `json:"media"`
	Tags      []string `json:"tags"`
	Created   time.Time `json:"created"`
	DefaultModel
}

type AliasList struct {
	Items []Alias `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type DefaultModel struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (k *Alias) Err() string {
	return k.Error
}

func (a *ApiClient) List(args []string) *AliasList {
	var page, size int
	var tags = []string{}
	for _, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			if page == 0 {
				page = num
			} else {
				size = num
			}
		} else {
			tags = append(tags, v)
		}
	}
	list := &AliasList{}
	var tagTokens []string
	for _, v := range tags {
		tagTokens = append(tagTokens, fmt.Sprintf("&tags=%s", v))
	}
	a.Request("GET", fmt.Sprintf("alias/%s?p=%d&s=%d%s", "richard", page, size, strings.Join(tagTokens, "")), "", list)
	return list
}

func (a *ApiClient) Get(fullKey string) *AliasList {
	k := &AliasList{}
	a.Request("GET", fmt.Sprintf("alias/%s/%s", "richard", url.QueryEscape(fullKey)), "", k)
	return k
}

func (a *ApiClient) Delete(key string) {
	a.Request("DELETE", fmt.Sprintf("alias/%s/%s", "richard", url.QueryEscape(key)), "", nil)
}

func (a *ApiClient) Create(uri string, path string) *Alias {
	body := struct {
		Url string `json:"url"`
		Key string `json:"key"`
	}{
		Url: uri,
		Key: path,
	}

	k := &Alias{}
	j, _ := json.Marshal(body)
	a.Request("POST", "hash", string(j), k)
	return k
}

func (a *ApiClient) Rename(key string, newKey string) *Alias {
	body := fmt.Sprintf(`{"newKey":"%s"}`, newKey)
	k := &Alias{}
	a.Request("PUT", fmt.Sprintf("alias/%s/%s/rename", "richard", url.QueryEscape(key)), body, k)
	return k
}

func (a *ApiClient) Patch(fullKey string, uri string) *Alias {
	body := struct {
		Uri string `json:"uri"`
	}{
		Uri: uri,
	}
	j, _ := json.Marshal(body)
	k := &Alias{}
	a.Request("PUT", fmt.Sprintf("alias/%s/%s/patch", "richard", url.QueryEscape(fullKey)), string(j), k)
	return k
}

func (a *ApiClient) Login(username string, password string) *models.User {
	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(gui.Colour(gui.LightBlue, "Your Kwk Username: "))
		usernameBytes, _, _ := reader.ReadLine()
		username = string(usernameBytes)
	}
	if password == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(gui.Colour(gui.LightBlue, "Your Password: "))
		passwordBytes, _, _ := reader.ReadLine()
		password = string(passwordBytes)
	}

	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	u := &models.User{}
	a.Request("POST", "users/login", body, u)
	if len(u.Token) > 50 {
		a.Settings.Upsert(userDbKey, u)
		fmt.Printf("%v signed in!", u.Username)
		return u
	}
	return nil
}

func (a *ApiClient) SignUp(email string, username string, password string) *models.User {
	body := fmt.Sprintf(`{"email":"%s", "username":"%s", "password":"%s"}`, email, username, password)
	u := &models.User{}
	a.Request("POST", "users", body, u)
	if len(u.Token) > 50 {
		a.Settings.Upsert(userDbKey, u)
		fmt.Printf("Welcome to kwk %s! You're signed in already.", u.Username)
		return u
	}
	return nil
}

func (a *ApiClient) Clone(fullKey string, newKey string) *Alias {
	return nil
}

func (a *ApiClient) Tag(fullKey string, tags ...string) *Alias {
	json, _ := json.Marshal(tags)
	body := fmt.Sprintf(`{"tags":%s}`, json)
	a.Request("POST", fmt.Sprintf("alias/%s/%s/tag", "richard", url.QueryEscape(fullKey)), body, nil)
	return nil
}

func (a *ApiClient) UnTag(fullKey string, tags ...string) *Alias {
	json, _ := json.Marshal(tags)
	body := fmt.Sprintf(`{"tags":%s}`, json)
	a.Request("DELETE", fmt.Sprintf("alias/%s/%s/tag", "richard", url.QueryEscape(fullKey)), body, nil)
	return nil
}

func (a *ApiClient) Signout() {
	a.Settings.Delete(userDbKey)
	fmt.Println("Logged out.")
}

func (a *ApiClient) PrintProfile() {
	u := &models.User{}
	err := a.Settings.Get(userDbKey, u)
	if err != nil {
		fmt.Println("You are not logged in please log in: kwk login <username> <password>")
	} else {
		fmt.Println("~~~~~~ Your Profile ~~~~~~~~~")
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(5, gui.Space) + gui.Build(5, gui.UBlock) + gui.Build(5, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))

		fmt.Printf("Email:      %v\n", u.Email)
		fmt.Printf("Username:   %v\n", u.Username)
		fmt.Printf("Host:       %v\n", u.Host)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
}

func (a *ApiClient) Request(method string, path string, body string, model interface{}) {
	url := fmt.Sprintf("%s%s", apiRoot, path)
	var req *http.Request
	if body != "" {
		b := []byte(body)
		buffer := bytes.NewBuffer(b)
		req, _ = http.NewRequest(method, url, buffer)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	u := &models.User{}
	if a.Settings.Get(userDbKey, u); u == nil {
		fmt.Println("You are not logged in.")
		return
	}

	req.Header.Set("x-kwk-key", u.Token)
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
	if model == nil {
		model = &DefaultModel{}
	}
	if e := json.Unmarshal(responseBytes, model); e != nil {
		fmt.Println(e)
		handleResponse(path, model, r)
		return
	}
	handleResponse(path, model, r)
}

func handleResponse(path string, i interface{}, r *http.Response) {
	switch {
	case r.StatusCode == http.StatusBadRequest :
		if i != nil {
			system.PrettyPrint(i)
		} else {
			fmt.Println(r.StatusCode)
		}
	case r.StatusCode == http.StatusForbidden :
		fmt.Println("Sign in please: 'kwk signin <username> <password>'")
	case r.StatusCode == http.StatusNotFound :
		fmt.Println(path + " not found.")
	case r.StatusCode != http.StatusOK :
		fmt.Println(r)
		fmt.Println(r.StatusCode)
	}
}

type ErrorResponse interface {
	Err() string
}

