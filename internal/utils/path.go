package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func HowManySlash(path, targetSegment string) string {
	idx := strings.LastIndex(path, targetSegment)
	if idx == -1 {
		return ""
	}

	nums := strings.Count(filepath.ToSlash(path[idx:]), "/")
	var ret string
	for i := nums; i >= 0; i-- {
		ret += "../"
	}
	return ret
}
