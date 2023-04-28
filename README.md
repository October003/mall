## go-zero RPC 服务


### 编写RPC

1. 编写pb文件，并生成代码
2. 完善配置结构体和配置文件(结构体和yaml文件一定要一致) 
3. 完善ServiceContext (生成代码 和 业务代码 之间的一座桥梁)
4. 完善rpc的业务逻辑

### RPC 服务测试工具

一个测试grpc服务的ui工具
https://github.com/fullstorydev/grpcui

安装:
```bash
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

其中 `localhost:8080` 是你rpc服务的地址
```bash
grpcui -plaintext localhost:8080
```

#### 如果出现下面这种情况:
```bash
grpcui -plaintext localhost:8080
Failed to compute set of methods to expose: server does not support the reflection API 
```
要想用 grpcui 测试RPC服务，需要让go-zero rpc服务以 dev或text
需要在配置文件中指定mode

```yaml
Name: user.rpc
Mode: dev
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: user.rpc
Mysql:
  DataSource: 
CacheRedis:
  - Host: 127.0.0.1:6379
```

### 订单服务的检索接口

/api/order/search :: 根据订单id查询订单信息
 - RPC  -->  userID  -->  user.GetUser

### go-zero 中通过RPC调用其他服务

1. 配置RPC客户端(配置结构体和yaml配置文件都要加RPC客户端配置，注意: ETCD 的 key要对应上)
2. 修改 ServiceContext (告诉生成的代码 我现在有RPC的客户端了)
 - go-zero 中的RPC服务会自动生成一份客户端代码
3. 编写业务逻辑 (可以直接通过RPC客户端发起RPC调用)


### RPC 调用传递metadata

#### 几个知识点

1. 什么是metadata？ 什么样的数据应该存入metadata？ 它和请求参数有什么区别
   元数据（metadata）是指在处理RPC请求和响应过程中需要但又不属于具体业务（例如身份验证详细信息）的信息，采用键值对列表的形式，其中键是string类型，值通常是[]string类型，但也可以是二进制数据。gRPC中的 metadata 类似于我们在 HTTP headers中的键值对，元数据可以包含认证token、请求标识和监控标签等。

2. gRPC的拦截器: 客户端的拦截器和服务器端的拦截器

#### go-zero 项目中添加client端拦截器

order服务的search接口中添加拦截器，添加一些userID、requestID、token等数据

#### go-zero 项目中添加server端拦截器



