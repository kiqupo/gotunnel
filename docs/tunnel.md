# 内网穿透通道

## 研究背景
有多个内网服务应用与一个外网应用，需要通过外网应用开放多个内网服务给外网使用，需要集成内网穿透通道代码。
### 要求
1. 支持TCP、HTTP、WS请求
2. 支持高并发

## 基本架构
1. gotunnel服务端（公网服务）：
> 1. 监听对外端口，把请求数据转发到通道
> 2. 监听通道端口，与客户端建立数据通道
2. gotunnel客户端（内网服务）：
> 1. 连接数据通道，建立数据通道
> 2. 转发数据至服务端口

## 流程图
![时序图](images/sequence.png)

## 主要机制
### 端口分配：
gotunnel客户端启动时发送附带配置密钥的HTTP请求通道接口，gotunnel服务端分配出不重复的外网客户端通道端口、与内网服务通道端口，并且返回配置的连接数N<br/>
![端口分配图](images/port.png)

### 多路复用/连接创建：
1. 背景：当客户端频繁使用短连接请求时，频繁创建与销毁TCP连接开销很大导致性能降低，所以提出对gotunnel服务端与客户端之间TCP连接复用机制。
2. 选型 [yamux](https://github.com/hashicorp/yamux) 包来实现
   1. 原理：把一个TCP连接(session)虚化成多个异步数据流(stream)<br/>
   ![多路复用原理图](images/yamux.png)
   2. yamux自带心跳检测
3. 机制：把N个TCP连接虚拟成N个session，接收到用户请求后再动态创建stream连接。<br/>
   ![连接创建逻辑图](images/newconnect.png)

## 测试结果：
1. websocket测试：
![ws测试结果图](images/wstest.png)
2. HTTP请求测试
![HTTP测试结果图1](images/testresult.png)
![HTTP测试结果图2](images/testresult2.png)

## TODO：
1. 安全性考虑：服务端可加上socks5认证等