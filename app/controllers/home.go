package controllers

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Home struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c Home) Index() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}

	return c.Render()
}

func (c Home) Message() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}

	ws := c.Request.Websocket

	Subscribe(c.currentUser.NotifyChannelId(), func(out interface{}) {
		if !ws.IsClientConn() {
			return
		}

		err := websocket.JSON.Send(ws, out)
		if err != nil {
			fmt.Println("WebSocket send error: ", err)
		}
	})
	return nil
}
