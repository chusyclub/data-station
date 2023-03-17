package log

import (
	"data-station/util"
	"fmt"
	"github.com/chusyclub/data-station/mysql"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

/**
 * @Author: chushiyang
 * @Description:
 * @File:  mysql_hook
 * @Version: 1.0.0
 * @Date: 2022/12/20 17:13
 */

type MySQLHook struct {
	Url        string
	SvcName    string
	SysLogChan chan SysLog
}

func NewMySQLHook(url, svcName string) *MySQLHook {
	hook := &MySQLHook{
		Url:        url,
		SvcName:    svcName,
		SysLogChan: make(chan SysLog, 100),
	}
	return hook
}

func (h *MySQLHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *MySQLHook) Fire(entry *logrus.Entry) error {
	log := SysLog{
		Project:  h.SvcName,
		Level:    strings.ToUpper(entry.Level.String()),
		DataTs:   time.Now().UnixMilli(),
		DataTime: util.GetCurrTimeStr(),
		Content:  entry.Message,
	}
	h.SysLogChan <- log
	return nil
}

func (h *MySQLHook) writeToDb() {
	engine, err := mysql.GetOrmEngine(h.Url)
	if err != nil {
		fmt.Println("write to db error:", err)
	}
	session := engine.NewSession()
	defer session.Close()
	for log := range h.SysLogChan {
		count, err := session.InsertOne(log)
		if err != nil {
			fmt.Println("insert result, count=", count, ", error=", err)
		}
	}

}
