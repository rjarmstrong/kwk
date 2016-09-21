package models

import "time"

type AliasList struct {
	Items []Alias `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
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
	Updated   time.Time `json:"updated"`
}
