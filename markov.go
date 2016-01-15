package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func main() {
	// Grab all text files from text corpus folder
	file_list := list_corpus("./text_corpus")

	// Grab contents from file
	for _, f := range file_list {
		file_contents := grab_contents(f)
		words := strings.Fields(file_contents)
		for _, word := range words {
			fmt.Println(word)
		}
	}
}


// grab_contents takes a file path for a text file and returns
// all contents of the file as a single string.
func grab_contents(path string) string {
	data, err := ioutil.ReadFile(path)
	check(err)
	return string(data)
}

// Error checking helper function.
func check (e error) {
	if e != nil {
		panic(e)
	}
}

// List_corpus takes a folder path to list and returns the
// full paths of the .txt files within that folder.  It does
// not search folders recursively.
func list_corpus(folder_path string) []string {
	file_list := make([]string, 0)
	files, _ := ioutil.ReadDir(folder_path)
	for _, f := range files {
		if f.IsDir() == true {
			continue
		} else {
			extension := path.Ext(f.Name())
			if extension == ".txt"{
				file_path := path.Join(folder_path, f.Name())
				file_list = append(file_list, file_path)
			}
		}
	}
	return file_list
}

