package durable

import (
	"fmt"
	"log"
	"strconv"
)

// // Add adds a fragment to the fragmentation.
func Fragmentation_Add(shard int64, KEY ...string) {
	ldb, err := LevelDBCreateDatabase("Database/", "GlobalSchema", "/")
	if err != nil {
		log.Print(err)
		return
	}
	for _, key := range KEY {
		err = ldb.Put([]byte(key), []byte(fmt.Sprintf("%d", shard)), nil)
		if err != nil {
			log.Print(err)
			return
		}
	}
	defer ldb.Close()

}

// Remove removes a fragment from the fragmentation.
func Fragmentation_Remove(key string) {
	ldb, err := LevelDBCreateDatabase("Database/", "GlobalSchema", "/")
	if err != nil {
		log.Print(err)
		return
	}
	err = ldb.Delete([]byte(key), nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer ldb.Close()
}

// // Get returns the shard for a given rowID.
func Fragmentation_Get(rowID string) (int64, error) {
	ldb, err := LevelDBCreateDatabase("Database/", "GlobalSchema", "/")
	if err != nil {
		log.Print(err)
		return -1, fmt.Errorf("error in creating database")
	}
	defer ldb.Close()
	shard, err := ldb.Get([]byte(rowID), nil)
	if err != nil {
		log.Print(err)
		return -1, fmt.Errorf("error in getting shard")
	}
	// conver shard to int64
	n, err := strconv.ParseInt(string(shard), 10, 64)
	if err != nil {
		log.Print(err)
		return -1, fmt.Errorf("error in converting shard to int64")
	}

	return n, nil
}

// // FRAGMENTATION TO FILE MARSHALJSON
// func Fragmentation_SaveToFile(rowID string) error {
// 	blob, err := f.fragments.MarshalJSON()
// 	if err != nil {
// 		log.Print(err)
// 		return err
// 	}
// 	// save blob to file
// 	err = os.WriteFile("Database/global_schema.json", blob, 0666)
// 	if err != nil {
// 		log.Print(err)
// 		return err
// 	}
// 	return nil

// }

// // get fragmentation from file
// func Fragmentation_LoadFromFile() error {
// 	// read blob from file
// 	ldb, err := durable.LevelDBCreateDatabase("Database/Common/", "global_schema")
// 	if err != nil {
// 		log.Print(err)
// 		return err
// 	}
// 	// unmarshal blob to fragmentation
// 	err = fragments.UnmarshalJSON(blob)
// 	if err != nil {
// 		log.Print(err)
// 		return err
// 	}
// 	return nil
// }
