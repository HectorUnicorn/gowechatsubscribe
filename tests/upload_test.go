package test

import (
	"gowechatsubscribe/uploader"
	"testing"
)

func TestDoUploadFile2Wechat(t *testing.T) {
	accessToken := "BjzDv40Dnst3jwjR7zGD9ivaOK0jdqxAscHgabLaPHOL15dGTcqoEXXetZeWa82eAPZo5nZ0ILmUj93k2zpbE7mF7ozhgclzztmAOxCIxrU2y0cpcL-6m0hkWNO0JO6RFSHcAHAPCR"
	uploader.DoUploadFile2Wechat("../qingming.jpg", accessToken)
	t.Error("test upload")
}
