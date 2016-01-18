package main

import (
	"fmt"
	db "github.com/patleeman/Go_Markov/database"
	"math/rand"
	"time"
	"strings"
)


// Generate a random sentence.
func main(){
	stmt := "SELECT target FROM markov WHERE targetminus3='%s' AND targetminus2='%s' AND targetminus1='%s';"
	sentence := []string{"there", "was", "once"}
	for {
		x := len(sentence)
		stmt_seed := fmt.Sprintf(stmt, sentence[x-3], sentence[x-2], sentence[x-1])
		word_options := db.Query(stmt_seed, 3)
		rand_word := choose_rand(word_options)

		sentence = append(sentence, rand_word)
		if strings.ContainsAny(rand_word, "!?.") {
			break
		}
		fmt.Println(sentence)
	}
}

// Choose a random word from slice and return it
func choose_rand(word_list []string) string {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	var word_list_length int
	word_list_length = len(word_list)
	rand_choice := r1.Intn(word_list_length)
	return word_list[rand_choice]
}