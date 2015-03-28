package controllers

import (
	"bytes"
	"fmt"
	"github.com/acsellers/inflections"
	"github.com/dchest/captcha"
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
	"reflect"
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
	c.RenderArgs["controller_name"] = inflections.Underscore(c.Name)
	c.RenderArgs["method_name"] = inflections.Underscore(c.MethodName)
	c.RenderArgs["route_name"] = fmt.Sprintf("%v#%v", inflections.Underscore(c.Name), inflections.Underscore(c.MethodName))
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
	userId := c.Session["user_id"]
	u := &User{}
	c.currentUser = u
	if len(userId) == 0 {
		return
	}

	DB.Where("id = ?", c.Session["user_id"]).First(u)
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
		revel.INFO.Println("current_user { id: ", c.currentUser.Id, ", login: ", c.currentUser.Login, " }")
		return nil
	}
}

func (c App) requireUserForJSON() revel.Result {
	if r := c.requireUser(); r != nil {
		return c.errorJSON(-1, "还未登录")
	} else {
		return r
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
	if c.currentUser.IsAdmin() {
		return true
	}

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
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c App) errorJSON(code int, msg string) revel.Result {
	result := AppResult{Code: code, Msg: msg}
	return c.RenderJson(result)
}

func (c App) errorsJSON(code int, errs []*revel.ValidationError) revel.Result {
	msgs := make([]string, len(errs))
	for i, err := range errs {
		msgs[i] = err.Message
	}
	result := AppResult{Code: code, Msg: strings.Join(msgs, "\n")}
	return c.RenderJson(result)
}

func (c App) successJSON(data interface{}) revel.Result {
	result := AppResult{Code: 0, Data: data}
	return c.RenderJson(result)
}

func (c App) Captcha(id string) revel.Result {
	captchaId := captcha.NewLen(4)
	c.Session["captcha_id"] = captchaId

	var buffer bytes.Buffer
	captcha.WriteImage(&buffer, captchaId, 200, 80)

	c.Response.ContentType = "image/png"
	c.Response.Status = 200

	return c.RenderText(buffer.String())
}

func (c App) validateCaptcha(code string) bool {
	return captcha.VerifyString(c.Session["captcha_id"], code)
}
