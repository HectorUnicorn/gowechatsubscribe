package dblite

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"regexp"
	"github.com/astaxie/beego"
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
	username := beego.AppConfig.String("dbusername")
	password := beego.AppConfig.String("dbpassword")
	database := beego.AppConfig.String("dbdatabase")

	var err error
	login := fmt.Sprintf("%s:%s@%s", username, password, database)
	fmt.Println("login command:", login)
	manager.db, err = sql.Open("mysql", login)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	//	_, err = manager.db.Exec(sqlStmt)
	//	if err != nil {
	//		log.Printf("%q: %s\n", err, sqlStmt)
	//		return false
	//	}
	return true
}

func (manager *DBManager) SelectPoetry(keyword string) string {

	rows, err := manager.db.Query("SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
		"FROM poetry " +
		"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
		"WHERE poetry.title = " + "'《" + keyword + "》'  " +
		"ORDER BY poet.poet_count desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var title string
	var author string
	var content string
	var poetUid string

	for rows.Next() {
		err = rows.Scan(&title, &author, &content, &poetUid)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(title, author)
		return manager.packFiledsAsString(title, author,content, poetUid)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	// not found exactly fuzzy match the title.
	if len(keyword) >= 2 {
		query := "SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
			"FROM poetry " +
			"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
			"WHERE poetry.title LIKE " + "'%" + keyword + "%' " +
			"ORDER BY poet.uid ASC, poet.poet_count DESC"
		fmt.Println("query:", query)
		err := manager.db.QueryRow(query).Scan(&title, &author, &content, &poetUid);
		if err != nil {
			log.Println(err)
		}
		fmt.Println("fuzzy title:", title, author)
		if len(title) > 0 {
			return manager.packFiledsAsString(title, author,content, poetUid)
		}

	}

	// title not found, fuzzy match the poetry content
	if len(keyword) >= 2 {
		query := "SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
			"FROM poetry " +
			"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
			"WHERE poetry.content LIKE " + "'%" + keyword + "%' " +
			"ORDER BY poet.uid ASC, poet.poet_count DESC"
		fmt.Println("query:", query)
		err := manager.db.QueryRow(query).Scan(&title, &author, &content, &poetUid);
		if err != nil {
			log.Println(err)
		}
		fmt.Println("fuzzy content:", title, content)
		return manager.packFiledsAsString(title, author, content, poetUid)
	}

	return "很抱歉，还没有这首诗哦~"
}

func (manager *DBManager)packFiledsAsString(title, author, content, poetUid string) string {
	if len(content) > 0 {
		content = strings.Replace(content, "。", "。\n", -1)
		content = strings.Replace(content, ";", ";\n", -1)
		content = strings.Replace(content, "；", "；\n", -1)
		content = strings.Replace(content, "？", "？\n", -1)
		content = strings.Replace(content, "?", "?\n", -1)
		content = strings.Replace(content, "!", "!\n", -1)
		content = strings.Replace(content, "！", "！\n", -1)
		var match []string
		reg := regexp.MustCompile("（(.*?)）")
		fmt.Println("reg1:",reg.FindAllString(content, -1))
		reg = regexp.MustCompile("\\((.*?)\\)")
		match = reg.FindAllString(content, -1)
		if len(match) > 0 {
			for _, piece := range match {
				content = strings.Replace(content, piece, "", -1)
			}
		}
		fmt.Println("reg2:",match)
		reg = regexp.MustCompile("\\[(.*?)\\]")
		match = reg.FindAllString(content, -1)
		if len(match) > 0 {
			for _, piece := range match {
				content = strings.Replace(content, piece, "", -1)
			}
		}
		fmt.Println("reg3:",match)

	}
	return title + "\n" + manager.SelectDynasty(strings.Split(poetUid, "_")[0]) + " · " + author + "\n" + content
}


func (manager *DBManager) SelectDynasty(id string) string {
	var name string
	query := "SELECT name " +
		"FROM dynasty " +
		"WHERE id = " + id
	fmt.Println("query:", query)
	err := manager.db.QueryRow(query).Scan(&name);
	if err != nil {
		log.Println(err)
	}
	return name

}