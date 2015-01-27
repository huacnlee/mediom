package controllers

import (
	"fmt"
	"github.com/revel/revel"
)

type Users struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Users).Before, revel.BEFORE)
//	revel.InterceptMethod((*Users).After, revel.AFTER)
//}

func (c Users) Show(username string) revel.Result {
	return c.RenderText(fmt.Sprintf("You want visit %s's home page.", username))
}
