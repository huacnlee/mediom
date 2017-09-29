package controllers

import (
	"fmt"
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

// Home controller
type Home struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

// Index - GET /
func (c Home) Index() revel.Result {
	return c.Render()
}

// Message - WS /msg
func (c Home) Message() revel.Result {
	if !c.isLogined() {
		return nil
	}

	ws := c.Request.Websocket

	Subscribe(c.currentUser.NotifyChannelId(), func(out interface{}) {
		err := websocket.JSON.Send(ws, out)
		if err != nil {
			fmt.Println("WebSocket send error: ", err)
		}
	})
	return nil
}

// Search GET /search
func (c Home) Search() revel.Result {
	return c.Redirect(fmt.Sprintf("https://google.com?q=site:ruby-china.org %v", c.Params.Get("q")))
}
