package controllers

import (
	"github.com/astaxie/beego"
	"gowechatsubscribe/models"
	"gowechatsubscribe/dblite"
	"strconv"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["IsHome"] = true
	login := checkAccount(c.Ctx)
	c.Data["IsLogin"] = login
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}

	tags, err := models.GetAllTags()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Tags"] = tags

	op := c.Input().Get("op")
	switch op {
	case "showtag":
		id := c.Input().Get("tag")
		tagName := c.Input().Get("tagname")
		beego.Debug("show poetry tag id:", id)
		if len(id) == 0 {
			break
		}
		tid, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			beego.Error(err)
		}
		poetries, err := models.GetPoetriesByTagId(int(tid))
		if err != nil {
			beego.Error(err)
		}
		for _, poetry := range poetries {
			poetry.Content = dblite.RenderContent(poetry.Content, "<br/>")
		}
		c.Data["Poetries"] = poetries
		c.Data["TagName"] = tagName
	}

	c.TplName = "mis.html"
}

func (c *HomeController) Post() {
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	beego.Info("title", title, "content:", content)
	// TODO get data from db

	tags, err := models.GetAllTags()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Tags"] = tags

	poetries, err := models.SearchPoetry(title, content)
	if err != nil {
		beego.Error(err)
		c.Redirect("/mis", 302)
		return
	}

	for _, poetry := range poetries {
		poetry.Content = dblite.RenderContent(poetry.Content, "<br/>")
	}

	c.Data["Poetries"] = poetries
	c.Data["SearchTitle"] = title
	c.Data["SearchContent"] = content

	beego.Debug("title", title, "content", content)
	c.Data["IsHome"] = true
	login := checkAccount(c.Ctx)
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	c.TplName = "mis.html"
}
