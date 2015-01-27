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
	u := &User{}

	newUser, v := u.Signup(c.Params.Get("login"), c.Params.Get("password"), c.Params.Get("password-confirm"))
	if v.HasErrors() {
		c.RenderArgs["a"] = 1
		c.RenderArgs["validation"] = v
		return c.RenderTemplate("accounts/new.html")
	}

	c.storeUser(&newUser)
	c.Flash.Success("注册成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Login() revel.Result {
	return c.Render()
}

func (c Accounts) LoginCreate() revel.Result {
	return c.Render()
}

func (c Accounts) Logout() revel.Result {
	c.clearUser()
	return c.Redirect(Home.Index)
}
