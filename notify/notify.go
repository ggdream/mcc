package notify

import (
	"context"
	"slices"

	"github.com/ggdream/mcc/payload"
)

var ins Notify
var scenes []string

type Notify interface {
	SendPushMessage(ctx context.Context, payload *payload.PushPayload) error
	SendTagMessage(ctx context.Context, payload *payload.PushPayload) error
}

func Init(scenes_ []string, token, secret string) (err error) {
	scenes = scenes_
	ins = NewDingTalk(token, secret)
	return nil
}

func SendPushMessage(ctx context.Context, payload *payload.PushPayload) error {
	if ins == nil {
		return nil
	}

	if !slices.Contains(scenes, "push") {
		return nil
	}

	return ins.SendPushMessage(ctx, payload)
}
