package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"math/rand"
)

const (
	mysqlDriver = "mysql"
)

type Poetry struct {
	Id        int64 `orm:"auto;index"`
	Url       string
	Content   string `orm:"type(text)"`
	Author    string `orm:"size(128)"`
	Interpret string `orm:"type(text)"`
	Title     string `orm:"size(128)"`
	Poetuid   string `orm:"size(128)"`
}

type Tag struct {
	Id          int `orm:"auto;index"`
	Tag         string `orm:"unique;size(128)"`
	TagCategory string `orm:"size(128)"`
}

type PoetryTag struct {
	Id        int `orm:"auto;index"`
	TagId     int
	PoetryId  int64
	BestLines string `orm:"size(128)"`
}

type TagState struct {
	Tag         Tag
	PoetryId    int64
	PoetryTagId int
	BestLines   string
	Active      bool
}

func RegisterDB() {
	mysqlConn := fmt.Sprintf("%s:%s@%s", beego.AppConfig.String("dbusername"),
		beego.AppConfig.String("dbpassword"), beego.AppConfig.String("dbdatabase"))
	beego.Info("mysql conn:", mysqlConn)
	orm.RegisterDataBase("default", mysqlDriver, mysqlConn, 30)
	orm.RegisterModel(new(Poetry), new(Tag), new(PoetryTag))

	// 自动建表
	orm.RunSyncdb("default", false, true)
}

func SearchPoetry(title, content string) ([]*Poetry, error) {
	o := orm.NewOrm()
	poetries := make([]*Poetry, 0)
	qs := o.QueryTable("poetry")
	var err error
	_, err = qs.Filter("title__contains", title).Filter("content__contains", content).Limit(20).All(&poetries)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return poetries, err
}

func GetPoetry(id int64) (*Poetry, error) {
	o := orm.NewOrm()
	poetry := Poetry{Id: id}
	var err error
	err = o.Read(&poetry)
	if err != nil {
		beego.Error(err)
	}
	return &poetry, err
}

func GetAllTags() ([]*Tag, error) {
	o := orm.NewOrm()
	tags := make([]*Tag, 0)
	qs := o.QueryTable("tag")
	_, err := qs.All(&tags)
	if err != nil {
		beego.Error(err)
	}
	return tags, err
}
func GetPoetriesByTagId(tagId int) ([]*Poetry, error) {
	o := orm.NewOrm()
	poetries := make([]*Poetry, 0)
	poetryTags := make([]*PoetryTag, 0)
	_, err := o.QueryTable("poetry_tag").Filter("tag_id", tagId).All(&poetryTags)

	ids := make([]int64, 0)
	for _, poetryTag := range poetryTags {
		ids = append(ids, poetryTag.PoetryId)
	}

	beego.Info("poetry ids:", ids)

	if len(ids) != 0 {
		_, err = o.QueryTable("poetry").Filter("id__in", ids).All(&poetries)
		if err != nil {
			beego.Error(err)
		}
	}
	return poetries, err
}

func AddTag(tag, tagcate string) (int64, error) {
	t := Tag{Tag: tag, TagCategory: tagcate}
	o := orm.NewOrm()
	id, err := o.Insert(&t)
	if err != nil {
		beego.Error(err)
	}
	return id, err
}

func DelTag(id int) error {
	t := Tag{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&t)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func SetPoetryTag(poetryId int64, tagId int, bestLines string) error {
	poetryTag := PoetryTag{PoetryId: poetryId, TagId: tagId, BestLines: bestLines}
	o := orm.NewOrm()
	_, err := o.Insert(&poetryTag)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func DelPoetryTag(poetryTagId int) error {
	poetryTag := PoetryTag{Id: poetryTagId}
	o := orm.NewOrm()
	_, err := o.Delete(&poetryTag)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func DelTagsOfPoetry(tagId int) error {
	poetryTag := PoetryTag{TagId: tagId}
	o := orm.NewOrm()
	_, err := o.Delete(&poetryTag)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func GetPoetryTagState(id int64) ([]*TagState, error) {
	poetryTags := make([]*PoetryTag, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("poetry_tag")
	_, err := qs.Filter("poetry_id", id).All(&poetryTags)
	if err != nil {
		beego.Error(err)
	}
	tags := make([]*Tag, 0)
	_, err = o.QueryTable("tag").All(&tags)
	if err != nil {
		beego.Error(err)
	}
	tagsStates := make([]*TagState, 0)
	beego.Info("tags", tags, " tags2", poetryTags)
	for _, t := range tags {
		tagSt := new(TagState)
		tagSt.Tag = *t
		tagSt.Active = false
		tagSt.PoetryId = id
		for _, t2 := range poetryTags {
			if t.Id == t2.TagId {
				tagSt.Active = true
				tagSt.BestLines = t2.BestLines
				tagSt.PoetryTagId = t2.Id
				break
			}
		}
		tagsStates = append(tagsStates, tagSt)
	}
	return tagsStates, err
}

func InTagMatch(keyword string) (int, error) {
	o := orm.NewOrm()
	tags := make([]*Tag, 0)
	_, err := o.QueryTable("tag").Filter("tag", keyword).All(&tags)

	if len(tags) > 0 {
		return tags[0].Id, nil
	}
	return -1, err
}

func RandomPoetry(tagId int) (string, error) {
	o := orm.NewOrm()
	poetryTags := make([]*PoetryTag, 0)
	qs := o.QueryTable("poetry_tag")
	qs.Filter("tag_id", tagId).All(&poetryTags)
	if len(poetryTags) > 0 {
		rand.Seed(time.Now().UnixNano())
		t := rand.Intn(len(poetryTags))
		return poetryTags[t-1].BestLines, nil
	}
	return "", nil
}
