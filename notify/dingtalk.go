package notify

import (
	"context"

	"github.com/blinkbean/dingtalk"

	"github.com/ggdream/mcc/payload"
)

var _ Notify = (*DingTalk)(nil)

type DingTalk struct {
	client *dingtalk.DingTalk
}

func NewDingTalk(token, secret string) *DingTalk {
	return &DingTalk{
		client: dingtalk.InitDingTalkWithSecret(token, secret),
	}
}

// SendPushMessage implements Notify.
func (d *DingTalk) SendPushMessage(_ context.Context, payload *payload.PushPayload) error {
	return d.client.SendActionCardMessage("代码部署: "+payload.Repo.FullName, payload.HeadCommit.Message, dingtalk.WithCardSingleURL(payload.HeadCommit.URL), dingtalk.WithCardSingleTitle("代码部署: "+payload.Repo.FullName))
}

// SendTagMessage implements Notify.
func (d *DingTalk) SendTagMessage(ctx context.Context, payload *payload.PushPayload) error {
	return d.client.SendActionCardMessage("发版提醒: "+payload.Repo.FullName, payload.HeadCommit.Message, dingtalk.WithCardSingleURL(payload.HeadCommit.URL), dingtalk.WithCardSingleTitle("发版提醒: "+payload.Repo.FullName))
}
