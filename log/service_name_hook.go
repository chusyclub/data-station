package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

/**
 * @Author: chushiyang
 * @Description:
 * @File:  app_name_hook
 * @Version: 1.0.0
 * @Date: 2022/12/20 17:13
 */

type ServiceHook struct {
	ServiceName string
}

func NewAppHook(serviceName string) *ServiceHook {
	fmt.Println("初始化service name hook...")
	return &ServiceHook{
		ServiceName: serviceName,
	}
}

func (h *ServiceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ServiceHook) Fire(entry *logrus.Entry) error {
	entry.Data["[service]"] = h.ServiceName
	return nil
}
