package pkg

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
	ulid "github.com/oklog/ulid/v2"
)


// ULID FUNCTIONS IS Universally Unique Lexicographically Sortable Identifier
// FOR FILE NAMING
func GetUlid() string {
	r := rand.New(new(rand.LockedSource))
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(r, 0)).String()
}

// read file from path to blob
func ReadFile(path string) []byte {
	// Read the contents of the file into a byte slice
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// READ PUBLIC KEY
func ReadServerPKI() (*model.AKS, error) {
	pk := &model.AKS{
		Id: "server_default",
	}
	err := pk.ReadToDatabase()
	if err != nil {
		return nil, err
	}
	return pk, nil
}

// Write Server PUBLIC KEY
func WriteServerPKI() (*model.AKS, error) {
	pk := &model.AKS{
		Id: "server_default",
	}
	err := pk.WriteToDatabase()
	if err != nil {
		return nil, err
	}
	return pk, nil
}


// HASH FUNCTIONS FOR STRING
func StringHashGeneration(str string) string {
	digest := sha256.New()
	digest.Write([]byte(str))

	return fmt.Sprintf("%x", digest.Sum(nil))
}

// HASH FUNCTIONS FOR FILE
func FileHashGeneration(filename string) string {
	file, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	digest := sha256.New()
	digest.Write(file)
	return fmt.Sprintf("%x", digest.Sum(nil))
}


// string to shard determiner
func StringToShard(s string) int64 {
	var n int64
	if s == "" {
		// UserId is missing, return 401
		n = 0
		} else {
		r := []rune(s)
		n = int64(r[0]) % 6
	}
	return n
}

func StringFullToShard(s string) int64 {
	var n int64
	if s == "" {
		// UserId is missing, return 401
		n = 0
		} else {
		for _, r := range []rune(s) {
			n += int64(r)
		}
		n = n % 6
	}
	return n
}





// func toChar(i int) rune {
// 	return rune('A' - 1 + i)
// }
