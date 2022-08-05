package server

import (
	"context"
	"errors"
	"linksvr/pkg/proto/linkpb"
)

var (
	ErrInvalidMetaData       = errors.New("request metadata is invalid, please check if parameter is right")
	ErrVerisonRecordNotFound = errors.New("version record is not find")
)

func (s *LinkServer) RegisterOSD(ctx context.Context, in *linkpb.RegisterOSDRequest) (*linkpb.RegisterOSDReply, error) {
	if err := s.selector.AddOSDService(in.Addr, in.OsdId); err != nil {
		return nil, err
	}
	return &linkpb.RegisterOSDReply{
		Result: "ok",
	}, nil
}
