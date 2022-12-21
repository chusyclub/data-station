package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
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

	// add line hook
	log.AddHook(NewLineHook(isDev))

	SetExtra(map[string]interface{}{
		"command": os.Args[0],
	})

	// set log level
	if lvl, err := logrus.ParseLevel(envLogLevel); err == nil {
		log.SetLevel(lvl)
	} else {
		if isDev {
			log.SetLevel(DebugLevel)
		} else {
			log.SetLevel(InfoLevel)
		}
	}
	// write to file
	if file := os.Getenv(envLogToFile); file != "" {
		logFilename := getLogFilename()
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
		log.AddHook(NewAppHook(svcName))
		if logDBUrl != "" {
			// add log to db hook
			mysqlHook := NewMySQLHook(logDBUrl, svcName)
			go mysqlHook.writeToDb()
			log.AddHook(mysqlHook)

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

func getLogFilename(filename ...string) string {
	var serviceName, logFilename string

	if len(filename) == 0 {
		if serviceName = os.Getenv(envServiceName); serviceName == "" {
			serviceName = "undefined"
		}
		logFilename = defLogPath + "_" + serviceName + "/" + defLogFile
	} else {
		logFilename = filename[0]
	}

	if !filepath.IsAbs(logFilename) {
		logFilename, _ = filepath.Abs(logFilename)
	}

	if path := filepath.Dir(logFilename); !isDirExists(path) {
		if err := os.MkdirAll(path, 0744); err != nil {
			log.Errorf("Mkdirall %s err: %v", path, err)
		}
	}

	return logFilename
}