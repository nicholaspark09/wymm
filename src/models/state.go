package models

import "time"

type State struct{
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Country string `json:"country"`
	Created time.Time `json:"created"`
	Active bool `json:"active"`
}