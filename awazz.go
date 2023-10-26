package main

// "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/post", getPost)
	r.POST("/post", savePost)
	r.GET("/person", getPerson)
	r.POST("/person", savePerson)
	r.POST("/community", saveCommunity)
	r.GET("/community", getCommunity)
	r.POST("/instance", saveInstance)
	r.GET("/instance", getInstance)
	r.POST("/comments", saveComment)
	r.GET("/comments", getComment)
	r.POST("/message", saveMessage)
	r.GET("/message", getMessage)
	r.POST("/follower", saveFollower)
	r.GET("/follower", getFollower)
	r.POST("/followee", saveFollowee)
	r.GET("/followee", getFollowee)
	r.POST("/notification", saveNotification)
	r.GET("/notification", getNotification)
	r.POST("/likes", saveLikes)
	r.GET("/likes", getLikes)

	// r.GET("/login", getMessage)
	// r.POST("/register", saveMessage)


	r.Run(":9091")
}
