package controllers

import (
	"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
	"reflect"
	"strconv"
)

type App struct {
	*revel.Controller
	currentUser *User
}

func (c *App) prependCurrentUser() {
	userId, _ := strconv.Atoi(c.Session["user_id"])
	u := &User{}
	DB.Where("id = ?", userId).First(u)
	c.currentUser = u
}

func (c App) storeUser(u *User) {
	if u.Id == 0 {
		return
	}
	c.Session["user_id"] = fmt.Sprintf("%v", u.Id)
}

func (c App) clearUser() {
	c.Session["user_id"] = ""
}

func (c App) isLogined() bool {
	return c.currentUser.Id > 0
}

func (c App) requireUser() revel.Result {
	if !c.isLogined() {
		c.Flash.Error("你还未登录哦")
		return c.Redirect(Accounts.Login)
	} else {
		fmt.Println("current_user: ", c.currentUser)
		return nil
	}
}

func (c App) isOwner(obj interface{}) bool {
	objType := reflect.TypeOf(obj)
	switch objType.Name() {
	case "models.Topic":
		return c.currentUser.Id == obj.(Topic).UserId
	case "models.User":
		return c.currentUser.Id == obj.(User).Id
	case "models.Reply":
		return c.currentUser.Id == obj.(Reply).UserId
	}
	return false
}

func (c App) renderValidation(tpl string, v revel.Validation) revel.Result {
	c.RenderArgs["validation"] = v
	return c.RenderTemplate(tpl)
}

func init() {
	revel.InterceptMethod((*App).Before, revel.BEFORE)
	revel.InterceptMethod((*App).After, revel.AFTER)
}

func (c *App) Before() revel.Result {
	c.prependCurrentUser()
	c.RenderArgs["validation"] = nil
	c.RenderArgs["logined"] = c.isLogined()
	c.RenderArgs["current_user"] = c.currentUser
	return c.Result
}

func (c *App) After() revel.Result {
	newParams := make(map[string]string, len(c.Params.Values))
	for key := range c.Params.Values {
		newParams[key] = c.Params.Get(key)
	}
	if len(newParams) > 0 {
		c.RenderArgs["params"] = newParams
	}
	return c.Result
}
