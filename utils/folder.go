package utils

import (
	"os"
	"path/filepath"
)

func GetWD() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
