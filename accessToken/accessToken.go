package accessToken

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type jsonStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type saveStruct struct {
	AccessToken string
	ExpiresTime int
}

var saveList map[[16]byte]saveStruct

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

func (s saveStruct) ToJson() jsonStruct {
	return jsonStruct{s.AccessToken, s.ExpiresIn()}
}

func (j jsonStruct) ExpiresTime() int {
	timeStamp := time.Now().Unix()
	return j.ExpiresIn + int(timeStamp)
}

func (j jsonStruct) ToSave() saveStruct {
	return saveStruct{j.AccessToken, j.ExpiresTime()}
}

func init() {
	saveList = make(map[[16]byte]saveStruct)
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

func getUrl(appID, appsecret string) string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appID + "&secret=" + appsecret
	return string(url)
}

func get(appID, appsecret string) (saveStruct, error) {
	key := md5.Sum([]byte(appID + appsecret))
	if v, ok := saveList[key]; ok && v.Check() {
		return v, nil
	} else {
		s, err := update(appID, appsecret)
		if err != nil {
			return saveStruct{}, err
		} else {
			saveList[key] = s
			return s, nil
		}
	}
}

func update(appID, appsecret string) (saveStruct, error) {
	fmt.Println("update access_token")
	url := getUrl(appID, appsecret)
	jsonStr := httpGet(url)
	var data jsonStruct
	if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
		return data.ToSave(), nil
	} else {
		return saveStruct{}, err
	}
}

func Get(appID, appsecret string) (jsonStruct, error) {
	accessToken, err := get(appID, appsecret)
	if err != nil {
		return jsonStruct{}, err
	}
	return accessToken.ToJson(), nil
}
