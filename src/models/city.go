package models

import "time"

type City struct{
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
	Country string `json:"country"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	CityCount int64 `json:"citycount"`
	Active bool `json:"active"`
}