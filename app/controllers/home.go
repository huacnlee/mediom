package controllers

import (
	"github.com/revel/revel"
)

type Home struct {
	*revel.Controller
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c Home) Index() revel.Result {
	return c.Render()
}
