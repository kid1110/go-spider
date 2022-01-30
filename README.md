# go-spider
使用go并发爬取qq新闻网站（学习使用）

获取需要的库

```go
go get github.com/go-sql-driver/mysql
go get github.com/sirupsen/logrus
go get github.com/opesun/goquery
go get github.com/axgle/mahonia
go get github.com/garyburd/redigo/redis
```

本项目属于学习go语言的一个比较入门的项目，本次的项目使用的都是比较原生的库，很多都是自己封装完成的，学习中可以加入其他框架以方便开发，本项目仍有许多的不足，本人仍在努力学习中。



##  项目描述

  开启协程并发爬取qq新闻的数据，保存入mysql中，客户使用的时候会从mysql中取出，并且使用了redis作为缓存



##  项目技术栈

 ###  mysql

* 建库建表脚本在**/scripts/mysql**文件夹中，使用了索引进行优化

* 使用连接池作为连接

### redis

* 自己封装了redis的增删改查以及序列化的函数
* redis使用了连接池作为连接
* 对缓存击穿，缓存穿透，缓存雪崩做了处理(采用加锁，布隆过滤器，随机存储时间)

### log

* 这部分采用的是logrus库，但使用不熟悉没有输出到文件中，需要优化

###  exception

*  go-web使用的是自定义的错误处理以及闭包统一错误处理

### go-web

* 本项目使用的是原生的net/http库作为go-web



## 项目特点

* 并发爬取QQ新闻的文章，通过使用索引和redis大大加快了爬取，查询的速度

*  对错误做了统一处理



##  部署

本后台使用的是docker部署，dockerfile在**/scripts/docker**中,部署方式是在window交叉编译了main函数，将main二进制文件和config.yaml放到go-web文件夹下。

![image-20220130140819121](https://gitee.com/kid1110/Imageshack/raw/master/img/image-20220130140819121.png)

在go-web文件下交叉编译

```shell
# Windows 下编译 Mac 和 Linux 64位可执行程序
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```

```shell
#Mac 下编译 Linux 和 Windows 64位可执行程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```

部署命令

```shell
sudo docker build -t go-web .
sudo docker run -di --name go-web -p 8081:8081 go-web
```



##  测试

请求头为 

![image-20220130141319861](https://gitee.com/kid1110/Imageshack/raw/master/img/image-20220130141319861.png)

![image-20220130141422448](https://gitee.com/kid1110/Imageshack/raw/master/img/image-20220130141422448.png)

![image-20220130141447420](https://gitee.com/kid1110/Imageshack/raw/master/img/image-20220130141447420.png)

![image-20220130141543235](https://gitee.com/kid1110/Imageshack/raw/master/img/image-20220130141543235.png)

