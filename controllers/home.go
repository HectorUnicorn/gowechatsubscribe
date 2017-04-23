package controllers

import "github.com/astaxie/beego"

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["IsHome"] = true
	login := checkAccount(c.Ctx)
	if !login {
		c.Redirect("/mis/login", 302)
		return
	}
	c.TplName = "mis.html"
}
