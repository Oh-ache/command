//
// main.go
// Copyright (C) 2021 ache <1751987128@qq.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type responseData struct {
	From        string  `json:"from"`
	To          string  `json:"to"`
	TransResult []*List `json:"trans_result"`
}
type List struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

var BaiDuUrl = "http://api.fanyi.baidu.com/api/trans/vip/translate"
var BaiDuAppId = "你的appid"
var BaiDuSercert = "你的密钥"

func main() {
	list1 := os.Args
	str := string(list1[1])

	getApiInfo(url.QueryEscape(str), checkLetter(str))
}

// 判断是否汉字
func checkLetter(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

// 请求百度翻译接口
func getApiInfo(q string, isHan bool) error {
	to := ""
	if isHan {
		to = "en"
	} else {
		to = "zh"
	}

	client := &http.Client{}

	appid := BaiDuAppId
	salt := getNum()
	sign := Md5(appid + q + salt + BaiDuSercert)
	// go %特殊字符串处理
	q = strings.Replace(url.QueryEscape(q), "%", "%%", -1)

	//生成要访问的url,token是api鉴权，每个api访问方式不同，根据api调用文档拼接URL
	str := fmt.Sprintf("?q=%v&from=auto&to=%s&appid=%s&salt=%s&sign=%s", q, to, appid, salt, sign)
	url := fmt.Sprintf(BaiDuUrl + str)
	//提交请求
	request, err := http.NewRequest("GET", url, nil)
	//异常捕捉
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, _ := client.Do(request)
	//关闭流
	defer response.Body.Close()
	//检出结果集
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	r := responseData{}
	json.Unmarshal(body, &r)
	fmt.Println(r.TransResult[0].Dst)
	return nil
}

// 获取随机数
func getNum() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	return strconv.Itoa(num)
}

// md5加密
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1
}
