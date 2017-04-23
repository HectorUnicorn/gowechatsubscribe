package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

const (
	mysqlDriver = "mysql"
	mysqlConn   = "root:guojialin@/poetry?charset=utf8"
)

func RegisterDB() {
	orm.RegisterDataBase("default", mysqlDriver, mysqlConn, 30)
	orm.RegisterModel(new(Poetry))

	// 自动建表
	orm.RunSyncdb("default", false, true)
}

type Poetry struct {
	Id        int64 `orm:"auto;index"`
	Url       string
	Content   string `orm:"type(text)"`
	Author    string `orm:"size(255)"`
	Interpret string `orm:"type(text)"`
	Title     string `orm:"size(255)"`
	Poetuid   string `orm:"size(255)"`
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
