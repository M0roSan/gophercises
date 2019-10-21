package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	dir := "sample"
	toRename := make(map[string][]string)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename[path] = append(toRename[path], info.Name())
		}
		return nil
	})

	for _, files := range toRename {
		n := len(files)
		sort.Strings(files)
		for i, filename := range files {
			res, _ := match(filename)
			newFilename := fmt.Sprintf("%s - %d of %d.%s", res.base, i+1, n, res.ext)
			oldPath := filepath.Join(dir, filename)
			newPtah := filepath.Join(dir, newFilename)
			fmt.Printf("mv %s => %s\n", oldPath, newPtah)
			err := os.Rename(oldPath, newPtah)
			if err != nil {
				fmt.Println("Error renaming:", oldPath, newPtah, err.Error())
			}
		}
	}

}

type matchResult struct {
	base  string
	ext   string
	index int
}

// match returns new filename or an error if the filename
// didn't match our pattern.
func match(fileName string) (*matchResult, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return &matchResult{strings.Title(name), ext, number}, nil
}
