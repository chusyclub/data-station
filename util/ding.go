package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/**
 * @Author: chushiyang
 * @Description: 钉钉消息写入群
 * @File:  ding
 * @Version: 1.0.0
 * @Date: 2023/3/17 23:43
 */

type DingParams struct {
	MsgType string              `json:"msgtype"`
	Text    map[string]string   `json:"text"`
	At      map[string][]string `json:"at"`
}

func SendDingTalkText(at []string, content, dingType, dingUrl string) {
	dp := DingParams{
		MsgType: dingType,
		Text:    map[string]string{"content": content},
		At:      map[string][]string{"atMobiles": at},
	}
	dpStr, _ := json.Marshal(dp)
	req, err := http.NewRequest("POST", dingUrl, bytes.NewBuffer(dpStr))
	if err != nil {
		fmt.Printf("\n获取钉钉请求req错误，url=[%s], params=[%s],err=[%v]\n", dingUrl, dpStr, err)
	}
	req.Header.Set("Content-Type", "application/json")
	var client = http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n钉钉请求错误，url=[%s], params=[%s]\n, err=[%v]", dingUrl, dpStr, err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("\n解析钉钉请求返回结果错误，url=[%s], params=[%s], err=[%v]\n", dingUrl, dpStr, err)
		} else {
			fmt.Printf("call dingTalk response:%s", string(body))
		}
	}
}
