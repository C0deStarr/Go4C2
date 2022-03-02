package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

var (
	g_strTeamServer    = "localhost"
	g_nAdminServerPort = 5555
)

func main() {
	var (
		conn        *grpc.ClientConn
		err         error
		adminClient grpcapi.AdminClient
	)

	// 1. connect to the team server
	server := fmt.Sprintf("%s:%d", g_strTeamServer, g_nAdminServerPort)
	conn, err = grpc.Dial(server, grpc.WithInsecure())
	if nil != err {
		log.Fatalf("grpc.Dial error: %v", err)
	}

	// 2. init client
	adminClient = grpcapi.NewAdminClient(conn)

	// 3. send grpcCmd
	grpcCmd := new(grpcapi.Command)
	grpcCmd.In = "whoami"
	var cmdResult *grpcapi.Command
	ctx := context.Background()
	cmdResult, err = adminClient.SendCommand(ctx, grpcCmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmdResult.Out)
}
