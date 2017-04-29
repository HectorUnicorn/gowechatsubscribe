package test

import (
	"testing"
	"gowechatsubscribe/dblite"
	"fmt"
)

func TestMatchTag(t *testing.T) {
	tag := "[捂脸]"
	mgr := dblite.NewDBManager()
	content := mgr.SelectPoetry(tag)
	fmt.Printf("tag %s matched content is %s", tag, content)
	t.Error("test matchtag")
}
