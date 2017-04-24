package controllers

import (
	"github.com/astaxie/beego"
	"gowechatsubscribe/models"
	"strconv"
)

type TagController struct {
	beego.Controller
}

func (c *TagController) Get() {
	login := checkAccount(c.Ctx)
	c.Data["IsLogin"] = login
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	op := c.Input().Get("op")
	switch op {
	case "del":
		id := c.Input().Get("id")
		beego.Debug("delete tag", id)
		if len(id) == 0 {
			break
		}
		tid, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
		    beego.Error(err)
		}
		err = models.DelTag(int(tid))
		if err != nil {
			beego.Error(err)
		}
		err = models.DelTagsOfPoetry(int(tid))
		if err != nil {
		    beego.Error(err)
		}
		c.Redirect("/mis/tag", 301)
		return
	}

	tags, err := models.GetAllTags()
	if err != nil {
	    beego.Error(err)
	}
	c.Data["Tags"] = tags
	c.TplName = "tags.html"
}


func (c *TagController) Post() {
	tagName := c.Input().Get("tag")
	tagCate := c.Input().Get("tag_category")
	beego.Info("tag", tagName, tagCate)
	id, err := models.AddTag(tagName, tagCate)
	if err != nil {
	    beego.Error(err)
	}
	beego.Info("rowId:", id)
	c.Redirect("/mis/tag", 302)
}

func (c *TagController) Delete() {

}