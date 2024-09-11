package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func Getwd() string {
	// using the function
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
	return mydir
}

func ExecuteDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	return exPath
}
