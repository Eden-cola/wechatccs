package ticket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ExpiresIn int = 7200
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

func (s saveStruct) ToJson() jsonStruct {
	return jsonStruct{0, "", s.Ticket, s.ExpiresIn()}
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

func (s *saveStruct) update(url string, jsonS jsonStruct) error {
	//fmt.Println("update access_token")
	jsonStr := s.httpGet(url)
	err := json.Unmarshal([]byte(jsonStr), &jsonS)
	if err == nil {
		s.load(jsonS)
	}
	return err
}

func (s *saveStruct) load(jsonS jsonStruct) {
	s.Ticket = jsonS.Ticket
	s.ExpiresTime = jsonS.ExpiresIn + int(time.Now().Unix())
}

func init() {
	saveList = make(map[string]saveStruct)
}

func (s saveStruct) httpGet(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
		err := v.update(getUrl(accessToken), jsonStruct{})
		if err == nil {
			saveList[accessToken] = v
		}
		return v, err
	}
}

func Get(accessToken string) (jsonStruct, error) {
	ticket, err := get(accessToken)
	if err != nil {
		return jsonStruct{}, err
	}
	return ticket.ToJson(), nil
}

//检查存储map，清理过期的ticket,返回多少秒后应该进行下一次检查
func Check() int {
	fmt.Println("ticket check")
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
