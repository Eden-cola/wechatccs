package ticket

import (
	"fmt"
	"testing"
	"wechatccs/accessToken"
)

func TestGet(t *testing.T) {
	appid := "wxf1563182e6379566"
	secret := "d4624c36b6795d1d99dcf0547af5443d"
	a_t, err := accessToken.Get(appid, secret)
	tic, err := Get(a_t.AccessToken)
	if err != nil {
		t.Errorf("出现错误")
		fmt.Println(err)
	}
	if tic.Ticket == "" {
		t.Errorf("未获取到ticket")
	}
	if tic.ExpiresIn != 7200 {
		t.Errorf("过期时间不正确")
	}
}
