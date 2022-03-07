package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"time"
	"strconv"
	"math/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const kDbFilePath string = "dictionary.db"

func open_db(db_path string) *sql.DB {
	db_pointer, _ := sql.Open("sqlite3", db_path)
	return db_pointer
}

func get_random(max int) int {
	rand.Seed(time.Now().UnixNano())
	var min int = 1
	return rand.Intn(max) + min
}

func get_db_max(db *sql.DB) int {
	row, err := db.Query("SELECT id FROM dictionary")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var max int = 0
	for row.Next() {
		max++
	}
	return max
}

func get_one_random_word(db *sql.DB) bool {
	max := get_db_max(db)
	rand_id := get_random(max)
	row, err := db.Query("SELECT * FROM dictionary WHERE id=" +
		strconv.Itoa(rand_id))
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var rows_available bool = false
	var word string
	for row.Next() {
		var id int
		var desc string
		row.Scan(&id, &word, &desc)
    fmt.Println(word + ":", desc)
		rows_available = true
	}
	return rows_available
}

func check_db(db_path string) bool {
	info, err := os.Stat(db_path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func insert_word(word *string, desc *string, db *sql.DB) {
	insert_word_sql := `INSERT INTO dictionary(word, desc) VALUES (?, ?)`
	statement, err := db.Prepare(insert_word_sql)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(*word, *desc)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	word := flag.String("w", "", "word")
	desc := flag.String("d", "", "description")
	flag.Parse()
  if !check_db(kDbFilePath) {
    fmt.Println("DB not available")
  } else {
    db := open_db(kDbFilePath)
    if len(*word) > 0 && len(*desc) > 0 {
      insert_word(word, desc, db)
    } else {
      for !get_one_random_word(db) {}
    }
  }
}
