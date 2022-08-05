package config

import "flag"

var (
	LINK_GRPC_ADDR    = flag.String("LINK_GRPC_ADDR", ":5001", "linksvr grpc server address")
	OSD_GRPC_ADDR     = flag.String("OSD_GRPC_ADDR", ":5000", "osdsvr grpc server address")
	LINK_HTTP_ADD     = flag.String("LINK_HTTP_ADD", ":4000", "linksvr http server address")
	UPLOAD_CHUNK_SIZE = flag.Int("UPLOAD_CHUNK_SIZE ", 1024*256, "upload chunk size")
	OSD_NUM           = flag.Int("OSD_NUM   ", 1, "osd service num")
)
