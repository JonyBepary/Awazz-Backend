package middlewares

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/blake2b"
)

func TokenGenerator(UserId, timeStamp, Server_seed string) string {
	var seed []byte
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)

	h, err := blake2b.New384([]byte(fmt.Sprint(r.Int63())))
	if err != nil {
		panic(err)
	}
	h.Write(seed)
	h.Write([]byte(UserId))
	h.Write([]byte(timeStamp))
	h.Write([]byte(Server_seed))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash

}
