package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"google.golang.org/protobuf/proto"
)

// Comments is the database model for comments.
func (cm *Comment) Get() error {
	//leveldb get
	db, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()
	data, err := db.Get([]byte(fmt.Sprintf("comment-%v-%v", cm.PostId, cm.Id)), nil)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, cm)
	if err != nil {
		return err
	}
	return nil
}

func (cm *Comment) Save() error {
	//leveldb put
	db, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()
	data, err := proto.Marshal(cm)
	if err != nil {
		return err
	}
	err = db.Put([]byte(fmt.Sprintf("comment-%v-%v", cm.PostId, cm.Id)), data, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetNComments(postId string, N int) ([]*Comment, error) {
	//leveldb get
	db, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	iterator := db.NewIterator(nil, nil)
	defer iterator.Release()
	var comments []*Comment
	for iterator.Next() && len(comments) < N+1 {
		var cm Comment
		err = proto.Unmarshal(iterator.Value(), &cm)
		if err != nil {
			return nil, err
		}
		if cm.PostId == postId {
			comments = append(comments, &cm)
		}

	}
	return comments, nil
}
