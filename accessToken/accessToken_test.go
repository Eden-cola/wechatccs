package accessToken

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	appid := "wxf1563182e6379566"
	secret := "d4624c36b6795d1d99dcf0547af5443d"
	data, err := Get(appid, secret)
	if err != nil {
		t.Errorf("出现错误")
		fmt.Println(err)
	}
	if data.AccessToken == "" {
		t.Errorf("未获取到accessToken")
	}
	t.Errorf("查看结果")
	fmt.Println(data)
}
