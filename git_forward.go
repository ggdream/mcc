package main

import "github.com/ggdream/mcc/payload"

type GitForward interface {
	Name() string
	Source() Source
	Event() Event
	GitEvent
}

type GitEvent interface {
	GetPushPayload(data []byte) (*payload.PushPayload, error)
}
