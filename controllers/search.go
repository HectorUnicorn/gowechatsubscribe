package controllers

import (
	"github.com/astaxie/beego"
	"gowechatsubscribe/models"
)

type SearchController struct {
	beego.Controller
}

func (c *SearchController) Post() {
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	// TODO get data from db

	poetries, err := models.SearchPoetry(title, content)
	if err != nil {
		beego.Error(err)
		c.Redirect("/mis", 302)
		return
	}

	c.Data["Poetries"] = poetries

	beego.Debug("title", title, "content", content)
	c.Data["IsHome"] = true
	login := checkAccount(c.Ctx)
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	c.TplName = "mis.html"
}
