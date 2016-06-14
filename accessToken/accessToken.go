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

func (s saveStruct) update(url string, jsonS jsonStruct) error {
	fmt.Println("update")
	jsonStr := httpGet(url)
	err := json.Unmarshal([]byte(jsonStr), &jsonS)
	if err == nil {
		s = jsonS.ToSave()
	}
	return err
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
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
		err := v.update(getUrl(appID, appsecret), jsonStruct{})
		fmt.Printf("%+v\n", v)
		if err == nil {
			saveList[key] = v
		}
		return v, err
	}
}

func Get(appID, appsecret string) (jsonStruct, error) {
	accessToken, err := get(appID, appsecret)
	if err != nil {
		return jsonStruct{}, err
	}
	return accessToken.ToJson(), nil
}
