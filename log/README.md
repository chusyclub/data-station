## Log


based on logrus, support sync to db, es, file

#### Config

### Usage

```
import (
	"github.com/chusyclub/data-station/log"
)

func main() {
	dbUrl := ""
	serviceName := "demo-sdk"
	log.Init(serviceName, dbUrl)
	log.Info("info")
	log.Infof("info-%d", 1)
	log.Error("error")
	log.Debug("debug")
	log.Warn("warn")
	select {}
}
```
you can also set service name by set environment `export SERVICE_NAME=demo-sdk`.  
if you want to write log to file, you can set environment `export LOG_TO_FILE=true`,  
then log default write to `{LOG_PATH}/{SERVICE_NAME}/default.log`.  
`LOG_PATH` default value is `/tmp/logs/`, or you can update by `export LOG_PATH=/tmp/service_logs/`.  
history log named by `logFilename+".%Y%m%d"`, default roll number is 5 , you can use `export LOG_ROLL=10` change it.

