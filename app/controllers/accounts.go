package controllers

import (
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Accounts struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Accounts).Before, revel.BEFORE)
//	revel.InterceptMethod((*Accounts).After, revel.AFTER)
//}

func (c Accounts) New() revel.Result {
	a := "hello world"
	b := "foobar"
	return c.Render(a, b)
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
	c.Flash.Success("注册成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Logout() revel.Result {
	c.clearUser()
	return c.Redirect(Home.Index)
}

func (c Accounts) Edit() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	return c.Render("accounts/edit.html")
}

func (c Accounts) Update() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
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
		return c.Render("accounts/edit.html")
	}
	c.Flash.Success("个人信息修改成功")
	return c.Redirect("/account/edit")
}
