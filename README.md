## 目录结构

1. api文件夹就是MVC框架的controller，负责协调各部件完成任务
2. model文件夹负责存储数据库模型和数据库操作相关的代码
3. service负责处理比较复杂的业务
4. serializer储存通用的json模型，把model得到的数据库模型转换成api需要的json对象
5. util一些通用的小工具
6. conf放一些静态存放的配置文件

## 环境

项目在启动的时候依赖以下环境变量，但是也可以在项目根目录创建.env文件设置环境变量便于使用(建议开发环境使用)

```
MYSQL_DSN="qxy_dy:123456@tcp(47.95.23.74:3306)/qxy_dy?charset=utf8&parseTime=True&loc=Local" # Mysql连接地址
SESSION_SECRET="womeishijiuchiliuliumei" # Seesion密钥，必须设置而且不要泄露
GIN_MODE="debug"
```


## 运行

```
go run main.go
```
