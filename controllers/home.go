package controllers

import (
	"github.com/astaxie/beego"
	"gowechatsubscribe/models"
	"gowechatsubscribe/dblite"
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
	c.TplName = "mis.html"
}

func (c *HomeController) Post() {
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	beego.Info("title", title, "content:", content)
	// TODO get data from db

	poetries, err := models.SearchPoetry(title, content)
	if err != nil {
		beego.Error(err)
		c.Redirect("/mis", 302)
		return
	}

	for _, poetry := range poetries {
		poetry.Content = dblite.RenderContent(poetry.Content, "   ")
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
