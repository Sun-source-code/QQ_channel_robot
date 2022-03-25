package main

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
)

func (p Processor) sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := p.api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println("回复信息", err)
	}
}
