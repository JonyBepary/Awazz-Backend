package main

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

// "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"

func getPost(c *gin.Context) {
	var post model.Post
	post.Id = c.Param("id")
	err := post.GetPost(post.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, post)
}

func savePost(c *gin.Context) {
	var post model.Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = post.SavePost()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, post)
}
func getPerson(c *gin.Context) {
	var p model.Person
	pid := c.Query("id")
	//! println("pid: " + pid)
	p, err := model.GetPerson(pid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	spew.Dump(p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, p)
}

func savePerson(c *gin.Context) {
	person := model.Person{}
	err := c.Bind(&person)
	if err != nil {
		println("---------------------------------------------------------------------------------------------------------")
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT PERSON OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	fmt.Printf("%+v", person)
	err = person.SavePerson()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, fmt.Sprintf("%+v", person))
}

func main() {
	r := gin.Default()
	r.GET("/post", getPost)
	r.POST("/post", savePost)
	r.GET("/person", getPerson)
	r.POST("/person", savePerson)
	r.Run(":9091")
}
