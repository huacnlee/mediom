package controllers

import (
	"fmt"
	_ "fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Accounts struct {
	*revel.Controller
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
