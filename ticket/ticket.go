package ticket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type jsonStruct struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func (j jsonStruct) ExpiresTime() int {
	timeStamp := time.Now().Unix()
	return j.ExpiresIn + int(timeStamp)
}

func (j jsonStruct) ToSave() saveStruct {
	return saveStruct{j.Ticket, j.ExpiresTime()}
}

type resultStruct struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

type saveStruct struct {
	Ticket      string
	ExpiresTime int
}

var saveList map[string]saveStruct

func (s saveStruct) ToResult() resultStruct {
	return resultStruct{s.Ticket, s.ExpiresIn()}
}

func (s saveStruct) Check() bool {
	timeStamp := time.Now().Unix()
	if s.ExpiresTime > int(timeStamp) {
		return true
	}
	return false
}

func (s saveStruct) ExpiresIn() int {
	timeStamp := time.Now().Unix()
	return s.ExpiresTime - int(timeStamp)
}

func init() {
	saveList = make(map[string]saveStruct)
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func getUrl(accessToken string) string {
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi"
	return string(url)
}

func get(accessToken string) (saveStruct, error) {
	if v, ok := saveList[accessToken]; ok && v.Check() {
		return v, nil
	} else {
		s, err := update(accessToken)
		if err != nil {
			return saveStruct{}, err
		} else {
			saveList[accessToken] = s
			return s, nil
		}
	}
}

func update(accessToken string) (saveStruct, error) {
	fmt.Println("update ticket")
	url := getUrl(accessToken)
	jsonStr := httpGet(url)
	var data jsonStruct
	if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
		return data.ToSave(), nil
	} else {
		return saveStruct{}, err
	}
}

func Get(accessToken string) (resultStruct, error) {
	ticket, err := get(accessToken)
	if err != nil {
		return resultStruct{}, err
	}
	return ticket.ToResult(), nil
}
