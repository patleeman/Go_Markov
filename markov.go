package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const markov_order = 3
const corpus_dir = "./text_corpus"

func main() {
	scanned_words := scan_text()
	fmt.Println(scanned_words)
}

// Scan text files inside corpus folder and return a map of schema:
// map[word : [word -n, ..., word -2, word -1]]
func scan_text() map[string][markov_order]string {
	// Grab all text files from text corpus folder
	file_list := list_corpus(corpus_dir)

	// Grab contents from file
	var word_map = make(map[string][markov_order]string)
	for _, f := range file_list {
		file_contents := grab_contents(f)
		punct_cleaned := replace_punctuation(file_contents)
		words := strings.Fields(punct_cleaned)

		for index, word := range words {
			var prior_words [markov_order]string
			if index <= markov_order {
				continue
			} else {
				for i := 0; i < markov_order ; i++ {
					word_index := index - (markov_order - i)
					prior_words[i] = words[word_index]
					word_map[word] = prior_words
				}
			}
		}
	}
	return word_map
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

// Replace punctuation marks with a space before so field can
// interpret it as a separate entity.
func replace_punctuation(text string) string {
	lookup := make(map[string]string)
	lookup["."] = " ."
	lookup[","] = " ,"
	lookup[";"] = " ;"
	lookup[":"] = " :"
	lookup["!"] = " !"
	lookup["?"] = " ?"
	lookup["\n"] = ""
	lookup["\t"] = ""
	lookup["\r"] = ""

	var new_str string = text
	for punct, repl := range lookup {
		new_str = strings.Replace(new_str, punct, repl, -1)
	}
	return new_str
}