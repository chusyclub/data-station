package log

import (
	"github.com/sirupsen/logrus"
)

/**
 * @Author: chushiyang
 * @Description:
 * @File:  vars.go
 * @Version: 1.0.0
 * @Date: 2022/12/20 14:40
 */

var (
	log      *logrus.Logger
	uniqueId string
	extra    map[string]interface{}
)

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)

const (
	typeField  = "log_type"
	defLogPath = "/tmp/logs/"
	defLogFile = "default.log"
)

const (
	defMaxRolls    = 5               // log file roll
	envLogToFile   = "LOG_TO_FILE"   // write log to file
	envLogEsServer = "LOG_ES_SERVER" // log server
	envServiceName = "SERVICE_NAME"  // service name
	envEnv         = "ENV"           // env
	envDev         = "dev"           // dev
	envBeta        = "beta"          // beta
	envTest        = "test"          // test
	envPro         = "product"       // product
)

type SysLog struct {
	Id       int64  `json:"id"`
	Project  string `json:"project"`
	Level    string `json:"level"`
	DataTs   int64  `json:"data_ts"`
	DataTime string `json:"data_time"`
	Content  string `json:"content"`
}

func (d SysLog) TableName() string {
	return "sys_logs"
}
