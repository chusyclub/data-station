package mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
	"time"
)

/**
 * @Author: chushiyang
 * @Description:
 * @File:  mysql
 * @Version: 1.0.0
 * @Date: 2022/12/20 17:08
 */

// 数据库引擎对应关系
var (
	engineMap      map[string]*xorm.Engine
	engineMapMutex sync.RWMutex
)

const (
	MaxOpenConn   = 100
	MaxIdleConn   = 40
	defaultDriver = "mysql"
)

func init() {
	if engineMap == nil {
		engineMap = make(map[string]*xorm.Engine)
	}
}

func GetOrmEngine(url string) (*xorm.Engine, error) {
	if url == "" {
		fmt.Println("datasource url is nil")
		return nil, errors.New("datasource url is nil")
	}
	var engine *xorm.Engine

	dataSourceName := url

	engineMapMutex.RLock()
	engine, ok := engineMap[dataSourceName]
	engineMapMutex.RUnlock()
	if ok {
		return engine, nil
	}

	engine, err := xorm.NewEngine(defaultDriver, dataSourceName)
	if nil != err {
		return nil, err
	}
	engine.SetMaxOpenConns(MaxOpenConn)
	engine.SetMaxIdleConns(MaxIdleConn)

	location, _ := time.LoadLocation("Asia/Shanghai")
	engine.DatabaseTZ = location
	engine.ShowSQL(true)

	err = engine.Ping()

	if err != nil {
		err := engine.Close()
		if err != nil {
			return nil, err
		}
	}

	engineMapMutex.Lock()
	defer engineMapMutex.Unlock()
	engineMap[dataSourceName] = engine
	return engine, nil
}
