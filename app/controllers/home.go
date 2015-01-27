package controllers

import (
	"github.com/revel/revel"
)

type Home struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c Home) Index() revel.Result {
	c.requireUser()
	return c.Render()
}
