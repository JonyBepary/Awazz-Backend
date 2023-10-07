package model

import (
	ac "github.com/go-ap/activitypub"
)

type Account struct {
	ac.Actor
	Password string
	Email string
	Bio string

}

