# 0. Overview
A C2 based on Go。

File tree：

* server：like cs teamserver；

* beacon：copy the name of cobalt strike，responsible for receiving and executing cmd，and sending the results back；

* client：like cobalt strike 's client；

# 1. Environment

1. Go，remember to put `%GOPATH%/bin` to %PATH%；
2. [Protobuf](https://developers.google.cn/protocol-buffers)：
	* [Tutorial](https://developers.google.cn/protocol-buffers/docs/gotutorial)，
	* [Compiler](https://github.com/protocolbuffers/protobuf/releases) 
	* Go plugin，%GOPATH%/bin/protoc-gen-go.exe:

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

If there is no exe generated, build in /protoc-gen-go：

```
go build -o protoc-gen-go main.go
```


3. [gRPC](https://www.grpc.io/docs/)：

```
google.golang.org/grpc
```



# 2. Usage

## Interface

Execute this cmd in /grpapi to compile beacon.proto, then /grpaapi/beacon.pb.go will be generated：

```
protoc -I . --go_out=plugins=grpc:./ ./beacon.proto
```

There are service interfaces in beacon.pb.go，eg：

```go
// beaconServer is the server API for beacon service.
type beaconServer interface {
	FetchCommand(context.Context, *Empty) (*Command, error)
	SendOutput(context.Context, *Command) (*Empty, error)
}
```



## Launch the server

Two ports need to be specified for two servers：

* BeaconServer，connecte to beacon on the target host；
* AdminServer，connect to the client；

example:  

```shell
> go run .\server\server.go -h
-beaconport uint
        team server port for beacon
-clientport uint
        team server port for client
-server string
        team server ip
> go run .\server\server.go -server localhost -clientport 4444 -beaconport 5555
```

## Beacon

```shell
> go run .\beacon.go -server localhost -port 5555
```

## Client

Only one cmd could be executed per process for now：

```shell
> go run .\client.go -server localhost -port 4444 -cmd "whoami"
```


# 3. Known issues

* Encoding troubles；

* If the results is too long，such as "ipconfig /all"，SendResult() will fail；