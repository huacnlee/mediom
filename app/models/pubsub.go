package models

import (
	"github.com/tuxychandru/pubsub"
)

var (
	PushCenter *pubsub.PubSub
)

func initPubsub() {
	PushCenter = pubsub.New(100)
}

func PushMessage(channel string, message interface{}) {
	PushCenter.Pub(message, channel)
}

func Subscribe(channel string, callback func(message interface{})) {
	ch := PushCenter.Sub(channel)
	for {
		out := <-ch
		callback(out)
	}
}
