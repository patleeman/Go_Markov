package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	db "github.com/patleeman/Go_Markov/database"
	"time"
	"strconv"
)

const markov_order = 3
const corpus_dir = "./text_corpus"


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
func scan_text() [][markov_order+1]string {
	// Grab all text files from text corpus folder
	file_list := list_corpus(corpus_dir)

	// Create a fixed length array of size markov_order + 1
	// Data := [word, word-1, word-2, ..., word-n]
	all_words := [][markov_order+1]string{}
	for _, f := range file_list {
		file_contents := grab_contents(f)
		punct_cleaned := replace_punctuation(file_contents)
		words := strings.Fields(punct_cleaned)

		for index, word := range words {
			var word_set [markov_order+1]string
			word_set[0] = word
			if index <= markov_order {
				continue
			} else {
				for i := 0; i < markov_order ; i++ {
					word_index := index - (markov_order - i)
					word_set[i+1] = words[word_index]
				}
			}
			all_words = append(all_words, word_set)
		}
	}
	return all_words
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

// Generate the rest of the insert statement.  Base statement generated with GenInsert.
func save_to_db(values [][markov_order+1]string) {
	// Generate values statement i.e.
	// INSERT INTO markov (target, wordminus1) VALUES (xxxxxxxxxx)
	base_statement := GenInsert(markov_order) + " ("
	var value_stmt string
	var all_statements []string
	for _, word_set := range values {
		value_stmt = ""
		for index, word := range word_set {
			value_stmt += "'" + word + "'"

			if index != markov_order {
				value_stmt += ", "
			} else {
				value_stmt += ")"
			}
		}
		full_stmt := base_statement + value_stmt
		all_statements = append(all_statements, full_stmt)
	}
	db.ExecuteTransaction(all_statements)
}

// Generate insert statements without values.
func GenInsert(markov_order int) string {
	stmt := `INSERT INTO markov (target, %s) VALUES `
	var variable_columns string
	var col_name string
	for i := 0; i < markov_order; i++{
		col_name = "targetminus" + strconv.Itoa(markov_order - i)
		if i != markov_order - 1 {
			col_name += ", "
		}
		variable_columns += col_name
	}
	stmt = fmt.Sprintf(stmt, variable_columns)
	return stmt
}