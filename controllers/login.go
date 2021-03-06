package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var success = 0

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
	c.Data["IsSuccess"] = success == 1
	c.Data["IsFailed"] = success == 2
}

func (c *LoginController) Post() {
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLogin") == "on"
	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<32 - 1
		}
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")

		c.Data["IsSuccess"] = true
		c.Data["IsFailed"] = false
		success = 0
		c.Redirect("/mis", 301)
		beego.Info("Login Successful! suc: ", success)
		return
	} else {
		c.Data["IsSuccess"] = false
		c.Data["IsFailed"] = true
		success = 2
		beego.Info("Login Failed! suc: ", success)
	}
	c.Redirect("/mis/login", 301)
	return
}

func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}
	pwd := ck.Value
	beego.Info("Get from Cookie: uname - pwd: ", uname, pwd)
	return beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd
}
