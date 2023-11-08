package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/syndtr/goleveldb/leveldb/util"
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

func GetFromByPost(cms *[]Comment, PostId string) error {
	ldb, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return nil
	}
	defer ldb.Close()

	iter := ldb.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("comment-%v", PostId))), nil)

	for iter.Next() {
		var cm Comment
		err = proto.Unmarshal(iter.Value(), &cm)
		if err != nil {
			return nil
		}
		*cms = append(*cms, cm)
	}

	iter.Release()
	err = iter.Error()
	if err != nil {
		return err
	}
	return nil
}

func (cm *Comment) SavetoSQL(frag_num int64) error {
	//mysql put
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		return err
	}
	sql := `CREATE TABLE IF NOT EXISTS Comment (
    Id  VARCHAR(250) PRIMARY KEY,
    User VARCHAR(250) NOT NULL,
    PostId VARCHAR(250) NOT NULL,
    UserId VARCHAR(250) NOT NULL,
    RepliedTo VARCHAR(250),
    Content TEXT,
    CreatedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Likes INTEGER DEFAULT 0,
    Replies INTEGER DEFAULT 0,
    IsDeleted BOOLEAN DEFAULT 0,
    IsUpdated BOOLEAN DEFAULT 0
)`
	defer db.Close()
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO Comment (Id,User,PostId,UserId,RepliedTo,Content,CreatedDate,UpdatedDate,Likes,Replies,IsDeleted,IsUpdated) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}
	cm.IsDeleted = false
	cm.IsUpdated = false
	_, err = statement.Exec(cm.Id, cm.User, cm.PostId, cm.UserId, cm.RepliedTo, cm.Content, cm.CreatedDate, cm.UpdatedDate, cm.Likes, len(cm.Replies), cm.IsDeleted, cm.IsUpdated)
	if err != nil {
		panic(err)
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

func (cm *Comment) Delete() error {
	//leveldb put
	db, err := durable.LevelDBCreateDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Delete([]byte(fmt.Sprintf("comment-%v-%v", cm.PostId, cm.Id)), nil)
	if err != nil {
		return err
	}
	return nil

}
func (cm *Comment) Update() error {
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
	err = db.Delete([]byte(fmt.Sprintf("comment-%v-%v", cm.PostId, cm.Id)), nil)
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
