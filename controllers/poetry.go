package controllers

import (
	"github.com/astaxie/beego"
	"gowechatsubscribe/models"
	"strconv"
	"gowechatsubscribe/dblite"
)

type PoetryController struct {
	beego.Controller
}

func (c *PoetryController) Get() {
	login := checkAccount(c.Ctx)
	c.Data["IsLogin"] = login
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	op := c.Input().Get("op")
	switch op {
	case "deltag":
		id := c.Input().Get("id")
		poetryId := c.Input().Get("poetry_id")
		beego.Debug("delete poetry tag", id)
		if len(id) == 0 {
			break
		}
		tid, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			beego.Error(err)
		}
		err = models.DelPoetryTag(int(tid))
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/mis/poetry/" + poetryId, 301)
		return
	}

	id := c.Ctx.Input.Param(":id")
	beego.Info("poetryId:", id)
	c.Data["PoetryId"] = id
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	poetry, err := models.GetPoetry(pid)
	if err != nil {
		beego.Error(err)
	}
	poetryTags, err := models.GetPoetryTagState(pid)
	if err != nil {
		beego.Error(err)
	}

	if poetry != nil {
		poetry.Content = dblite.RenderContent(poetry.Content, "<br/>")
	}

	c.Data["PoetryTags"] = poetryTags
	c.Data["Poetry"] = poetry
	c.TplName = "poetry_view.html"
}

func (c *PoetryController) Post() {
	login := checkAccount(c.Ctx)
	c.Data["IsLogin"] = login
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	// for add tag to poetry
	poetryId := c.Input().Get("poetry_id")
	tagId := c.Input().Get("tag_id")
	bestLines := c.Input().Get("best_lines")
	if len(poetryId) == 0 || len(tagId) == 0 {
		return
	}
	pid, err := strconv.ParseInt(poetryId, 10, 64)
	tid, err := strconv.ParseInt(tagId, 10, 32)
	err = models.SetPoetryTag(pid, int(tid), bestLines)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/mis/poetry/"+poetryId, 301)
}
