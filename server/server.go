package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type beaconServer struct {
}

func NewbeaconServer() *beaconServer {
	newbeaconServer := new(beaconServer)
	return newbeaconServer
}

func (beaconServer *beaconServer) FetchCommand(context context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	cmd := new(grpcapi.Command)
	cmd.Out = "This is FetchCommand"
	log.Printf("recv cmd from server")
	return cmd, nil
}
func (beaconServer *beaconServer) SendOutput(context context.Context, empty *grpcapi.Command) (*grpcapi.Empty, error) {
	return &grpcapi.Empty{}, nil
}

var (
	g_nbeaconPort = 4444
)

func main() {
	var ()

	// register beacon server
	beaconServer := NewbeaconServer()
	grpcbeaconServer := grpc.NewServer()
	grpcapi.RegisterBeaconServer(grpcbeaconServer, beaconServer)

	// listen
	strbeaconAddr := fmt.Sprintf("localhost:%d", g_nbeaconPort)
	beaconListener, err := net.Listen("tcp", strbeaconAddr)
	if nil != err {
		log.Fatal(err)
	}
	grpcbeaconServer.Serve(beaconListener)
}
