# 红包雨开发文档
## 项目内容
开发一个春节红包雨后端系统。
项目测试接口如下：

抢红包：http://221.194.149.43:80/snatch

拆红包：http://221.194.149.43:80/open

查看红包列表：http://221.194.149.43:80/wallet_list

## 项目接口说明
本项目实现了三个接口：抢红包接口、拆红包接口、钱包列表接口。
- 抢红包接口
  - 每个用户最多只能抢到N个红包，次数可以配置
  - 用户有一定几率可以抢到红包，概率可配置
``` JSON
  Request: POST ../snatch
  {
    "uid": 123
  }
   
  Response:
  {
    "code": 0,
    "msg": "success",
    "data": {
      "envelope_id": 123,
      "max_count": 5,
      "cur_count": 3
    }
  }
  
  {
    "code": 2,
    "msg": "用户ID不存在"
  }
  
  {
    "code": 3,
    "msg": "您因为作弊被系统封禁！"
  }
  
  {
    "code": 4,
    "msg": "很遗憾，您运气不太好，没能抢到红包！"
  }
  
  {
    "code": 5,
    "msg": "您的可抢红包数已达上限！"
  }
   
```

- 拆红包接口
  - 拆开红包并加入钱包
  - 检测红包是否存在以及是否与用户对应

``` JSON
  Request: POST ../open
  {
    "uid": 123,
    "envelope_id": 123
  }
  
  Response:
  {
    "code": 0,
    "msg": "success",
    "data": {
      "value": 50
    }
  }
  
  {
    "code": 2,
    "message": "这个红包已经被打开了！"
  }
  
  {
    "code": 3,
    "message": "这个红包不存在或不属于您，您无权打开！"
  }
```

- 钱包列表接口
  - 列出用户拥有的所有红包
    - 未打开的红包不显示金额，已打开红包显示金额
  - 列出用户抢到的总金额

``` JSON
  Request: POST ../open
  {
    "uid": 123
  }
  
  Response:
  {
    "code": 0,
    "msg": "success",
    "data": {
      "amount": 112,
      "envelope_list": [
        {
          "envelope_id": 123,
          "value": 50,
          "opened": true,
          "snatch_time": 1634551711
        },
        {
          "envelope_id": 124,
          "opened": false,
          "snatch_time": 1634551812
        }
      ]
    }
  }
  
  {
    "code": 2,
    "msg": "用户ID不存在"
  }
```
此外，红包雨的总金额，总红包数可以进行动态配置与调整。

## 项目设计
### 项目构想
在该项目中，我们将以MySQL作为主要的数据持久化手段，并利用Redis来缓存数据，加速数据查询。为了应对高并发场景，进一步使用了消息队列RocketMQ来实现MySQL的异步写入，提高系统的响应速度，使其在高压下能有更好的表现。

<img src="/assets/arch.jpg">

### 数据库设计
<img src="/assets/database.jpg">

### 数据库选择
在本次项目中，我们的选择是MySQL。之所以这样选择，是因为，MySQL 适合中小型软件，而我们的红包雨后端系统实际上就是中小型系统。虽然MySQL 的容量略逊于 Oracle 数据库，但是项目中实际上只有两个表，所以并没有必要选择其他容量较大的数据库。此外，MySQL能够快速、有效和安全的处理大量的数据。相对于 Oracle 等数据库来说，MySQL 的使用是非常简单的。MySQL 具有快速、健壮和易用等特点。
### Redis
式存储，Key为User:{uid}，包括amount与cur_count两个字段。红包信息同样采用了Hash的方式存储，Key为Envelope:{eid}，包括了open、snatch_time、value、uid四个字段。用户红包列表采用了Set来存储，Key为User:{uid}:Envelopes。
我们在碰到需要执行耗时特别久，且结果不频繁变动的SQL，就特别适合将运行结果放入缓存。这样，后面的请求就去缓存中读取，使得请求能够迅速响应。在大并发的情况下，所有的请求直接访问数据库，数据库会出现连接异常。这个时候，就需要使用Redis做一个缓冲操作，让请求先访问到Redis，而不是直接访问数据库。考虑这个项目大并发，大访问量的情况下，使用Redis是一个很好的选择。
### 红包金额分配算法
对于红包金额，经过研究，我们认为红包金额的分布大致呈现截尾正态分布。简单来说，就是有上下限的正态分布。这样的分布能够保证用户获得金额不同的同时又不会相差很大，且可以很好的满足尽可能用完预算这一要求。该算法具体来讲，使用Box-Muller转换公式，将一个0~1之间的随机分布转化为一个高斯分布，再把高斯分布转化为一个均值为mu方差为sigma的正态分布。
每次计算红包金额时，我们会从Redis中读取剩余金额与剩余红包数的信息，并根据这两者计算出红包金额的均值，最后根据均值与设定的方差来随机生成红包金额。
### MySQL与Redis的数据同步
我们在系统中使用了Redis来缓存数据，加速查询，但这可能会带来数据不同步的问题。这种不同步通常是由于数据的更新导致的，我们应当采取一定的策略来保证MySQL与Redis中的数据的同步。我们采用了写更新这一思路来解决这一问题，当一个接口准备更新MySQL中的数据时，我们也会同步更新Redis中的缓存数据，这样就能够保证两者的同步。
### 反作弊机制设计
在红包雨这类服务的实际使用场景中，由于其接口直接对外暴露，可能会有用户使用脚本工具等作弊手段来获取不正当的利益，因此我们的系统需要找出这些用户并对其进行一定的惩罚。首先可以明确的是，作弊用户通常会出现在抢红包这一场景中，因此我们首先在该场景中加入反作弊机制。
当一名用户访问抢红包服务时，我们将它的ID信息缓存在redis中（User:{uid}:Snatch），并设置10秒的过期时间。这个Key可以一定程度上衡量一名用户在10秒内抢红包的数量。一旦某一名用户发出了抢红包请求，且我们检测到该用户所对应的key超过了一个阈值（例如10），我们就判定该用户作弊，并阻止他的抢红包行为。进一步地，我们将在Redis中设置一个Key来标识作弊用户（User:{uid}:Cheat），在一定时间段中拒绝用户的抢红包请求。
### 高并发情况下的更新冲突解决
我们知道，在并发的情况下，许多在非并发下不存在的问题也会随之而来，最典型的就是场景就是并发地更新同一个值。在红包雨项目中，如果用户正常使用系统，则很难出现同一个值的并发更新，而如果用户采取了作弊手段，例如使用脚本来抢红包，则有可能会导致并发更新同一个值（例如用户已抢到的红包数），最终导致更新错误（系统中记录的已抢红包数小于实际分配给该用户的红包数），系统出现不一致性，为我们的业务带来影响。
以更新用户已抢到的红包数CurCount这一值为例，在Redis中，我们可以使用Redis提供的Incr操作来完成这个任务，这是一个原子操作，因此可以很好的解决我们的问题。而在更新了Redis中的值后，我们也需要再更新MySQL中的值，这里我选择使用GORM提供的Lock来为更新的数据上锁，同时使用了GORM提供的Transaction来保证操作的原子性。通过这些保护工作，我确保了每分配出一个红包，对应的用户在Redis与MySQL中的Curcount都能同步且正确地变化。
针对打开红包这一功能，我也使用了如上的方法来避免读与写的冲突。
### 限流策略
我们在系统中采取了令牌桶的方式进行限流，在启动服务前为router添加一个令牌桶中间件。
``` Go
r.Use(tokenbucket.NewLimiter(2000, 2000, 500*time.Millisecond))
```
我们采用的令牌桶总容量为2000块令牌，每秒产生2000块令牌，若某一请求到达时没有令牌则等待500ms。
### 压测情况
我们使用了wrk进行了open接口的压力测试，具体结果如下：
``` Bash
root@node-mvph6z:~/wrk# wrk -c600 -d60s -t6 --latency -s snatch.lua http://172.16.237.193:80/snatch
Running 1m test @ http://172.16.237.193:80/snatch
  6 threads and 600 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   142.53ms   90.98ms 480.80ms   77.80%
    Req/Sec   756.91     90.60     1.17k    68.53%
  Latency Distribution
     50%  103.83ms
     75%  191.22ms
     90%  327.49ms
     99%  361.64ms
  271363 requests in 1.00m, 49.62MB read
Requests/sec:   4517.21
Transfer/sec:    845.82KB

root@node-mvph6z:~/wrk# wrk -c400 -d60s -t6 --latency -s snatch.lua http://172.16.237.193:80/get_wallet_list
Running 1m test @ http://172.16.237.193:80/get_wallet_list
  6 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   194.42ms   23.31ms 588.03ms   96.68%
    Req/Sec   340.41     57.00     0.91k    89.08%
  Latency Distribution
     50%  197.93ms
     75%  198.30ms
     90%  198.69ms
     99%  203.05ms
  122097 requests in 1.00m, 75.37MB read
Requests/sec:   2032.78
Transfer/sec:      1.25MB
```

``` lua
-- snatch.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

math.randomseed(os.time())
request = function()
    uid = math.random(100000)
    body = "uid=" .. uid
    return wrk.format(nil, nil, nil, body)
 end
```
在测试的过程中，平均响应延迟为88.46ms，标准差为23.37ms，最大延迟为285.75ms，有82.49%的数据位于正负一个标准差的范围内。在该测试条件下，系统的QPS为2091.84，且经过后期检查，系统在该过程中成功地保持了业务逻辑的正确性与数据的一致性。

### 项目的可配置参数
|Key|说明|
|---|----|
|TotalMoney|当前剩余的总金额|
|EnvelopeNum|当前还准备发放的红包数|
|MaxCount|当前每个用户可以抢到的红包数|
|EnvelopeId|当前当前已分配的最大红包ID，用于实现红包ID的自增分配|

### 请求错误码规定
|错误码\Handler|OpenHandler|SnatchHandler|WalletListHandler|
|---|---|---|---|
|0|正常|正常|正常|
|1|系统内部错误|系统内部错误|系统内部错误|
|2|红包已被打开|用户不存在|用户不存在|
|3|红包不存在或不属于该用户|用户由于作弊被封禁||
|4||手气不佳，没能抢到红包||
|5||用户红包数量已达上限||

### 项目结构说明
项目的主要目录结构如下，其中allocate中实现了红包分配算法，dao中实现了MySQL的CRUD操作的封装，handler中实现了三个接口的处理函数以及RocketMQ、Redis的初始化，rocketmq中实现了RocketMQ中的消费者，负责接收消息并向MySQL中写入数据，tokenbucket中实现了令牌桶。
```
├── Dockerfile
├── allocate
│   └── money.go
├── config.yaml
├── dao
│   ├── ConnectionInit.go
│   ├── EnvelopeCrud.go
│   └── UserCrud.go
├── devops-deploym.yaml
├── go.mod
├── go.sum
├── handler
│   ├── ConnectHandler.go
│   ├── HandlerHelper.go
│   ├── OpenHandler.go
│   ├── PingTest.go
│   ├── RocketMQClient.go
│   ├── SnatchHandler.go
│   └── WalletListHandler.go
├── main.go
├── rocketmq
│   ├── Dockerfile
│   ├── main
│   └── main.go
└── tokenbucket
    └── limit.go
```