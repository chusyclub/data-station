package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: chushiyang
 * @Description: log util
 * @File:  log.go
 * @Version: 1.0.0
 * @Date: 2022/12/20 14:40
 */

func init() {
	log = logrus.New()

	// set log formatter
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "[T]",
		logrus.FieldKeyLevel: "[L]",
		logrus.FieldKeyMsg:   "[Msg]",
	}
	log.SetFormatter(formatter)

	var isDev bool
	switch strings.ToLower(os.Getenv(envEnv)) {
	case envDev:
		isDev = true
	case envBeta:
		isDev = false
	case envPro:
		isDev = false
	case envTest:
		isDev = false
	default:
		isDev = true
	}

	// add es log hook
	if server := os.Getenv(envLogEsServer); server != "" {
		log.AddHook(NewEsLogHook(server))
	}

	// set default roll number
	if roll := os.Getenv(envLogRoll); roll != "" {
		defMaxRolls, _ = strconv.Atoi(roll)
	}

	// add line hook
	log.AddHook(NewLineHook(isDev))

	SetExtra(map[string]interface{}{
		"command": os.Args[0],
	})

}

// ---------------------------------------------------------------------------------------------------------------------

func Init(serviceName, logDBUrl string) error {
	// add service name hook
	svcName := ""
	if serviceName == "" {
		if os.Getenv(envServiceName) != "" {
			svcName = os.Getenv(envServiceName)
		}
	} else {
		svcName = serviceName
	}
	// add service name hook
	if svcName != "" {
		log.Infof("[%s] start...", serviceName)
		log.AddHook(NewAppHook(svcName))
		if logDBUrl != "" {
			// add log to db hook
			mysqlHook := NewMySQLHook(logDBUrl, svcName)
			go mysqlHook.writeToDb()
			log.AddHook(mysqlHook)

		}

		// write to file
		if file := os.Getenv(envLogToFile); file != "" {
			logFilename := getLogFilename(svcName)
			if path := filepath.Dir(logFilename); !isDirExists(path) {
				if err := os.MkdirAll(path, 0744); err != nil {
					log.Errorf("Mkdirall %s err: %v", path, err)
				}
			}
			rotateLog, _ := rotatelogs.New(
				logFilename+".%Y%m%d",
				rotatelogs.WithLinkName(logFilename),
				rotatelogs.WithMaxAge(24*time.Duration(defMaxRolls)*time.Hour),
				rotatelogs.WithRotationTime(24*time.Hour),
			)
			log.Out = rotateLog
		}
	}

	return nil
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetExtra(h map[string]interface{}) {
	extra = h
	RefreshUniqueId()
}

func RefreshUniqueId() {
	uniqueId = xid.New().String()
}

func Type(typ string) *logrus.Entry {
	return log.WithField(typeField, typ)
}

func With(pairs ...string) *logrus.Entry {
	if len(pairs)%2 != 0 {
		pairs = append(pairs, "unknown")
	}

	fields := logrus.Fields{}
	for i := 0; i < len(pairs); i += 2 {
		fields[pairs[i]] = pairs[i+1]
		if pairs[i] == "type" {
			fields[typeField] = pairs[i+1]
		}
	}

	return log.WithFields(fields)
}

func WithField(key string, value interface{}, typ string) *logrus.Entry {
	return log.WithField(typeField, typ).WithField(key, value)
}

func WithFields(fields logrus.Fields, typ string) *logrus.Entry {
	fields[typeField] = typ
	return log.WithFields(fields)
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	log.Infoln(v...)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	log.Debugln(v...)
}

func Warn(v ...interface{}) {
	log.Warning(v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

func Warnln(v ...interface{}) {
	log.Warnln(v...)
}

func Warning(v ...interface{}) {
	log.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

func Warningln(v ...interface{}) {
	log.Warningln(v...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	log.Errorln(v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

// ---------------------------------------------------------------------------------------------------------------------

func isDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func getLogFilename(serviceName string) string {
	if path := os.Getenv(envLogPath); path != "" {
		defLogPath = envLogPath
	}
	var logFilename string
	logFilename = defLogPath + serviceName + "/" + defLogFile

	if !filepath.IsAbs(logFilename) {
		logFilename, _ = filepath.Abs(logFilename)
	}

	if path := filepath.Dir(logFilename); !isDirExists(path) {
		if err := os.MkdirAll(path, 0744); err != nil {
			log.Errorf("Mkdirall %s err: %v", path, err)
		}
	}
	fmt.Println("logFilename==", logFilename)
	return logFilename
}
