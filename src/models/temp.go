package models

import (
	"time"
)

type Temp struct{
		User string `json:"user"`
		Email string `json:"email"`
		Type string `json:"type"`
		Temptoken string `json:"temptoken"`
		Created time.Time `json:"created"`
		Modified time.Time `json:modified"`
		Activated bool `json:"activated"`
		Active bool `json:"active"`
}