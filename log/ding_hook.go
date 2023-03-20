package log

import (
	"fmt"
	"github.com/chusyclub/data-station/util"
	"github.com/sirupsen/logrus"
	"strings"
)

/**
 * @Author: chushiyang
 * @Description: 钉钉群机器人报警群
 * @File:  ding_hook.go
 * @Version: 1.0.0
 * @Date: 2023/3/17 23:38
 */

type DingDingHook struct {
	ServiceName string
	DingFlag    string
	DingUrl     string
	DingType    string
	DingAts     []string
}

func NewDingHook(serviceName, flag, dingUrl, dingType string, dingAts []string) *DingDingHook {
	hook := &DingDingHook{
		ServiceName: serviceName,
		DingFlag:    flag,
		DingUrl:     dingUrl,
		DingType:    dingType,
		DingAts:     dingAts,
	}
	return hook
}

func (d *DingDingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (d *DingDingHook) Fire(entry *logrus.Entry) error {
	// send message to ding ding robot
	level := strings.ToUpper(entry.Level.String())
	if level == "ERROR" {
		content := fmt.Sprintf("[%s]项目[%s]在时间[%s]报错【%s】", d.DingFlag, d.ServiceName, util.GetCurrTimeStr(), entry.Message)
		util.SendDingTalkText(d.DingAts, content, d.DingType, d.DingUrl)
	}
	return nil
}
