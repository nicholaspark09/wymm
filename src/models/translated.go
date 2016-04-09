package models

import "time"

type Translated struct{
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Controller string `json:"controller"`
	Action string `json:"action"`
	User string `json:"user"`
	Locale string `json:"locale"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Flags int `json:"flags"`
	Likes int `json:"likes"`
	Ranking float32 `json:"ranking"`
	Active bool `json:"active"`
}