package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Processor struct {
	api openapi.OpenAPI
}

var pre string

func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData, idiom map[string]string) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	firstWord := getRandKey(idiom)

	//// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}

	switch cmd.Cmd {
	case "hi":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "成语接龙开始":
		toStart := &dto.MessageToCreate{
			Content: firstWord,
			MessageReference: &dto.MessageReference{
				MessageID:             data.ID,
				IgnoreGetMessageError: false,
			},
		}
		p.sendReply(ctx, data.ChannelID, toStart)
		pre = firstWord

	default:
		if idiomExists(cmd.Cmd, idiom) && idiomIsRight(cmd.Cmd, pre) {
			nextWord := idiomSelect(cmd.Cmd, idiom)
			toNext := &dto.MessageToCreate{
				Content: nextWord + " (可回复查看释义)",
				MessageReference: &dto.MessageReference{
					MessageID:             data.ID,
					IgnoreGetMessageError: false,
				},
			}
			p.sendReply(ctx, data.ChannelID, toNext)
			pre = nextWord
		} else if cmd.Cmd == "查看释义" {
			value := getMeaning(pre, idiom)
			toNext := &dto.MessageToCreate{
				Content: value,
				MessageReference: &dto.MessageReference{
					MessageID:             data.ID,
					IgnoreGetMessageError: false,
				},
			}
			p.sendReply(ctx, data.ChannelID, toNext)
		} else if !idiomExists(cmd.Cmd, idiom) {
			toNext := &dto.MessageToCreate{
				Content: "你输入的成语不存在",
				MessageReference: &dto.MessageReference{
					MessageID:             data.ID,
					IgnoreGetMessageError: false,
				},
			}
			p.sendReply(ctx, data.ChannelID, toNext)
		} else if idiomIsRight(cmd.Cmd, pre) == false {
			toNext := &dto.MessageToCreate{
				Content: "你输入的成语不合法",
				MessageReference: &dto.MessageReference{
					MessageID:             data.ID,
					IgnoreGetMessageError: false,
				},
			}
			p.sendReply(ctx, data.ChannelID, toNext)
		}
	}

	return nil
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s

消息来自：%s
`
	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		getIP(),
	)
}

func (p Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func getIdiom() map[string]string {
	Map := make(map[string]string)
	f, err := os.Open("/Users/mac/Desktop/GolangProject/Idiom/data.txt")
	if err != nil {
		fmt.Println("err", err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		words := strings.Split(line, "\t")
		Map[words[0]] = words[2]
	}
	fmt.Println("map的长度为 ", len(Map))
	return Map
}
