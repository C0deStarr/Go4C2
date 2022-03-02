package main

import (
	"Go4C2/grpcapi"
	"context"
	"fmt"
	"log"
	"net"

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

	cmd := <-beaconServer.m_chanCmd

	log.Printf("cmd sent: %s", cmd.In)
	return cmd, nil
}
func (beaconServer *BeaconServer) SendResult(context context.Context, cmdResult *grpcapi.Command) (*grpcapi.Empty, error) {
	log.Printf("SendResult()")
	log.Printf(cmdResult.Out)
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
		adminServer.m_chanCmd <- cmd
	}()
	res = <-adminServer.m_chanResult
	return res, nil
}

var (
	g_nBeaconPort = 4444
	g_nAdminPort  = 5555
)

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
	)

	// 0. init channel
	chanCmd = make(chan *grpcapi.Command)    // no buffer to block the client input goroutine
	chanResult = make(chan *grpcapi.Command) // no buffer to block the client input goroutine

	// 1. Beacon server
	// 1.1 register beacon server
	beaconServer = NewBeaconServer(chanCmd, chanResult)
	grpcBeaconServer = grpc.NewServer()
	grpcapi.RegisterBeaconServer(grpcBeaconServer, beaconServer)

	// 1.2 listen
	strBeaconAddr := fmt.Sprintf("localhost:%d", g_nBeaconPort)
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
	strAdminAddr := fmt.Sprintf("localhost:%d", g_nAdminPort)
	adminListener, err = net.Listen("tcp", strAdminAddr)
	if nil != err {
		log.Fatal(err)
	}
	// 2.3 run admin server
	grpcAdminServer.Serve(adminListener)
}
