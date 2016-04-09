package models

import (
	"time"
)

type Session struct{
		Token string `json:"token"`
		Device string `json:"device"`
		Created time.Time `json:"created"`
		Modified time.Time `json:"modified"`
		Online bool `json:"online"`
		Ttl time.Time `json:"ttl"`
		Active bool `json:active"`
}