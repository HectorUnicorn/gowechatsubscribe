package rpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"net"
	"github.com/astaxie/beego"
)

const (
	Create = iota  // 完全按照诗句去创建图片
	Overlay   // Overlay在图片库中的图片
	CreateFromSource // 在传入的图片上处理
)

type ImageCreation struct {
	Action  int
	Poetry  string
	Author  string
	Danysty string
	Tags    []string
}

func (i *ImageCreation) ToString() {
	params := "Action=" + i.Action + "&Poetry=" + i.Poetry + "&Author=" + i.Author + "&Danysty=" + i.Danysty + "&Tags=" + i.Tags
}

// action - create image
//
func NewJsonRpcSocketClient(creation ImageCreation) (string, error) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3456", 1000*1000*1000*30)
	if err != nil {
	    beego.Error("create client err", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)
	var filePath string
	err = client.Call("CreateImage.Create", )
}
