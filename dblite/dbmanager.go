package dblite

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"regexp"
	"github.com/astaxie/beego"
	"github.com/going/toolkit/log"
	"gowechatsubscribe/models"
)

const (
	sqlStmt string = `CREATE TABLE IF NOT EXISTS dynasty (
  id INTEGER PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  url VARCHAR(128) NOT NULL,
  poet_count VARCHAR(128) NOT NULL,
  poet_page VARCHAR(128) NOT NULL
) charset = utf8mb4;
CREATE TABLE IF NOT EXISTS poet (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  uid VARCHAR(128),
  url VARCHAR(128) NOT NULL,
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
  type VARCHAR(128),
  url VARCHAR(128) NOT NULL,
  poetuid VARCHAR(128)
) charset = utf8mb4;
CREATE TABLE IF NOT EXISTS poetry (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  url VARCHAR(128) NOT NULL,
  content MEDIUMTEXT,
  author VARCHAR(128),
  interpret MEDIUMTEXT,
  title VARCHAR(128),
  poetuid VARCHAR(128),
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

	emoji := map[string]string{
		"/::)":  "[微笑]",
		"/::~":  "[撇嘴]",
		"/::B":  "[色]",
		"/::|":  "[发呆]",
		"/:8-)": "[得意]",
		"/::<":  "[流泪]",
		"/::$":  "[害羞]",
		"/::X":  "[闭嘴]",
		"/::Z":  "[睡]",
		"/::’(": "[大哭]",
		"/::-|": "[尴尬]",
		"/::@":  "[发怒]",
		"/::P":  "[调皮]",
		"/::D":  "[呲牙]",
		"/::O":  "[惊讶]",
		"/::(":  "[难过]",
		"/:–b":  "[囧]",
		"/::Q":  "[抓狂]",
		"/::T":  "[吐]",
		"/:,@P": "[偷笑]",
		"/:,@-D" : "[愉快]" ,
		"/::d" : "[白眼]" ,
		"/:,@o" : "[傲慢]" ,
		"/:|-)" : "[困]" ,
		"/::!" : "[惊恐]" ,
		"/::L" : "[流汗]" ,
		"/::>" : "[憨笑]" ,
		"/::,@" : "[悠闲]" ,
		"/:,@f" : "[奋斗]" ,
		"/::-S" : "[咒骂]" ,
		"/:?" : "[疑问]" ,
		"/:,@x" : "[嘘]" ,
		"/:,@@" : "[晕]" ,
		"/:,@!" : "[衰]" ,
		"/:!!!" : "[骷髅]" ,
		"/:xx" : "[敲打]" ,
		"/:bye" : "[再见]" ,
		"/:wipe" : "[擦汗]" ,
		"/:dig" : "[抠鼻]" ,
		"/:handclap" : "[鼓掌]" ,
		"/:B-)" : "[坏笑]" ,
		"/:<@" : "[左哼哼]" ,
		"/:@>" : "[右哼哼]" ,
		"/::-O" : "[哈欠]" ,
		"/:>-|" : "[鄙视]" ,
		"/:P-(" : "[委屈]" ,
		"/::'|" : "[快哭了]" ,
		"/::*" : "[亲亲]" ,
		"/:8*" : "[可怜]" ,
		"/:pd" : "[菜刀]" ,
		"/:<W>" : "[西瓜]" ,
		"/:beer" : "[啤酒]" ,
		"/:coffee" : "[咖啡]" ,
		"/:pig" : "[猪头]" ,
		"/:rose" : "[玫瑰]" ,
		"/:fade" : "[凋谢]" ,
		"/:showlove" : "[嘴唇]" ,
		"/:heart" : "[爱心]" ,
		"/:break" : "[心碎]" ,
		"/:cake" : "[蛋糕]" ,
		"/:bome" : "[炸弹]" ,
		"/:shit" : "[便便]" ,
		"/:moon" : "[月亮]" ,
		"/:sun" : "[太阳]" ,
		"/:hug" : "[拥抱]" ,
		"/:strong" : "[强]" ,
		"/:weak" : "[弱]" ,
		"/:share" : "[握手]" ,
		"/:v" : "[胜利]" ,
		"/:@)" : "[抱拳]" ,
		"/:jj" : "[勾引]" ,
		"/:@@" : "[拳头]" ,
		"/:jump" : "[跳跳]" ,
		"/:shake" : "[发抖]" ,
		"/:<O>" : "[怄火]" ,
		"/:circle" : "[转圈]" ,
		"😄" : "😄" ,
		"😷" : "😷" ,
		"😂" : "😂" ,
		"😝" : "😝" ,
		"😳" : "😳" ,
		"😱" : "😱" ,
		"😔" : "😔" ,
		"😒" : "😒" ,
		"[嘿哈]" : "[嘿哈]" ,
		"[捂脸]" : "[捂脸]" ,
		"[奸笑]" : "[奸笑]" ,
		"[机智]" : "[机智]" ,
		"[皱眉]" : "[皱眉]" ,
		"[耶]" : "[耶]" ,
		"👻" : "👻" ,
		"🙏" : "🙏" ,
		"💪" : "💪" ,
		"🎉" : "🎉" ,
		"🎁" : "🎁" ,
		"[红包]" : "[红包]" ,
		"[鸡]" : "[鸡]" ,
	}

	newkey, exists := emoji[keyword]
	if exists {
		keyword = newkey
	}

	tagId, err := models.InTagMatch(keyword)
	beego.Info("has poetry tag:", tagId)
	if tagId != -1 {
		content, err := models.RandomPoetry(tagId)
		beego.Debug("random poetry is:", content)
		if err != nil {
		    beego.Error(err)
		}
		if len(content) > 0 {
			return content
		}
	}

	rows, err := manager.db.Query("SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
		"FROM poetry " +
		"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
		"WHERE poetry.title = " + "'《" + keyword + "》'  " +
		"ORDER BY poet.poet_count desc")
	if err != nil {
		log.Warn(err)
	}
	defer rows.Close()

	var title string
	var author string
	var content string
	var poetUid string

	for rows.Next() {
		err = rows.Scan(&title, &author, &content, &poetUid)
		if err != nil {
			log.Warn(err)
		}
		fmt.Println(title, author)
		return manager.packFiledsAsString(title, author, content, poetUid)
	}

	err = rows.Err()
	if err != nil {
		log.Warn(err)
	}

	if len(keyword) >= 2 {

		// not found exactly fuzzy match the title.
		query := "SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
			"FROM poetry " +
			"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
			"WHERE poetry.title LIKE " + "'%" + keyword + "%' " +
			"ORDER BY poet.uid ASC, poet.poet_count DESC"
		fmt.Println("query:", query)
		err := manager.db.QueryRow(query).Scan(&title, &author, &content, &poetUid);
		if err != nil {
			log.Warn(err)
		}
		fmt.Println("fuzzy title:", title, author)
		if len(title) > 0 {
			return manager.packFiledsAsString(title, author, content, poetUid)
		}

		// title not found, fuzzy match the poetry content
		query = "SELECT poetry.title, poetry.author, poetry.content, poetry.poetuid " +
			"FROM poetry " +
			"LEFT JOIN poet ON poetry.poetuid = poet.uid " +
			"WHERE poetry.content LIKE " + "'%" + keyword + "%' " +
			"ORDER BY poet.uid ASC, poet.poet_count DESC"
		fmt.Println("query:", query)
		err = manager.db.QueryRow(query).Scan(&title, &author, &content, &poetUid);
		if err != nil {
			log.Warn(err)
		}
		fmt.Println("fuzzy content:", title, content)
		if len(title) > 0 {
			return manager.packFiledsAsString(title, author, content, poetUid)
		}

	} else {
		return "您输入的诗词名称太短哦~"
	}

	return "很抱歉，还没有这首诗哦~"
}

func (manager *DBManager) packFiledsAsString(title, author, content, poetUid string) string {
	if len(content) > 0 {
		content = RenderContent(content, "\n")
	}
	result := title + "\n" + manager.SelectDynasty(strings.Split(poetUid, "_")[0]) + " · " + author + "\n" + content
	if len(result) >= 600 {
		result = result[0:600]
	}
	return result
}

func (manager *DBManager) SelectDynasty(id string) string {
	var name string
	query := "SELECT name " +
		"FROM dynasty " +
		"WHERE id = " + id
	fmt.Println("query:", query)
	err := manager.db.QueryRow(query).Scan(&name);
	if err != nil {
		log.Warn(err)
	}
	return name

}

func RenderContent(content string, returnChar string) string {
	content = strings.Replace(content, "。", "。"+returnChar, -1)
	content = strings.Replace(content, ";", ";"+returnChar, -1)
	content = strings.Replace(content, "；", "；"+returnChar, -1)
	content = strings.Replace(content, "？", "？"+returnChar, -1)
	content = strings.Replace(content, "?", "?"+returnChar, -1)
	content = strings.Replace(content, "!", "!"+returnChar, -1)
	content = strings.Replace(content, "！", "！"+returnChar, -1)
	var match []string
	reg := regexp.MustCompile("（(.*?)）")
	fmt.Println("reg1:", reg.FindAllString(content, -1))
	reg = regexp.MustCompile("\\((.*?)\\)")
	match = reg.FindAllString(content, -1)
	if len(match) > 0 {
		for _, piece := range match {
			content = strings.Replace(content, piece, "", -1)
		}
	}
	fmt.Println("reg2:", match)
	reg = regexp.MustCompile("\\[(.*?)\\]")
	match = reg.FindAllString(content, -1)
	if len(match) > 0 {
		for _, piece := range match {
			content = strings.Replace(content, piece, "", -1)
		}
	}
	fmt.Println("reg3:", match)
	return content
}
