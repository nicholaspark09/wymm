package models

import "time"

type Like struct{
	Name string `json:"name"`
	Controller string `json:"controller"`
	Action string `json:"action"`
	User string `json:"user"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Active bool `json:"active"`
}