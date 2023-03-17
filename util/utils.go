package util

import "time"

/**
 * @Author: chushiyang
 * @Description: 常用工具方法
 * @File:  utils
 * @Version: 1.0.0
 * @Date: 2023/3/17 23:47
 */

func GetCurrTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
