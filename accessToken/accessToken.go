package accessToken

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ExpiresIn int = 7200
)

//json对应单元
type jsonStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

//var jsonMap map[string]interface{}

//存储单元
type saveStruct struct {
	AccessToken string
	ExpiresTime int
}

//检查是否过期//待废弃
func (s saveStruct) Check() bool {

	if s.ExpiresIn() > 1 {
		return true
	}
	return false
}

//距离过期剩余时间
func (s saveStruct) ExpiresIn() int {
	timeStamp := time.Now().Unix()
	return s.ExpiresTime - int(timeStamp)
}

//更新存储对象信息
func (s *saveStruct) update(url string) error {
	//jsonM := jsonMap
	var jsonS jsonStruct
	jsonStr := s.httpGet(url)
	err := json.Unmarshal([]byte(jsonStr), &jsonS)
	if err == nil {
		s.load(jsonS)
	}
	//fmt.Printf("result: %+v\n", s)
	return err
}

//从json对应单元中加载数据
func (s *saveStruct) load(jsonS jsonStruct) {
	s.AccessToken = jsonS.AccessToken
	s.ExpiresTime = jsonS.ExpiresIn + int(time.Now().Unix())
	//s.ExpiresTime = 25 + int(time.Now().Unix())
}

func (s saveStruct) httpGet(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

var saveList map[[16]byte]saveStruct //存储map
var preUploadLock map[[16]byte]bool  //预加载锁

func init() {
	saveList = make(map[[16]byte]saveStruct)
	preUploadLock = make(map[[16]byte]bool)
}

func getUrl(appID, appsecret string) string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appID + "&secret=" + appsecret
	return string(url)
}

func preUpdate(appID, appsecret string) {
	key := md5.Sum([]byte(appID + appsecret))
	if _, ok := preUploadLock[key]; ok {
		return
	} else {
		preUploadLock[key] = true
	}
	var v saveStruct
	err := v.update(getUrl(appID, appsecret))
	if err == nil {
		saveList[key] = v
	}
	delete(preUploadLock, key)
}

func get(appID, appsecret string) (saveStruct, error) {
	key := md5.Sum([]byte(appID + appsecret))
	if v, ok := saveList[key]; ok && v.ExpiresIn() > 2 {
		if v.ExpiresIn() < 10 {
			go preUpdate(appID, appsecret)
		}
		return v, nil
	} else {
		err := v.update(getUrl(appID, appsecret))
		if err == nil {
			saveList[key] = v
		}
		return v, err
	}
}

//获取某个appid的json对应对象
func Get(appID, appsecret string) (jsonStruct, error) {
	a_t, err := get(appID, appsecret)
	if err != nil {
		return jsonStruct{}, err
	}
	return jsonStruct{a_t.AccessToken, a_t.ExpiresIn()}, nil
}

//检查存储map，清理过期的accessToken,返回多少秒后应该进行下一次检查
func Check() int {
	fmt.Println("accessToken check")
	firstIn := ExpiresIn //第一个将要过期的时间
	for key, value := range saveList {
		ei := value.ExpiresIn()
		if ei < 1 {
			delete(saveList, key)
		} else {
			if ei < firstIn {
				firstIn = ei
			}
		}
	}
	return firstIn
}
