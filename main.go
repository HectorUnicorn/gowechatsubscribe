package main

import (
	_ "gowechatsubscribe/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"

)

func main() {
	beego.Any("/", hello)
	beego.Run()
}

func hello(ctx *context.Context) {
	//配置微信参数
	config := &wechat.Config{
		AppID: beego.AppConfig.String("AppID"),
		AppSecret: beego.AppConfig.String("AppSecret"),
		Token: beego.AppConfig.String("Token"),
		EncodingAESKey: beego.AppConfig.String("EncodingAESKey"),
	}
	wc := wechat.NewWechat(config)

	// 传入request和responseWriter
	server := wc.GetServer(ctx.Request, ctx.ResponseWriter)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		switch msg.MsgType {
		// 文本消息
		case message.MsgTypeText:
			//回复消息：演示回复用户发送的消息
			//text := message.NewText(msg.Content)
			return &message.Reply{message.MsgTypeText, "王照文你好"}
		}

	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

