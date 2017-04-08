package main

import (
	_ "gowechatsubscribe/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
	"gowechatsubscribe/dblite"
	"github.com/going/toolkit/log"
)

func main() {
	beego.Any("/", hello)
	beego.Run()
}

var accessToken string

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

		//openId := server.GetOpenID()

		switch msg.MsgType {
		// 文本消息
		case message.MsgTypeText:
			//回复消息：演示回复用户发送的消息
			text := message.NewText(msg.Content)
			dbManager := dblite.NewDBManager()
			var result string
			fmt.Println("input:", text.Content)
			if (len(text.Content) < 2) {
				result = "请输入两个字以上哦！"
			} else {
				result = dbManager.SelectPoetry(text.Content)
			}
			reply := message.NewText(result)
			return &message.Reply{message.MsgTypeText, reply}
		case message.MsgTypeEvent:

			reply := message.NewText("Hi 亲爱的国学爱好者，" + "谢谢你的关注！我是您的较为智能的国学小助手。尝试回复诗句或者词牌名，看看都有什么吧！" )
			return &message.Reply{message.MsgTypeText, reply}
		}
		return &message.Reply{message.MsgTypeText, "没有找到哦，亲~\n"}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	accessToken, err = server.GetAccessToken()

	if err != nil {
		log.Warn(err)
	}

	//发送回复的消息
	server.Send()
}

