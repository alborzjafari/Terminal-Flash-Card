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

func update_word(word *string, desc *string, db *sql.DB) {
  update_word_sql := `UPDATE dictionary SET desc = ? WHERE word = ?`;
  statement, err := db.Prepare(update_word_sql)
  if err != nil {
    log.Fatalln(err.Error())
  }
  _, err = statement.Exec(*desc, *word)
  if err != nil {
    log.Fatalln(err.Error())
  }
}

func remove(word *string, db *sql.DB) {
  update_word_sql := `DELETE FROM dictionary WHERE word = ?`;
  statement, err := db.Prepare(update_word_sql)
  if err != nil {
    log.Fatalln(err.Error())
  }
  _, err = statement.Exec(*word)
  if err != nil {
    log.Fatalln(err.Error())
  }
}

func main() {
  word := flag.String("w", "", "word")
  desc := flag.String("d", "", "description")
  word_to_update := flag.String("u", "", "word to update")
  database_path := flag.String("b", "dictionary.db", "database")
  remove_word := flag.String("r", "", "word for remove")
  flag.Parse()
  if !check_db(*database_path) {
    fmt.Println("DB not available")
  } else {
    db := open_db(*database_path)
    if len(*word) > 0 && len(*desc) > 0 && len(*word_to_update) == 0 {
      insert_word(word, desc, db)
    } else if len(*word_to_update) > 0 && len(*desc) > 0 && len(*word) == 0 {
      update_word(word_to_update, desc, db)
    } else if len(*remove_word) > 0 {
      remove(remove_word, db)
    } else {
      for !get_one_random_word(db) {}
    }
  }
}
