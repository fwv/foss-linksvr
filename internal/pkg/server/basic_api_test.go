package server

import (
	"context"
	"flag"
	"linksvr/pkg/zlog"
	"log"
	"testing"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	osdAddr = flag.String("osdAddr", ":50051", "osd serivce address")
)

func BenchmarkUploadFile(b *testing.B) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*osdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := osdpb.NewOsdServiceClient(conn)

	str := "hello,1231231231231231231231231231231231123123123"
	data := []byte(str)
	for i := 0; i < b.N; i++ {
		UploadFile(context.Background(), c, data)
	}
}

func TestUploadFile(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*osdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := osdpb.NewOsdServiceClient(conn)

	str := "hello,1231231231231231231231231231231231123123123"
	data := []byte(str)
	UploadFile(context.Background(), c, data)
}

func UploadFile(ctx context.Context, c osdpb.OsdServiceClient, data []byte) error {
	maxLen := len(data)
	maxChunkSize := 1024

	stream, err := c.UploadFile(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < maxLen; i += maxChunkSize {
		head := i
		tail := i + maxChunkSize
		if tail > maxLen {
			tail = maxLen
		}
		stream.Send(&osdpb.FileUploadRequest{
			MetaData: &osdpb.MetaData{},
			Chunk:    data[head:tail],
		})
	}
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	zlog.Info("upload completed", zap.Any("reslut code", rsp.ResultCode))
	return nil
}
