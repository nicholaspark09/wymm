package models

import "time"

type Line struct{
	Name string `json:"name"`
	Translated string `json:"translated"`
	English string `json:"english"`
	Likes int32 `json:"likes"`
	Dislikes int32 `json:"dislikes"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	User string `json:"user"`
	Locale string `json:"locale"`
	Flags int `json:"flag"`
	Image string `json:"image"`
	Views int64 `json:"views"`
	Rank float32 `json:"rank"`
	Active bool `json:"active"`
}