package controllers

import (
	"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
	"strconv"
)

type App struct {
	*revel.Controller
	CurrentUser User
}

func (c App) currentUser() *User {
	fmt.Println("is nil", c.CurrentUser.Id == 0)
	if c.CurrentUser.Id > 0 {
		return &c.CurrentUser
	}
	userId, _ := strconv.Atoi(c.Session["user_id"])
	DB.Where("id = ?", userId).First(&c.CurrentUser)
	return &c.CurrentUser
}

func (c App) storeUser(u *User) {
	fmt.Println("storeUser ", u)
	if u.Id == 0 {
		return
	}
	fmt.Println("will store session", u.Id)
	c.Session["user_id"] = fmt.Sprintf("%v", u.Id)
}

func (c App) clearUser() {
	c.Session["user_id"] = ""
}

func (c App) requireUser() revel.Result {
	u := c.currentUser()
	if u.Id == 0 {
		c.Flash.Error("你还未登录哦")
		return c.Redirect(Accounts.Login)
	} else {
		fmt.Println("current_user: ", u)
		return nil
	}
}

func init() {
	revel.InterceptMethod((*App).Before, revel.BEFORE)
	revel.InterceptMethod((*App).After, revel.AFTER)
}

func (c App) Before() revel.Result {
	u := c.currentUser()
	c.RenderArgs["validation"] = nil
	if u.Id > 0 {
		c.RenderArgs["current_user"] = u
	}
	return c.Result
}

func (c App) After() revel.Result {
	newParams := make(map[string]string, len(c.Params.Values))
	for key := range c.Params.Values {
		newParams[key] = c.Params.Get(key)
	}
	c.RenderArgs["params"] = newParams
	return c.Result
}
