package model

import (
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	cr "github.com/cloudflare/circl/sign/eddilithium3"
)

type AKS struct {
	Id      string
	Types   string
	PrivKey cr.PrivateKey
	PubKey  cr.PublicKey
}

// This Interface.takes filepath as parameter and save public key and private key in the database
func (pk *AKS) WriteToDatabase(param ...string) error {
	filepath := "Database/blob/"
	if len(param) > 1 {
		filepath = param[0]
	}
	ldb, err := durable.LevelDBCreateDatabase(filepath)
	if err != nil {
		return err
	}
	defer ldb.Close()
	blob, err := pk.PubKey.MarshalBinary()
	if err != nil {
		return err
	}
	err = ldb.Put([]byte("publicKey_"+pk.Id), blob, nil)
	if err != nil {
		return err
	}
	blob, err = pk.PrivKey.MarshalBinary()
	if err != nil {
		return err
	}
	err = ldb.Put([]byte("privateKey_"+pk.Id), blob, nil)
	if err != nil {
		return err
	}
	return nil
}

func (pk *AKS) ReadToDatabase(param ...string) error {
	filepath := "Database/blob/"
	if len(param) > 1 {
		filepath = param[0]
	}
	ldb, err := durable.LevelDBCreateDatabase(filepath)
	if err != nil {
		return err
	}
	defer ldb.Close()
	blob, err := ldb.Get([]byte("publicKey_"+pk.Id), nil)
	if err != nil {
		return err
	}
	err = pk.PubKey.UnmarshalBinary(blob)
	if err != nil {
		return err
	}
	blob, err = ldb.Get([]byte("privateKey_"+pk.Id), nil)
	if err != nil {
		return err
	}
	err = pk.PrivKey.UnmarshalBinary(blob)
	if err != nil {
		return err
	}
	pk.Types = pk.PubKey.Scheme().Name()
	return nil
}
