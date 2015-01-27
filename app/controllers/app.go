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

func (c App) requireUser() {
	u := c.currentUser()
	if u.Id == 0 {
		c.Flash.Error("你还未登录哦")
		c.Redirect("/signin")
	} else {
		fmt.Println("current_user: ", u)
	}
}
