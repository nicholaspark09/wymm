package models

import "time"

type View struct{
	Name string `json:"name"`
	Controller string `json:"controller"`
	Action string `json:"action"`
	Safekey string `json:"safekey"`
	User string `json:"user"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Device string `json:"device"`
	IP string `json:"ip"`
	Host string `json:"host"`
	Active bool `json:"active"`
}