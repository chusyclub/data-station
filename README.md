# data-station
sdk内容为data-station项目使用

1. 钉钉推送配置方式
```
[DingDing]
DingFlag="日志"
DingTalkUrl="https://oapi.dingtalk.com/robot/send?access_token=*****"
DingType="text"
Ats=[]
```
2. 数据库配置
```
[DB]
QuantMainURL="user:password@(ip:port)/dbname?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8"
...
```
3. main.go中添加
```
// 在初始化config以后调用，传入config对象
log.Init("项目名", *config.Conf)
```