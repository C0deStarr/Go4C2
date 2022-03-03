package main

import (
	"Go4C2/grpcapi"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	var (
		conn          *grpc.ClientConn
		err           error
		adminClient   grpcapi.AdminClient
		strFlagServer string
		nFlagPort     uint
		strFlagCmd    string
	)
	flag.StringVar(&strFlagServer, "server", "", "team server ip")
	flag.UintVar(&nFlagPort, "port", 0, "team server port")
	flag.StringVar(&strFlagCmd, "cmd", "", "\"CMD [ARG,ARG,...]\"")
	flag.Parse()
	if len(os.Args) <= 1 {
		log.Fatalf("-h to see usage")
	}
	// 1. connect to the team server
	server := fmt.Sprintf("%s:%d", strFlagServer, nFlagPort)
	conn, err = grpc.Dial(server, grpc.WithInsecure())
	if nil != err {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	log.Printf("Dial() ok")

	// 2. init client
	adminClient = grpcapi.NewAdminClient(conn)
	log.Printf("NewAdminClient() ok")

	// 3. send grpcCmd
	grpcCmd := new(grpcapi.Command)
	grpcCmd.In = strFlagCmd
	var cmdResult *grpcapi.Command
	ctx := context.Background()
	log.Printf("SendCommand()")
	cmdResult, err = adminClient.SendCommand(ctx, grpcCmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmdResult.Out)
}
