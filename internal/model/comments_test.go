package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

func TestComments_Save_Get(t *testing.T) {
	// Create a new Comments object with some data
	cm := &Comments{

		PostId: "1",
		Id:     "2",
		Text:   "This is a comment",
	}

	// Call the Save method
	cm.Save()

	// Open the database and check that the data was saved correctly
	db, err := leveldb.OpenFile("Database/comments", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	data, err := db.Get([]byte(fmt.Sprintf("comment-%v-%v", cm.PostId, cm.Id)), nil)
	if err != nil {
		panic(err)
	}

	// Unmarshal the data and check that it matches the original Comments object
	var savedCm Comments
	err = proto.Unmarshal(data, &savedCm)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, cm.GetId(), savedCm.GetId())
	assert.Equal(t, cm.GetPostId(), savedCm.GetPostId())
	assert.Equal(t, cm.GetText(), savedCm.GetText())
}

func TestGetNComments(t *testing.T) {
	// Create 10 Comments objects with different post IDs
	var comments []*Comments
	for i := 1; i <= 10; i++ {
		cm := &Comments{
			PostId: "1",
			Id:     fmt.Sprint(i),
			Text:   fmt.Sprintf("This is comment %d", i),
		}
		cm.Save()
		comments = append(comments, cm)
	}

	// Call the GetNComments function for post 1
	results, err := GetNComments("1", 10)
	assert.NoError(t, err)

	// Check that the correct comments were returned
	assert.Len(t, results, 10)
	for i := 1; i < 10; i++ {
		assert.Equal(t, comments[i].GetId(), results[i].GetId())
		assert.Equal(t, comments[i].GetText(), results[i].GetText())

	}

	// Call the Get10Comments function for post 2

}
