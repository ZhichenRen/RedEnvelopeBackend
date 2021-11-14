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
```