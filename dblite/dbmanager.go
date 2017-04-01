package dblite

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	sqlStmt string = `CREATE TABLE IF NOT EXISTS dynasty (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  url VARCHAR(1024) NOT NULL,
  poet_count VARCHAR(1024) NOT NULL,
  poet_page VARCHAR(1024) NOT NULL
) charset = utf8mb4;
CREATE TABLE IF NOT EXISTS poet (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(1024) NOT NULL,
  uid VARCHAR(1024),
  url VARCHAR(1024) NOT NULL,
  poet_count INT,
  dynasty_id INT,
  INDEX(id),
  FOREIGN KEY(dynasty_id) REFERENCES dynasty(id)
)  charset = utf8mb4;
CREATE TABLE IF NOT EXISTS peotry_index (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  no INT,
  name VARCHAR(100) NOT NULL,
  summary VARCHAR(100),
  type VARCHAR(255),
  url VARCHAR(1024) NOT NULL,
  poetuid VARCHAR(255)
) charset = utf8mb4;
CREATE TABLE IF NOT EXISTS poetry (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  url VARCHAR(1024) NOT NULL,
  content MEDIUMTEXT,
  author VARCHAR(255),
  interpret MEDIUMTEXT,
  title VARCHAR(255),
  poetuid VARCHAR(255),
  INDEX(id)
)  charset = utf8mb4;`
)

type DBManager struct {
	db *sql.DB
}

func NewDBManager() *DBManager {
	manager := &DBManager{}
	manager.CreateTableIfNeeded()
	return manager
}

func (manager *DBManager) CreateTableIfNeeded() bool {
	var err error
	manager.db, err = sql.Open("sqlite3", "./poetry.db")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer manager.db.Close()
	_, err = manager.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return false
	}
	return true
}


func (manager *DBManager) SelectPoetry(keyword string) string {

	rows, err := manager.db.Query("select title, author, content from poetry where title = " + "'《" + keyword + "》'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		var author string
		var content string
		err = rows.Scan(&title, &author, &content)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(title, author)
		return title + "\n" + content
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}




