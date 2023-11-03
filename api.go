package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/middlewares"
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
	"github.com/SohelAhmedJoni/Awazz-Backend/pkg"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

// set expire time to 1 day
var expireTime int64 = 3600 * 24

// check_login checks if the user is logged in.
func check_login(c *gin.Context) error {

	UserId := c.GetHeader("Anon-User")
	token := c.GetHeader("Anon-Token")

	if token == "" {
		// Token is missing, return 401
		return errors.New("missing Anon-Token")
	}
	if UserId == "" {
		// UserId is missing, return 401
		return errors.New("missing Anon-User")
	}
	ldb, err := durable.LevelDBCreateDatabase("Database/", "Token", "/")
	if err != nil {
		return err
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("token_%v", UserId)), nil)
	if err != nil {
		return err
	}
	t := model.Token{}
	err = proto.Unmarshal(blob, &t)
	if err != nil {
		return err
	}

	if token != t.Token {
		return errors.New("token mismatch")
	}
	remainingTime := t.GetGenerateTime() + expireTime - time.Now().Unix()
	if remainingTime < 0 {
		// remove token from database
		err = ldb.Delete([]byte(fmt.Sprintf("token_%v", UserId)), nil)
		if err != nil {
			return errors.New("token expired and token removal failed")
		}

		return errors.New("token expired and removed from database")
	}
	return nil
}

// login handles user login requests.
// It takes a gin.Context object as input and retrieves the username and password from the query Queryeters.
// If either the username or password is missing, it returns a 401 error.
// It then checks if the user exists in the database and retrieves the user object.
// If the password is incorrect, it returns a 401 error.
// If the user exists and the password is correct, it generates a token for the user and saves it in the database.
// It returns the generated token in the response body.
func login(c *gin.Context) {

	// function implementation
	username := c.Query("username")
	password := c.Query("password")
	if username == "" {
		// UserId is missing, return 401
		c.JSON(401, gin.H{"error": "missing user name"})
		return
	}
	if password == "" {
		// UserId is missing, return 401
		c.JSON(401, gin.H{"error": "missing password"})
		return
	}
	// check if user exists
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("user_%v", username)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := model.User{}
	err = proto.Unmarshal(blob, &user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user.Password != password {
		c.JSON(401, gin.H{"error": "wrong password"})
		return
	}
	ldb.Close()

	ldb, err = durable.LevelDBCreateDatabase("Database/", "Token", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	token := model.Token{}
	token.UserName = username
	token.GenerateTime = time.Now().Unix()
	token.Token = middlewares.TokenGenerator(token.UserName, string(token.GenerateTime), "Awazz")
	err = token.SaveToken()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	blob, err = proto.Marshal(&token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = ldb.Put([]byte("token_"+username), blob, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ldb.Close()

	c.JSON(200, token)
}

// register is a handler function for the registration endpoint.
// It checks if the user is logged in, binds the user object to the gin context,
// checks if the user exists, saves the user data, and puts the user object in leveldb.
// It returns an error if any of the above steps fail.
func register(c *gin.Context) {

	// function implementation
	u := model.User{}
	//binding the user object to the gin context
	err := c.Bind(&u)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT user OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	spew.Config.Indent = "\t"
	spew.Dump(u)

	if u.UserName == "" {
		// UserId is missing, return 401
		c.JSON(401, gin.H{"error": "missing user name"})
		return
	}
	if u.Password == "" {
		// UserId is missing, return 401
		c.JSON(401, gin.H{"error": "missing password"})
		return
	}

	// check if user exists
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	_, err = ldb.Get([]byte(fmt.Sprintf("user_%v", u.UserName)), nil)
	if err == nil {
		c.JSON(401, gin.H{"error": "username already exists"})
		return
	}
	u.AccountTime = time.Now().Unix()

	err = u.SaveUserData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// putting user object in leveldb
	blob, err := proto.Marshal(&u)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = ldb.Put([]byte(fmt.Sprintf("user_%v", u.UserName)), blob, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ldb.Close()

	c.JSON(200, u)

}

// getPost function gets the post object from the LevelDB database and returns the post object.
func getPost(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var post model.Post
	post.Id = c.Query("Id")
	// err := post.GetPost(post.Id)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("post_%v", post.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = proto.Unmarshal(blob, &post)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, post)
}

// delPost function deletes the post object from the database and returns the post object.
func delPost(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var post model.Post
	post.Id = c.Query("Id")
	err := post.DeletePost(post.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("post_%v", post.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Post Deleted")
}

// savePost function saves the post object to the LevelDB database and returns the saved post object.
func savePost(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

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
	lbd, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer lbd.Close()
	blob, err := proto.Marshal(&post)
	if err != nil {
		log.Print(err)
	}
	lbd.Put([]byte(fmt.Sprintf("post_%v", post.Id)), blob, nil)
	c.JSON(200, post)
}

// getPerson function gets the person object from the LevelDB database and returns the person object.
func getPerson(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Person
	pid := c.Query("Id")
	println("pid: " + pid)
	// err := p.GetPerson(pid)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	//! spew.Dump(p)
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("person_%v", pid)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, p)
}

func delPerson(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var person model.Person
	person.Id = c.Query("Id")
	err := person.DeletePerson(person.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("person_%v", person.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Person Deleted")
}

// savePerson function saves the person object to the LevelDB database and returns the saved person object.
func savePerson(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	person := model.Person{}
	err := c.Bind(&person)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT community OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	spew.Config.Indent = "\t"
	err = person.SavePerson()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// getting array of string from QueryArray
	person.Replies = c.QueryArray("Replies")
	person.Tag = c.QueryArray("Tag")
	person.Url = c.QueryArray("Url")
	person.Too = c.QueryArray("Too")
	person.Bto = c.QueryArray("Bto")
	person.Cc = c.QueryArray("Cc")
	person.Bcc = c.QueryArray("Bcc")
	person.Following = c.QueryArray("Following")
	person.Followers = c.QueryArray("Followers")
	person.Liked = c.QueryArray("Liked")

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()

	blob, err := proto.Marshal(&person)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("person_%v", person.Id)), blob, nil)
	c.JSON(200, fmt.Sprintf("%+v", person))
}

// saveCommunity function saves the community object to the LevelDB database and returns the saved community object.
// It takes a gin.Context object as input and binds the community object to it. It then creates a LevelDB database and saves the community object to it.
// It also saves the admin, mod, and member IDs of the community to the database. Finally, it returns the saved community object.
func saveCommunity(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	community := model.Community{}
	err := c.Bind(&community)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT community OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	spew.Config.Indent = "\t"
	spew.Dump(community)

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()

	aId := c.QueryArray("admin")
	blob, err := json.Marshal(aId)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("community_admin_%v", community.Id)), blob, nil)
	modId := c.QueryArray("mod")
	blob, err = json.Marshal(modId)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("community_mod_%v", community.Id)), blob, nil)
	memId := c.QueryArray("member")
	blob, err = json.Marshal(memId)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("community_member_%v", community.Id)), blob, nil)

	blob, err = proto.Marshal(&community)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("community_%v", community.Id)), blob, nil)

	err = community.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, fmt.Sprintf("%+v", community))
}

func getCommunity(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Community
	cid := c.Query("id")
	//! println("pid: " + pid)
	// err := p.GetCommunity(cid)
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("community_%v", cid)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	spew.Dump(p)

	c.JSON(200, p)
}
func delCommunity(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
	}
	var community model.Community
	community.Id = c.Query("Id")
	err := community.DeleteCommunity(community.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("community_%v", community.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Community Deleted")
}

// saveCommunity function saves the community object to the LevelDB database and returns the saved community object.
// It takes a gin.Context object as input and binds the community object to it. It then creates a LevelDB database and saves the community object to it.
// It also saves the admin, mod, and member IDs of the community to the database. Finally, it returns the saved community object.
func saveInstance(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Instance{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// getting array of string from QueryArray
	p.Owner = c.QueryArray("Owner")
	p.CreatedBy = c.QueryArray("CreatedBy")
	p.CommunityIds = c.QueryArray("CommunityIds")
	p.Labels = c.QueryArray("Labels")
	p.PublicDomain = c.QueryArray("PublicDomain")
	p.Tags = c.QueryArray("Tags")

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("instance_%v", p.Id)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

// func saveCommunity(c *gin.Context) {
func getInstance(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Instance
	Iid := c.Query("Id")
	//! println("pid: " + pid)
	// err := p.GetCommunity(Iid)
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("instance_%v", Iid)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb.Close()
	c.JSON(200, p)
}
func delInstance(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})

	}
	var instance model.Instance
	instance.Id = c.Query("Id")
	err := instance.DeleteInstance(instance.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("instance_%v", instance.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Instance Deleted")
}

func saveComment(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Comment{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	spew.Config.Indent = "\t"
	spew.Dump(p)

	// getting array of string from QueryArray
	p.Replies = c.QueryArray("Replies")
	p.IsDeleted = c.Query("IsDeleted") == "true" || false
	p.IsDeleted = c.Query("IsUpdated") == "true" || false
	// Save instance in sqlite database
	err = p.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, p)
}

func getComment(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Comment
	p.Id = c.Query("cid")
	p.PostId = c.Query("pid")
	//! println("pid: " + pid)
	// err := p.GetCommunity(Iid)
	err := p.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, p)
}

func delComment(c *gin.Context) {
	var comment model.Comment
	comment.Id = c.Query("cid")
	comment.PostId = c.Query("pid")
	err := comment.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("comment_%v", comment.Id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Comment Deleted")
}

func saveMessage(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Messages{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.SaveMessages()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("message_%v", p.MsgId)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

func getMessage(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Messages
	msg_id := c.Query("MsgId")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("message_%v", msg_id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb.Close()
	c.JSON(200, p)
}

func delMessage(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Messages
	p.MsgId = c.Query("MsgId")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("message_%v", p.MsgId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ldb.Close()

	err = p.DeleteMessages(p.MsgId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Message Deleted")
}

func saveNotification(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Notifications{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.SaveNotifications()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("notification_%v", p.Source)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

func getNotification(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Notifications
	notification_id := c.Query("Receiver")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("notification_%v", notification_id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb.Close()
	c.JSON(200, p)
}
func delNotification(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Notifications
	notification_id := c.Query("NotificationId")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("notification_%v", notification_id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ldb.Close()

	err = p.DeleteNotifications(notification_id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Notification Deleted")
}

func saveFollower(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Follower{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.SaveFollower()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("follower_%v", p.UserId)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

func getFollower(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	var p model.Follower
	user_id := c.Query("UserId")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("follower_%v", user_id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb.Close()
	c.JSON(200, p)
}
func delFollower(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})

	}
	var follower model.Follower
	follower.UserId = c.Query("UserId")
	err := follower.DeleteFollowee(follower.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("follower_%v", follower.UserId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Followee Deleted")
}
func saveFollowee(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Followee{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.SaveFollowee()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("followee_%v", p.UserId)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

func getFollowee(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Followee{}
	user_id := c.Query("UserId")
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("followee_%v", user_id)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = proto.Unmarshal(blob, &p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb.Close()
	c.JSON(200, p)
}

func delFollowee(c *gin.Context) {
	var followee model.Followee
	followee.UserId = c.Query("UserId")
	err := followee.DeleteFollower(followee.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("followee_%v", followee.UserId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Followee Deleted")
}

func uploadFile(c *gin.Context) {
	/*
	  UploadFile function handles the upload of a single file.
	  It gets the file from the form data, saves it to the defined path,
	  generates a unique identifier for the file, saves the file metadata to the database,
	  and returns a success message and the file metadata.
	*/
	// Get the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ulid := pkg.GetUlid()
	dir := filepath.Join("Database", "assets")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file directory"})
			return
		}
	}
	// Define the path where the file will be saved
	filePath := filepath.Join(dir, ulid)
	// Save the file to the defined path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Save file metadata to database
	fileMetadata := model.File{
		Uuid:      ulid,
		Name:      file.Filename,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: 0,
		Hash:      pkg.FileHashGeneration(filePath),
		MimeType:  file.Header.Get("Content-Type"),
		Ext:       filepath.Ext(file.Filename),
	}
	err = fileMetadata.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open Database directory"})
		return
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&fileMetadata)
	if err != nil {
		log.Print(err)
	}

	err = ldb.Put([]byte(fmt.Sprintf("file_%v", fileMetadata.Uuid)), blob, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	ldb.Close()

	// Return a success message and the file metadata
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "Details": fileMetadata})
}
func deleteFile(c *gin.Context) {
	/*
		DeleteFile function handles the deletion of a single file.
		It gets the unique identifier of the file to be deleted,
		deletes the file from the defined path, deletes the file metadata from the database,
		and returns a success message.
	*/
	// Get the unique identifier of the file to be deleted
	ulid := c.Query("ulid")
	f := model.File{
		Uuid: ulid,
	}
	// Define the path of the file to be deleted
	filePath := filepath.Join("Database", "assets", ulid)
	spew.Dump(filePath)
	spew.Dump(ulid)
	// Delete the file from the defined path
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "1Failed to delete file"})
		return
	}
	err := f.Delete(ulid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "3Failed to open Database directory"})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "3Failed to open Database directory"})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("file_%v", ulid)), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "4Failed to delete file metadata"})
		return
	}
	ldb.Close()
	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
func uploadFiles(c *gin.Context) {
	/*
		UploadFiles function handles the upload of multiple files.
		It gets the files from the form data, saves each file to the defined path,
		generates a unique identifier for each file, saves the file metadata to the database,
		and returns a success message and the file metadata.
	*/
	// Get the files from the form data
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	var fileMetadata model.FileList
	// Save each file to the defined path and generate a unique identifier for each file
	dir := filepath.Join("Database", "assets")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file directory"})
			return
		}
	}
	for i, file := range files {
		filePath := filepath.Join(dir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		fileMetadata.Files[i] = &model.File{
			Uuid:      pkg.GetUlid(),
			Name:      file.Filename,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: 0,
			Hash:      pkg.FileHashGeneration(filePath),
			MimeType:  file.Header.Get("Content-Type"),
			Ext:       filepath.Ext(file.Filename),
		}
		// Save file metadata to database
		fileMetadata.Files[i].Save()
	}
	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully", "Details": fileMetadata})
}

func deleteFiles(c *gin.Context) {
	/*
		DeleteFiles function handles the deletion of multiple files.
		It gets the unique identifiers of the files to be deleted,
		deletes the files from the defined path, deletes the file metadata from the database,
		and returns a success message.
	*/
	// Get the unique identifiers of the files to be deleted
	ulids := c.QueryArray("ulids")
	// Define the path of the files to be deleted
	for _, ulid := range ulids {
		filePath := filepath.Join("Database", "assets", ulid)
		// Delete the files from the defined path
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
			return
		}
		// Delete file metadata from database
		// REMOVE FILE USING SYSTEM
		err := os.Remove(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
			return
		}

		ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open Database directory"})
			return
		}
		defer ldb.Close()
		err = ldb.Delete([]byte(fmt.Sprintf("file_%v", ulid)), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file metadata"})
			return
		}
		ldb.Close()
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Files deleted successfully"})
}

// DownloadFile function handles the download of a single file.
// It gets the file metadata from the database, gets the file from the defined path,
// and returns the file.
func downloadFile(c *gin.Context) {
	/*
	   GetFile function retrieves a file from the server.
	   It gets the unique identifier of the file to be retrieved,
	   retrieves the file metadata from the database,
	   defines the path of the file to be retrieved,
	   opens the file, reads the first 512 bytes of the file to determine its content type,
	   gets the file info, sets the headers for the file transfer, and returns the file.
	*/
	// Get the unique identifier of the file to be retrieved
	ulid := c.Query("id")
	var file model.File
	// Retrieve the file metadata from the database
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open Database directory"})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("file_%v", ulid)), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file metadata" + err.Error() + ulid})
		return
	}
	err = proto.Unmarshal(blob, &file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file metadata" + err.Error() + ulid})
		return
	}

	// Define the path of the file to be retrieved
	filePath := filepath.Join("Database", "assets", file.GetUuid())
	// redirect to file
	c.Redirect(http.StatusMovedPermanently, filePath)

}

func saveLike(c *gin.Context) {

	// check if user is logged in
	// if err := check_login(c); err != nil {
	// 	c.JSON(401, gin.H{"authentication error": err.Error()})
	// 	return
	// }

	p := model.Like{}
	err := c.Bind(&p)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// PRINT instance OBJECT TO CONSOLE INTENDED FOR DEBUGGING
	// spew.Config.Indent = "\t"
	// spew.Dump(p)

	// Save instance in sqlite database
	err = p.SaveLikes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// CREATING LEVELDB DATABASE
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("like_%v_%v)", p.UserId, p.EntityId)), blob, nil)
	ldb.Put([]byte(fmt.Sprintf("like_%v_%v)", p.EntityId, p.UserId)), blob, nil)
	ldb.Close()
	c.JSON(200, p)
}

func getLike(c *gin.Context) {

	// check if user is logged in
	// if err := check_login(c); err != nil {
	// 	c.JSON(401, gin.H{"authentication error": err.Error()})
	// 	return
	// }

	UserID := c.Query("UserId")
	EntityId := c.Query("EntityId")
	if EntityId == "" || UserID == "" {
		c.JSON(500, gin.H{"error": errors.New("Query is empty")})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("like_%v_%v", UserID, EntityId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var p *model.Like
	proto.Unmarshal(blob, p)
	c.JSON(200, p)
}
func delLike(c *gin.Context) {
	UserID := c.Query("UserId")
	EntityId := c.Query("EntityId")
	if EntityId == "" || UserID == "" {
		c.JSON(500, gin.H{"error": errors.New("Query is empty")})
		return
	}

	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	err = ldb.Delete([]byte(fmt.Sprintf("like_%v_%v", UserID, EntityId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	p := model.Like{
		UserId:   UserID,
		EntityId: EntityId,
	}
	err = p.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Like deleted successfully"})
}

func getLIKESByUserId(c *gin.Context) {
	/*
		likeEntity function handles the liking of a Entity.
		It gets the Entity ID from the form data, gets the user ID from the session,
		checks if the user has already liked the Entity, and returns a success message.
	*/
	// // check if user is logged in
	// if err := check_login(c); err != nil {
	// 	c.JSON(401, gin.H{"authentication error": err.Error()})
	// 	return
	// }
	// Get the Entity ID from the form data
	user_id := c.Query("UserId")
	// Get the user ID from the session
	likes := model.Likes{}
	err := likes.GetByUserId(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"Details": likes.Likes})
}

func getLIKESByEntityId(c *gin.Context) {
	/*
		likeEntity function handles the liking of a Entity.
		It gets the Entity ID from the form data, gets the user ID from the session,
		checks if the user has already liked the Entity, and returns a success message.
	*/
	// check if user is logged in
	// if err := check_login(c); err != nil {
	// 	c.JSON(401, gin.H{"authentication error": err.Error()})
	// 	return
	// }
	// Get the Entity ID from the form data
	entity_id := c.Query("EntityId")
	// Get the user ID from the session
	likes := model.Likes{}
	err := likes.GetByEntityId(entity_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Make a copy of the likes object before passing it to the JSON response
	// likesCopy := make([]*model.Like, len(likes.Likes))
	// copy(likesCopy, likes.Likes)

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"Details": likes.Likes})

}

func setLikeByEntityId(c *gin.Context) {
	/*
		likeEntity function handles the liking of a Entity.
		It gets the Entity ID from the form data, gets the user ID from the session,
		checks if the user has already liked the Entity, and returns a success message.
	*/
	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}
	// Get the Entity ID from the form data
	entity_id := c.Query("entity_id")
	user_id := c.Query("user_id")
	// Get the user ID from the session
	v := model.Like{EntityId: entity_id, UserId: user_id, CreatedAt: time.Now().Unix()}
	err := v.SaveLikes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Entity liked successfully", "Details": v})
}
