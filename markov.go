package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	db "github.com/patleeman/Go_Markov/database"
	"time"
)

const markov_order = 3
const corpus_dir = "./text_corpus"

// TODO: Fix scan_text data storage schema.  Need to nest map inside a slice in order to not overwrite keys or organize it differently.

func main() {
	start := time.Now()
	// Initialize Database
	db.InitDB(markov_order)

	// Scan Text
	fmt.Println("Scanning files in corpus directory.")
	scanned_words := scan_text()

	// Save Data to Database
	fmt.Println("Saving data to database.")
	save_to_db(scanned_words)

	elapsed := time.Since(start)
	fmt.Printf("Done.  Runtime: %s", elapsed)
}

// Scan text files inside corpus folder and return a slice with maps of schema:
// map[word : [word -n, ..., word -2, word -1]]
func scan_text() map[string][markov_order]string {
	// Grab all text files from text corpus folder
	file_list := list_corpus(corpus_dir)

	// Grab contents from file

	// map with a string type key with a markov_order sized string array
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
			fmt.Println(index, " ", word, " ",  prior_words)
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
	lookup := map[string]string{
		"." : " . ",
		"," : " , ",
		":" : " : ",
		";" : " ; ",
		"!" : " ! ",
		"?" : " ? ",
		"\n" : " ",
		"\t" : " ",
		"\r" : "",
		"--" : " -- ",
		"(" : " ( ",
		")" : " ) ",
		"'" : "''",   // Replace single quotes with double quote for db.
	}

	var new_str string = text
	for punct, repl := range lookup {
		new_str = strings.Replace(new_str, punct, repl, -1)
	}
	return new_str
}

func save_to_db(word_dict map[string][markov_order]string) {
	// Generate values statement i.e.
	// INSERT INTO markov (target, wordminus1) VALUES (xxxxxxxxxx)
	base_statement := db.GenInsert(markov_order) + " ("
	var value_stmt string
	for target := range word_dict {
		value_stmt = ""
		value_stmt += "'" + target + "', "

		data_values := word_dict[target]
		for index, nth_order_word := range data_values {

			value_stmt += "'" + nth_order_word + "'"

			if index != markov_order - 1 {
				value_stmt += ", "
			} else {
				value_stmt += ")"
			}
		}
		full_stmt := base_statement + value_stmt
		db.Execute(full_stmt)
	}
}

