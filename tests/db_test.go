package test

import (
	"testing"
	"gowechatsubscribe/dblite"
	"fmt"
)

func TestSelectPoetry(t *testing.T) {
	db := dblite.NewDBManager()
	result := db.SelectPoetry("伐檀")
	fmt.Println("result:", result)
	if len(result) > 12 {
		t.Log("testing passed")
	} else {
		t.Error("testing not passed")
	}
}