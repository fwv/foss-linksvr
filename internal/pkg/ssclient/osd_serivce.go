package ssclient

import (
	"context"
	"errors"
	"io"
	"linksvr/pkg/zlog"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"go.uber.org/zap"
)

type OsdSerivce struct {
	osdClient osdpb.OsdServiceClient
}

func NewOsdService(osdClient osdpb.OsdServiceClient) *OsdSerivce {
	return &OsdSerivce{
		osdClient: osdClient,
	}
}

func (s *OsdSerivce) SayHello(ctx context.Context, message string) error {
	if s.osdClient == nil {
		return errors.New("osd service is empty")
	}
	rsp, err := s.osdClient.SayHello(ctx, &osdpb.HelloRequest{
		Name: message,
	})
	if err != nil {
		return err
	}
	zlog.Info("server reply", zap.String("message", rsp.Message))
	return nil
}

func (s *OsdSerivce) UploadFile(ctx context.Context, data []byte) error {
	maxLen := len(data)
	maxChunkSize := 1024

	stream, err := s.osdClient.UploadFile(ctx)
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

func (s *OsdSerivce) UploadFileFromStream(ctx context.Context, src io.Reader, objectName string) error {
	if s.osdClient == nil {
		return errors.New("osd client is nil")
	}
	stream, err := s.osdClient.UploadFile(ctx)
	if err != nil {
		return err
	}
	maxChunkSize := 1024
	chunk := make([]byte, maxChunkSize)
	for {
		n, err := src.Read(chunk)
		zlog.Info("", zap.Int("read size", n), zap.String("content", string(chunk[:n])))
		stream.Send(&osdpb.FileUploadRequest{
			MetaData: &osdpb.MetaData{
				Name: objectName,
			},
			Chunk: chunk[:n],
		})
		if err != nil {
			if err == io.EOF {
				zlog.Info("read completed")
				_, err := stream.CloseAndRecv()
				if err != nil {
					return err
				}
				break
			}
			zlog.Error("failed to read request body", zap.Error(err))
			return err
		}
	}
	return nil
}
