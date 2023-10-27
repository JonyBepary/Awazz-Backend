package pkg

import (
	"os"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
)

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
