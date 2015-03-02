package controllers

import (
	"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
	"reflect"
	"strconv"
	"strings"
)

type App struct {
	*revel.Controller
	currentUser *User
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

func (c App) requireAdmin() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}

	if !c.currentUser.IsAdmin() {
		c.Flash.Error("此功能需要管理员权限。")
		return c.Redirect("/")
	}

	return nil
}

func (c App) isOwner(obj interface{}) bool {
	objType := reflect.TypeOf(obj)
	switch objType.String() {
	case "models.Topic":
		return c.currentUser.Id == obj.(Topic).UserId
	case "*models.Topic":
		return c.currentUser.Id == obj.(*Topic).UserId
	case "models.User":
		return c.currentUser.Id == obj.(User).Id
	case "*models.User":
		return c.currentUser.Id == obj.(*User).Id
	case "models.Reply":
		return c.currentUser.Id == obj.(Reply).UserId
	case "*models.Reply":
		return c.currentUser.Id == obj.(*Reply).UserId
	default:
		panic(fmt.Sprintf("Invalid isOwner type: %v, %v, name: %v", obj, objType, objType.Name()))
	}

	return false
}

func (c App) renderValidation(tpl string, v revel.Validation) revel.Result {
	c.RenderArgs["validation"] = v
	return c.RenderTemplate(tpl)
}

type AppResult struct {
	code int
	msg  string
	data interface{}
}

func (c App) errorJSON(code int, msg string) revel.Result {
	result := &AppResult{code: code, msg: msg}
	return c.RenderJson(result)
}

func (c App) errorsJSON(code int, errs []*revel.ValidationError) revel.Result {
	msgs := make([]string, len(errs))
	for i, err := range errs {
		msgs[i] = err.Message
	}
	result := &AppResult{code: code, msg: strings.Join(msgs, "\n")}
	return c.RenderJson(result)
}

func (c App) successJSON(data interface{}) revel.Result {
	result := &AppResult{code: 0, msg: "", data: data}
	return c.RenderJson(result)
}
