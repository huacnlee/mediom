package controllers

import (
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
)

type Accounts struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Accounts).Before, revel.BEFORE)
//	revel.InterceptMethod((*Accounts).After, revel.AFTER)
//}

func (c Accounts) New() revel.Result {
	return c.Render()
}

func (c Accounts) Create() revel.Result {
	u := User{}
	newUser := User{}

	v := revel.Validation{}

	if !c.validateCaptcha(c.Params.Get("captcha")) {
		v.Error("验证码不正确")
		return c.renderValidation("accounts/new.html", v)
	}

	newUser, v = u.Signup(c.Params.Get("login"), c.Params.Get("password"), c.Params.Get("password-confirm"))
	if v.HasErrors() {
		return c.renderValidation("accounts/new.html", v)
	}

	c.storeUser(&newUser)
	c.Flash.Success("注册成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Login() revel.Result {
	return c.Render()
}

func (c Accounts) LoginCreate() revel.Result {
	u := User{}
	newUser := User{}
	v := revel.Validation{}

	if !c.validateCaptcha(c.Params.Get("captcha")) {
		v.Error("验证码不正确")
		return c.renderValidation("accounts/login.html", v)
	}

	newUser, v = u.Signin(c.Params.Get("login"), c.Params.Get("password"))
	if v.HasErrors() {
		return c.renderValidation("accounts/login.html", v)
	}

	c.storeUser(&newUser)
	c.Flash.Success("登录成功，欢迎再次回来。")
	return c.Redirect(Home.Index)
}

func (c Accounts) Logout() revel.Result {
	c.clearUser()
	c.Flash.Success("登出成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Edit() revel.Result {
	c.requireUser()
	return c.Render()
}

func (c Accounts) Update() revel.Result {
	c.requireUser()
	c.currentUser.Email = c.Params.Get("email")
	c.currentUser.GitHub = c.Params.Get("github")
	c.currentUser.Twitter = c.Params.Get("twitter")
	c.currentUser.Tagline = c.Params.Get("tagline")
	c.currentUser.Location = c.Params.Get("location")
	c.currentUser.Description = c.Params.Get("description")
	var u User
	u = *c.currentUser
	_, v := UpdateUserProfile(u)
	if v.HasErrors() {
		return c.renderValidation("accounts/edit.html", v)
	}
	c.Flash.Success("个人信息修改成功")
	return c.Redirect("/account/edit")
}
