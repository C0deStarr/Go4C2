package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

var (
	g_strTeamServer     = "localhost"
	g_nBeaconServerPort = 4444
)

func main() {
	var (
		conn         *grpc.ClientConn
		err          error
		beaconClient grpcapi.BeaconClient
		cmd          *grpcapi.Command
	)
	server := fmt.Sprintf("%s:%d", g_strTeamServer, g_nBeaconServerPort)
	conn, err = grpc.Dial(server, grpc.WithInsecure())
	if nil != err {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	beaconClient = grpcapi.NewBeaconClient(conn)

	// begin polling
	ctx := context.Background()
	req := new(grpcapi.Empty)
	cmd, err = beaconClient.FetchCommand(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("response: %s", cmd.Out)
}
