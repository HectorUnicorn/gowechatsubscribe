package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

const (
	mysqlDriver = "mysql"
	mysqlConn   = "root:guojialin@/poetry?charset=utf8"
)

type Poetry struct {
	Id        int64 `orm:"auto;index"`
	Url       string
	Content   string `orm:"type(text)"`
	Author    string `orm:"size(255)"`
	Interpret string `orm:"type(text)"`
	Title     string `orm:"size(255)"`
	Poetuid   string `orm:"size(255)"`
}

type Tag struct {
	Id          int `orm:"auto;index"`
	Tag         string `orm:"unique;size(255)"`
	TagCategory string `orm:"size(255)"`
}

type PoetryTag struct {
	Id        int `orm:"auto;index"`
	TagId     int
	PoetryId  int64
	BestLines string `orm:"size(255)"`
}

type TagState struct {
	Tag         Tag
	PoetryId    int64
	PoetryTagId int
	BestLines   string
	Active      bool
}

func RegisterDB() {
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