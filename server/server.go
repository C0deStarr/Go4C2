package main

import (
	"Go4C2/grpcapi"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

////////////////////////////
// BeaconServer
////////////////////////////
type BeaconServer struct {
	m_chanCmd    chan *grpcapi.Command // recv cmd from admin server
	m_chanResult chan *grpcapi.Command // send result to admin server
}

func NewBeaconServer(_chanCmd, _chanResult chan *grpcapi.Command) *BeaconServer {
	newBeaconServer := new(BeaconServer)
	newBeaconServer.m_chanCmd = _chanCmd
	newBeaconServer.m_chanResult = _chanResult
	return newBeaconServer
}

func (beaconServer *BeaconServer) FetchCommand(context context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	log.Printf("FetchCommand()")
	cmd := new(grpcapi.Command)
	select {
	case cmd, ok := <-beaconServer.m_chanCmd:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("channel cmd closed")
	default:
		// no jobs
		return cmd, nil
	}

	log.Printf("cmd sent: %s", cmd.In)
	return cmd, nil
}
func (beaconServer *BeaconServer) SendResult(context context.Context, cmdResult *grpcapi.Command) (*grpcapi.Empty, error) {
	log.Printf("SendResult()")
	beaconServer.m_chanResult <- cmdResult
	return &grpcapi.Empty{}, nil
}

////////////////////////////
// AdminServer
////////////////////////////
type AdminServer struct {
	m_chanCmd    chan *grpcapi.Command // send cmd to beacon server
	m_chanResult chan *grpcapi.Command // recv result from beacon server
}

func NewAdminServer(_chanCmd, _chanResult chan *grpcapi.Command) *AdminServer {
	newAdminServer := new(AdminServer)
	newAdminServer.m_chanCmd = _chanCmd
	newAdminServer.m_chanResult = _chanResult
	return newAdminServer
}
func (adminServer *AdminServer) SendCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	log.Printf("SendCommand()")
	var res *grpcapi.Command
	go func() {
		log.Printf("wait for cmd")
		adminServer.m_chanCmd <- cmd
		log.Printf("cmd sent")
	}()
	log.Printf("wait for result")
	res = <-adminServer.m_chanResult
	log.Printf("result sent")
	return res, nil
}

func main() {
	var (
		beaconServer     *BeaconServer
		grpcBeaconServer *grpc.Server
		beaconListener   net.Listener

		adminServer     *AdminServer
		grpcAdminServer *grpc.Server
		adminListener   net.Listener

		chanCmd    chan *grpcapi.Command
		chanResult chan *grpcapi.Command

		err error

		strFlagServer   string
		nFlagBeaconPort uint
		nFlagClientPort uint
	)
	flag.StringVar(&strFlagServer, "server", "", "team server ip")
	flag.UintVar(&nFlagBeaconPort, "beaconport", 0, "team server port for beacon")
	flag.UintVar(&nFlagClientPort, "clientport", 0, "team server port for client")
	flag.Parse()
	if len(os.Args) <= 1 {
		log.Fatalf("-h to see usage")
	}
	// 0. init channel
	chanCmd = make(chan *grpcapi.Command)    // no buffer to block the client input goroutine
	chanResult = make(chan *grpcapi.Command) // no buffer to block the client input goroutine

	// 1. Beacon server
	// 1.1 register beacon server
	beaconServer = NewBeaconServer(chanCmd, chanResult)
	grpcBeaconServer = grpc.NewServer()
	grpcapi.RegisterBeaconServer(grpcBeaconServer, beaconServer)

	// 1.2 listen
	strBeaconAddr := fmt.Sprintf("%s:%d", strFlagServer, nFlagBeaconPort)
	beaconListener, err = net.Listen("tcp", strBeaconAddr)
	if nil != err {
		log.Fatal(err)
	}

	// 1.3 run beacon server
	go func() {
		grpcBeaconServer.Serve(beaconListener)
	}()

	// 2. Admin server
	// 2.1 register admin server
	adminServer = NewAdminServer(chanCmd, chanResult)
	grpcAdminServer = grpc.NewServer()
	grpcapi.RegisterAdminServer(grpcAdminServer, adminServer)

	// 2.2 listen
	strAdminAddr := fmt.Sprintf("%s:%d", strFlagServer, nFlagClientPort)
	adminListener, err = net.Listen("tcp", strAdminAddr)
	if nil != err {
		log.Fatal(err)
	}
	// 2.3 run admin server
	grpcAdminServer.Serve(adminListener)
}
