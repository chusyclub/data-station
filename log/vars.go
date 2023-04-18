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
	log           *logrus.Logger
	uniqueId      string
	extra         map[string]interface{}
	defLogPath    = "/tmp/logs/"
	defMaxRolls   = 5
	FlagTrade     = "[quant-trade]"
	LogTypeCommon = 1 // 普通类型
	LogTypeTrade  = 2 // 交易类型
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
	LogTypeField = "log_type"
	TradeIdField = "trade_id"

	defLogFile = "default.log"
)

const (
	envLogToFile   = "LOG_TO_FILE"   // write log to file
	envLogPath     = "LOG_PATH"      // log to file path
	envLogRoll     = "LOG_ROLL"      // log to file path
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
	LogType  int    `json:"log_type"` // 类别 1.普通，2.交易日志
	TradeId  string `json:"trade_id"` // 交易ID
}

func (d SysLog) TableName() string {
	return "sys_logs"
}
