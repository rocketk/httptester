近来有需求要对HTTP接口进行压测，于是去了解了一下JMeter，发现虽然功能强大，但本身依赖JAVA，并且依赖图形化界面，不够轻量化，所以想着自己写一个基于命令行的轻量级压测工具，于是就有了httptester。

httptester是采用GO语言来编写的，尽管目前仍是早期版本，但是基本功能已经没有问题了。

### 快速开始

httptester是一个二进制可执行文件，无需安装。

下载最新版本：https://github.com/rocketk/httptester/releases

找到自己平台对应的版本下载至本地，将其所在目录放入系统环境变量中。

#### 常用命令

获取帮助信息

```shell
httptester -h
httptester run -h
httptester serve -h
```

最简单的压测命令（以压测百度首页为例）

```shell
httptester run -u 'https://www.baidu.com/'
```

设定并发量、循环次数、超时时间

```shell
httptester run --loop 10 --concurrency 100 --timeout 500ms -u 'https://www.baidu.com/'
#httptester run --loop 10 --concurrency 100 --timeout 2s -u 'https://www.baidu.com/'
```

#### 结果

```
 100% [==============================]
-- Configuration --
Concurrency: 100	Loop: 10	Timeout: 2000 ms	KeepAlive: true	TimeUnit: ms	Method: GET	URL: https://www.baidu.com/
Headers: []
Body:

-- Conclusion --
total count: 1000
success count: 981
failed count: 0
error count: 19
nature duration: 2280 ms
total cost: 195688 ms
max: 1539 ms
min: 0 ms
median: 110 ms
mean: 195 ms
standard deviation: 248.309631
throughput: 430 requests/second
```

| 字段               | 含义                                                         |
| ------------------ | ------------------------------------------------------------ |
| total count        | 总共的http请求数量，等于loop*concurrency                     |
| success count      | 总共成功的http请求数量                                       |
| failed count       | 失败的请求数量，不同于 error count ，只有被【断言】校验不通过的才算失败 |
| error count        | 错误数量，一般是超时或http接口不可用                         |
| nature duration    | 自然耗时（区别于下面的总体耗时）                             |
| total cost         | 总体耗时，每个请求的耗时加总起来的总耗时，在并发情况下会大于自然耗时 |
| max                | 最大单次请求耗时                                             |
| min                | 最小单次请求耗时                                             |
| median             | 单次请求耗时中位数                                           |
| mean               | 平均每次请求耗时                                             |
| standard deviation | 每次请求耗时标准差                                           |
| throughput         | 吞吐量，数值等于 success_count / nature_duration             |

### 启动示例服务

为了更好地测试各种Assertion表达式，你可以通过以下命令启动一个Restful风格的API服务：

```shell
httptester serve 
#httptester serve --
```

此示例服务是一个典型的Restful风格的API服务，包含以下操作（以curl命令为例）：

#### 列出全部用户

```shell
curl http://localhost:1234/users
```

#### 增加一个新用户

```shell
curl -X POST \
 -d '{"name":"NewUser","age":18,"stature":175,"weight":60,"available":true}' \
 'http://localhost:1234/users'
```

#### 更新一个已有用户

```shell
curl -X PUT \
 -d '{"name":"NewUser","age":18,"stature":175,"weight":60,"available":true}' \
 'http://localhost:1234/users/{id}'
```

注意将`{id}`改为实际的用户id

#### 删除一个已有用户

 ```shell
curl -X DELETE 'http://localhost:1234/users/{id}
 ```

### Assertion

在默认情况下，httptester并不会做断言检测，也就是说只要http请求得到了响应，不论其返回的响应是什么，不论响应码是什么，都会按照成功来计算。但很多情况下，你需要判断其结果是否正确。

当前版本的httptester支持3种断言，即 **响应码断言** / **JSON断言** / **正则表达式断言**。

以上一节Restful-API服务接口中的“列出所有用户”为例，`curl`格式如下：

```shell
curl -i http://localhost:1234/users
```

返回结果：

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 13 Mar 2021 13:28:48 GMT
Content-Length: 459

[
  {
    "age" : 18,
    "id" : "732e930c-59ce-4087-b509-6288b5d2d6c5",
    "stature" : 180,
    "weight" : 62.5,
    "name" : "Jack",
    "available" : true
  },
  {
    "age" : 25,
    "id" : "30160f4c-a4ee-420b-ba61-fb0f78a3e312",
    "stature" : 175,
    "weight" : 60.5,
    "name" : "Mary",
    "available" : true
  },
  {
    "age" : 32,
    "id" : "04c7aca2-c27b-4ff1-8756-1bdf57ab7993",
    "stature" : 185,
    "weight" : 65,
    "name" : "Benjamin",
    "available" : true
  },
  {
    "age" : 15,
    "id" : "51784d48-4abd-4f01-83df-4cd43869027c",
    "stature" : 160,
    "weight" : 50.799999999999997,
    "name" : "Lee",
    "available" : false
  }
]
```

接下来我们来看在`httptester`中如何来写断言。

#### 响应码断言

使用`--assert-status-codes`来设定响应码断言，下面的例子表示，只有当http响应码为`200`或`201`时才算请求成功

```shell
httptester run -u 'http://localhost:1234/users' \
  -c 100 -l 100 \
  --assert-status-codes '200 201' 
```

#### JSON断言

此次我们使用“新增一个用户”为例，它的返回值如下：

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 12 Mar 2021 07:01:50 GMT
Content-Length: 116

{"id":"ddb6b9fe-f0af-40ef-8b44-90e5b150b3ac","name":"NewUser","age":18,"stature":175,"weight":60.5,"available":true}
```

假设我们认为返回值中的`name`值要等于`NewUser`，那么我可以使用`--assert-json-expression`来达到这一目的。注意双等号两侧的空格是必须的。

```shell
httptester run --method POST -u 'http://localhost:1234/users' \
  -b '{"name":"NewUser","age":18,"stature":175,"weight":60,"available":true}' \
  -c 100 -l 100 \
  --assert-json-expression '$.name == NewUser'
```

### 更为复杂的例子

添加`header` 设定`method` 添加`body` 添加`timeout`

```shell
httptester run --method 'POST' -u 'http://localhost:1234/users' \
  -H 'Content-Type:application/json' \
  -H 'accept:application/json' \
  -b '{"name":"NewUser","age":18,"stature":175,"weight":60,"available":true}' \
  --timeout 2s \
  --loop 100 \
  --concurrency 100 \
  --assert-status-codes '200 201' \
  --assert-json-expression '$.name == NewUser'
  -e
```

`-e`表示如果出现失败或报错，将错误信息打印出来



---
如果对这个小工具感兴趣，欢迎给我点赞。
如果有任何问题或建议，也欢迎给在此项目中给我提issue或者pull request。
