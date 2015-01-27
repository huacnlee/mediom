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

	v := u.Signup(c.Params.Get("login"), c.Params.Get("password"), c.Params.Get("password-confirm"))
	if v.HasErrors() {
		c.RenderArgs["errors"] = v.Errors
		return c.RenderTemplate("accounts/new.html")
	}

	//v = u.

	//if v.HasErrors() {
	//// Store the validation errors in the flash context and redirect.
	//c.Validation.Keep()
	//fmt.Println(c.Validation.ErrorMap())
	//c.FlashParams()
	//return c.RenderTemplate("accounts/new.html")
	//}

	return c.RenderTemplate("accounts/new.html")
}

func (c Accounts) Login() revel.Result {
	return c.Render()
}

func (c Accounts) LoginCreate() revel.Result {
	return c.Render()
}
