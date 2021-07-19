# 数据通道demo

## 测试端口：
1. 用户访问端口：8007
2. SC控制SA端口：8009
3. 双向隧道端口：8008

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
    |-- utils
    |   |-- network.go  --工具类
    |-- README.md
```