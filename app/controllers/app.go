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
}

func (c App) Before() revel.Result {
	u := c.currentUser()
	if u.Id > 0 {
		c.RenderArgs["current_user"] = u
	}
	return c.Result
}
