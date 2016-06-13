package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
	"wechatccs/accessToken"
	"wechatccs/ticket"
)

func getToken(c echo.Context) error {
	appid := c.QueryParam("appid")
	secret := c.QueryParam("secret")
	accessToken, err := accessToken.Get(appid, secret)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return c.JSON(http.StatusCreated, accessToken)
	}
}

func getTicket(c echo.Context) error {
	token := c.QueryParam("access_token")
	ticket, err := ticket.Get(token)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return c.JSON(http.StatusCreated, ticket)
	}
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		html := "<p>获取access_token请访问: /token?appid=APPID&secret=APPSECRET</p>"
		html += "<p>获取jsapi_ticket请访问: /ticket?access_token=ACCESS_TOKEN</p>"
		return c.HTML(http.StatusOK, html)
	})
	e.GET("/token", getToken)
	e.GET("/ticket", getTicket)
	fmt.Println("listen:2266")
	e.Run(standard.New(":2266"))
}
