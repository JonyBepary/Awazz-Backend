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

func StringHashGeneration(str string) string {
	digest := sha256.New()
	digest.Write([]byte(str))
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func FileHashGeneration(filename string) string {
	hash := sha256.New()
	file, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	hash.Write([]byte(fmt.Sprintf("%v", file)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
