package main

import "github.com/gin-gonic/gin"

// "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"

func main() {

	r := gin.Default()
	r.GET("/post", getPost)
	r.POST("/post", savePost)
	r.DELETE("/post", delPost)
	r.GET("/posts_n", getNpost)
	r.GET("/person", getPerson)
	r.POST("/person", savePerson)
	r.DELETE("/person", delPerson)
	r.POST("/community", saveCommunity)
	r.GET("/community", getCommunity)
	r.DELETE("/community", delCommunity)
	r.POST("/instance", saveInstance)
	r.GET("/instance", getInstance)
	r.DELETE("/instance", delInstance)
	r.POST("/comments", saveComment)
	r.GET("/comments", getComment)
	r.DELETE("/comments", delComment)
	r.POST("/message", saveMessage)
	r.GET("/message", getMessage)
	r.DELETE("/message", delMessage)
	r.POST("/follower", saveFollower)
	r.GET("/follower", getFollower)
	r.DELETE("/follower", delFollower)
	r.POST("/followee", saveFollowee)
	r.GET("/followee", getFollowee)
	r.DELETE("/followee", delFollowee)
	r.POST("/notification", saveNotification)
	r.GET("/notification", getNotification)
	r.DELETE("/notification", delNotification)
	r.POST("/like", saveLike)
	r.GET("/like", getLike)
	r.DELETE("/like", delLike)
	r.GET("/likes_entity/", getLIKESByEntityId)
	r.GET("/likes_user/", getLIKESByUserId)
	r.GET("/login", login)
	r.POST("/register", register)
	r.POST("/file", uploadFile)
	r.GET("/file", downloadFile)
	r.DELETE("/file", deleteFile)
	r.POST("/files", uploadFiles)
	r.DELETE("/files", deleteFiles)
	r.Static("/Database/assets/", "./Database/assets/")

	r.Run(":9091")
}
