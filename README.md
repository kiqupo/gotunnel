# 数据通道demo

## 测试端口：
1. 用户访问端口：8007
2. 双向隧道端口：8008

## 测试：
>1. postman
>2. curl

## 项目目录：
```
|-- tunnelDemo
    |-- cmd
    |   |-- server
    |       |-- main.go --sc数据通道demo
    |   |-- client
    |       |-- main.go --sa通道demo
    |-- docs
    |   |-- tunnel.md   --预研文档
    |-- pkg
    |   |-- tunnel
    |       |-- client.go --通道客户端工具类
    |       |-- server.go --通道服务端工具类
    |       |-- tunnel.go --通道公用工具类
    |-- README.md
```

## 快速使用
```
go mod tidy
```
server（SC端）
```
tunnel.ServerTunnel(conf)

// 模拟分配端口
tunnel.RegisterController("213", tunnelPost, visitPost)
```
client（SA端）
```
tunnel.ClientRun(conf)
```

## 性能监控地址
> server: 127.0.0.1:6060

> client: 127.0.0.1:6061

## TODO
