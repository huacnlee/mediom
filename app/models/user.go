package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/revel/revel"
	"strings"
	"time"
)

type User struct {
	Id        int32
	Login     string `sql:"size:255;not null"`
	Password  string `sql:"size:255;not null"`
	Topics    []Topic
	Replies   []Reply
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) EncodePassword(raw string) (md5Digest string) {
	data := []byte(raw)
	result := md5.Sum(data)
	md5Digest = hex.EncodeToString(result[:])
	return
}

func (u User) Signup(login string, password string, passwordConfirm string) (user User, v revel.Validation) {
	u.Login = strings.Trim(login, " ")

	v.MinSize(login, 5).Key("用户名").Message("最少要 5 个字符")
	v.MinSize(password, 6).Key("密码").Message("最少要 6 个字符")

	if password != passwordConfirm {
		v.Error("密码与确认密码不一致")
	}

	var existCount int
	db.Model(&User{}).Where("login = ?", login).Count(&existCount)
	fmt.Println("login same as: ", login, " have ", existCount)
	if existCount > 0 {
		v.Error("帐号已经被注册")
	}

	if v.HasErrors() {
		return u, v
	}

	u.Password = u.EncodePassword(password)

	err := db.Create(&u).Error
	if err != nil {
		v.Error(fmt.Sprintf("服务器异常, %v", err))
	}
	fmt.Println("created user: ", u)
	return u, v
}

func (u User) Signin(login string, password string) (user User, v revel.Validation) {
	login = strings.Trim(login, " ")

	if len(password) == 0 {
		v.Error("还未输入密码")
	}

	db.First(&user, "login = ? and password = ?", login, u.EncodePassword(password))
	fmt.Println("first user:", user)
	if user.Id == 0 {
		v.Error("帐号密码不正确")
	}
	return user, v
}
