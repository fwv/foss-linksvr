package ssclient

import (
	"context"
	"errors"
	"io"
	"linksvr/internal/pkg/config"
	"linksvr/pkg/zlog"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"go.uber.org/zap"
)

var ()

type OsdSerivce struct {
	osdClient osdpb.OsdServiceClient
}

func NewOsdService(osdClient osdpb.OsdServiceClient) *OsdSerivce {
	return &OsdSerivce{
		osdClient: osdClient,
	}
}

func (s *OsdSerivce) UploadFileFromStream(ctx context.Context, src io.Reader, objectName string, bucketID string) error {
	if s.osdClient == nil {
		return errors.New("osd client is nil")
	}
	stream, err := s.osdClient.UploadFile(ctx)
	if err != nil {
		return err
	}
	// send metadata
	stream.Send(&osdpb.FileUploadRequest{
		MetaData: &osdpb.MetaData{
			Name:     objectName,
			BucketId: bucketID,
		},
	})
	totalSize := 0
	chunk := make([]byte, *config.UPLOAD_CHUNK_SIZE)
	for {
		n, err := src.Read(chunk)
		// zlog.Debug("read chunk data from http request body", zap.Int("read size", n))
		// send trunk data
		stream.Send(&osdpb.FileUploadRequest{
			Chunk: chunk[:n],
		})
		if err != nil {
			if err == io.EOF {
				zlog.Info("read http request body completed", zap.String("object name", objectName), zap.Int("total size", totalSize))
				_, err := stream.CloseAndRecv()
				if err != nil {
					return err
				}
				break
			}
			zlog.Error("failed to read request body", zap.Error(err))
			return err
		}
		totalSize += n
	}
	return nil
}

func (s *OsdSerivce) DownloadFileFromStream(ctx context.Context, dst io.Writer, objectName string, bucketID string, version int64) error {
	if s.osdClient == nil {
		return errors.New("osd client is nil")
	}
	stream, err := s.osdClient.DownloadFile(ctx, &osdpb.FileDownloadRequest{
		MetaData: &osdpb.MetaData{
			Name:     objectName,
			BucketId: bucketID,
			Version:  version,
		},
	})
	if err != nil {
		return err
	}
	totalSize := 0
	for {
		rsp, err := stream.Recv()
		if rsp != nil && len(rsp.Chunk) > 0 {
			zlog.Debug("read chunk data from osd stream", zap.Int("chunk size", len(rsp.Chunk)))
			dst.Write(rsp.Chunk)
			totalSize += len(rsp.Chunk)
		}
		if err == io.EOF {
			zlog.Info("read object data completed from stream completed", zap.String("object name", objectName), zap.Int("total size", totalSize))
			break
		} else if err != nil {
			zlog.Error("failed to recv data from stream", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *OsdSerivce) UploadFile(ctx context.Context, data []byte) error {
	maxLen := len(data)

	stream, err := s.osdClient.UploadFile(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < maxLen; i += *config.UPLOAD_CHUNK_SIZE {
		head := i
		tail := i + *config.UPLOAD_CHUNK_SIZE
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
