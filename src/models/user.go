package models

import(
	"time"
)

type User struct{	
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	OauthProvider string `json:"oauthprovider"`
	OauthUid int64 `json:"oauthuid"`
	Token string `json:"token"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Ttl time.Time `json:"ttl"`
	Locale string `json:"locale"`
	Group int `json:"group"`
	Pic string `json:"pic"`
	Active bool `json:"active"`
}