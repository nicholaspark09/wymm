package models

import "time"

type Country struct{
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Created time.Time `json:"created"`
	Active bool `json:"active"`
}