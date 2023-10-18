package pkg

import (
	"os"
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
