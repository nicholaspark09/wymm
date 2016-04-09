package models

import "time"

type Permission struct{
	Name string `json:"name"`
	VideoKey string `json:"videokey"`
	Level int `json:"level"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Active bool `json:"active"`
}