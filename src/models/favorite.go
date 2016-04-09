package models

import "time"

type Favorite struct{
	Name string `json:"name"`
	Line string `json:"line"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Active bool `json:"active"`
}