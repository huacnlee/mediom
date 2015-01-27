package controllers

import (
	"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
	"strconv"
)

type App struct {
	*revel.Controller
}

func (c App) currentUser() *User {
	var u User
	userId, _ := strconv.Atoi(c.Session["user_id"])
	DB.Where("id = ?", userId).First(&u)
	return &u
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
