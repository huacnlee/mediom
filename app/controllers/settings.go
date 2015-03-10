package controllers

import (
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Settings struct {
	App
}

func init() {
	revel.InterceptMethod((*Settings).Before, revel.BEFORE)
}

func (c *Settings) Before() revel.Result {
	if r := c.requireAdmin(); r != nil {
		return r
	}

	return nil
}

func (c Settings) Index() revel.Result {
	settings := []Setting{}
	DB.Model(Setting{}).Order("`key` desc").Find(&settings)
	c.RenderArgs["settings"] = settings
	return c.Render("settings/index.html")
}

func (c Settings) Edit(key string) revel.Result {
	setting := FindSettingByKey(key)
	c.RenderArgs["setting"] = setting
	return c.Render("settings/edit.html")
}

func (c Settings) Update(key string) revel.Result {
	setting := FindSettingByKey(key)
	c.Params.Bind(&setting.Val, "val")
	c.RenderArgs["setting"] = setting
	if err := DB.Save(&setting).Error; err != nil {
		return c.Render("settings/edit.html")
	}
	c.Flash.Success("设置更新成功")
	return c.Redirect("/settings")
}
