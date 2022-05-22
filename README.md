
## 目录结构

1. api文件夹就是MVC框架的controller，负责协调各部件完成任务
2. model文件夹负责存储数据库模型和数据库操作相关的代码
3. service负责处理比较复杂的业务（可能用不到，为了简化）
4. util一些通用的小工具
5. conf放一些静态存放的配置文件
6. public放视频和图片，音频等的静态资源目录
7. e-r.png是数据库的e-r图
8. .env是环境配置，如果想修改数据库连接则改这里
9. middleware是中间件，用于登陆验证，跨域操作
10. serializer准备放一些序列化统一格式的操作，暂时还没写好
11. server放要挂载的路由

## 项目运行流程
1. 首先进入main.go
2. 然后在conf.go中初始化数据库链接
3. 接下来挂载路由，去server/router.go中
4. 现在我已经写好了基本的路由框架，你们就专注于每个路由对应的api函数具体实现就可以了

## 写一个接口的流程
1. 首先看自己负责的路由调用的函数
2. 然后看接口的入参和响应
3. 根据响应格式写结构体
4. 在自己的api中调用的函数写代码

## 我写了两个demo可以在router文件和api文件看到，照着写就行


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
