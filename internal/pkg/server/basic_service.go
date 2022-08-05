package server

import (
	"linksvr/internal/pkg/config"
	"linksvr/internal/pkg/osd"
	"linksvr/pkg/proto/linkpb"
	"linksvr/pkg/zlog"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type LinkServer struct {
	linkpb.UnimplementedLinkServiceServer
	selector *osd.OSDSelector
}

func NewLinkServer(selector *osd.OSDSelector) *LinkServer {
	s := &LinkServer{
		selector: selector,
	}
	return s
}

func (s *LinkServer) Serve() error {
	// register Link
	lis, err := net.Listen("tcp", *config.LINK_GRPC_ADDR)
	if err != nil {
		zlog.Fatal("failed to listen.", zap.Error(err))
	}
	server := grpc.NewServer()
	linkpb.RegisterLinkServiceServer(server, s)
	zlog.Info("server start listening", zap.String("addr", *config.LINK_GRPC_ADDR))
	if err := server.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
		return err
	}
	return nil
}

func (s *LinkServer) Shutdown() {

}
