package buptncov

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Buptncov struct {
	Client  http.Client
	Cookies []*http.Cookie
}

func New() *Buptncov {
	return &Buptncov{http.Client{Timeout: 10 * time.Second}, []*http.Cookie{}}
}
func (self *Buptncov) Login(username string, password string) error {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	req, _ := http.NewRequest("POST", "https://app.bupt.edu.cn/uc/wap/login/check", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := self.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return err2
	}
	response := GeneralResponse{}
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		fmt.Println("error:", json_err)
	}
	if response.Error == 0 {
		fmt.Println("登录成功")
		self.Cookies = resp.Cookies()
	} else {
		fmt.Println(string(body))
		return errors.New("登录失败")
	}
	return nil
}

func (self *Buptncov) GetAndPostInfo() error {
	// 获取提交信息
	req, _ := http.NewRequest("GET", "https://app.bupt.edu.cn/xisuncov/wap/open-report/index", nil)
	for i := range self.Cookies {
		req.AddCookie(self.Cookies[i])
	}
	resp, requestErr := self.Client.Do(req)
	if requestErr != nil {
		return requestErr
	}
	defer resp.Body.Close()
	body, parseErr := io.ReadAll(resp.Body)
	if parseErr != nil {
		return parseErr
	}
	var map_resp GeneralResponse
	jsonErr := json.Unmarshal(body, &map_resp)
	if jsonErr != nil {
		return jsonErr
	}
	// get请求会返回之前提交的情况，这里只要改一下日期date和时间段flag即可。
	// 返回的请求格式大致如下：
	// - e
	// - m
	// - d
	//		- date 当前日期
	//		- info 之前提交的情况
	//			- sfzx
	//			- tw 体温
	//			- area 区域
	//			- city 城市
	//			- province 省份
	//			- address 地址
	//			- ....
	//			- flag 时间段
	// ......
	// ......

	result := map_resp.Data["info"].(map[string]interface{})
	date, _ := map_resp.Data["date"].(string)
	period := date[11:]

	// 修改时间段
	var flag int
	if strings.Contains(period, "上午") {
		flag = 0
	} else if strings.Contains(period, "下午") {
		flag = 1
	} else if strings.Contains(period, "晚上") {
		flag = 2
	}
	result["flag"] = flag
	// 修改日期
	result["date"] = strings.Replace(date[:10], "-", "", -1)

	jsonString, _ := json.Marshal(result)
	fmt.Println("提交表单：", string(jsonString))

	// 提交
	submitreq, _ := http.NewRequest("GET", "https://app.bupt.edu.cn/xisuncov/wap/open-report/save", bytes.NewBuffer(jsonString))
	for i := range self.Cookies {
		submitreq.AddCookie(self.Cookies[i])
	}
	submitreq.Header.Set("Content-Type", "application/json")
	submitresp, requestErr := self.Client.Do(submitreq)
	if requestErr != nil {
		return requestErr
	}
	submitresult, submitresulterr := io.ReadAll(submitresp.Body)
	if submitresulterr != nil {
		return submitresulterr
	}
	fmt.Println("提交结果：", string(submitresult))
	var submitResultJson GeneralResponse
	submitResultJson_err := json.Unmarshal(submitresult, &submitResultJson)
	if submitResultJson_err != nil {
		return submitResultJson_err
	}
	if submitResultJson.Error == 1 || submitResultJson.Error == 0 {
		fmt.Println("提交成功")
	} else {
		return errors.New("提交失败")
	}
	return nil
}
