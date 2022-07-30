package osd

import (
	"linksvr/internal/pkg/ssclient"
	"log"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OSDStream struct {
	osdService *ssclient.OsdSerivce
}

func NewOSDStream(osdAddr string) (*OSDStream, error) {
	conn, err := grpc.Dial(osdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	// defer conn.Close()
	c := osdpb.NewOsdServiceClient(conn)
	osdSvc := ssclient.NewOsdService(c)
	return &OSDStream{
		osdService: osdSvc,
	}, nil
}

// func (s *OSDStream) Write(p []byte) (int, error) {
// 	n := 0
// 	// 遍历输入数据，按字节写入目标资源
// 	for _, b := range p {
// 		w.ch <- b
// 		n++
// 	}
// 	return n, nil
// }
