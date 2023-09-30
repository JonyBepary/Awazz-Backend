package model

import (
	ac "github.com/go-ap/activitypub"
)

type User struct {
	ac.Actor
}

func (user *User) CreateUser() string {
	return "user created by saria san"
}
