package models

import "time"

type Language struct{
	Name string `json:"name"`
	Locale string `json:"locale"`
	Native string `json:"native"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Active bool `json:"active"`
}