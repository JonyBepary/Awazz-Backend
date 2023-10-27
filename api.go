package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/middlewares"
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "Token", "/")
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
// It takes a gin.Context object as input and retrieves the username and password from the query parameters.
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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

	ldb, err = durable.LeveldbCreateDatabase("Database/", "Token", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	post.Id = c.Param("id")
	// err := post.GetPost(post.Id)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	lbd, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	notification_id := c.Query("NotificationId")
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
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

func saveLikes(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Likes{}
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
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		panic(err)
	}
	defer ldb.Close()
	blob, err := proto.Marshal(&p)
	if err != nil {
		log.Print(err)
	}
	ldb.Put([]byte(fmt.Sprintf("like_%v", p.UserId)), blob, nil)
	ldb.Close()
	c.JSON(200, p)

}

func getLikes(c *gin.Context) {

	// check if user is logged in
	if err := check_login(c); err != nil {
		c.JSON(401, gin.H{"authentication error": err.Error()})
		return
	}

	p := model.Likes{}
	ldb, err := durable.LeveldbCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte(fmt.Sprintf("like_%v", p.UserId)), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	proto.Unmarshal(blob, &p)
	ldb.Close()
	c.JSON(200, p)
}
