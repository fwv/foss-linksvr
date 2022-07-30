package main

import (
	"flag"
	"linksvr/internal/pkg/object"
	"linksvr/pkg/zlog"
	"net/http"

	"go.uber.org/zap"
)

var (
	osdAddr      = flag.String("osdAddr", ":5000", "osd serivce address")
	linkHttpAddr = flag.String("osdHttpAddr", ":4000", "linksvr http server address")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	// conn, err := grpc.Dial(*osdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := osdpb.NewOsdServiceClient(conn)
	// osdSvc := ssclient.NewOsdService(c)
	// osdSvc.SayHello(context.Background(), "fengwei")

	// // upload test
	// content := "134568901234567890akjsdhfakjsdhflak;jsdf;alsdjf;alskdjf;laksdjf;alksdjf;alksdjf;alksdjfa;lksdjalskhfuioweyhfqiuwehfuihcvzhdxcvuilkawheuiolfhalskdjfhkwjleh1234567890"
	// data := []byte(content)
	// osdSvc.UploadFile(context.Background(), data)
	http.HandleFunc("/object/", object.Handler)
	zlog.Info("osd http server start to serving", zap.String("addr", *linkHttpAddr))
	http.ListenAndServe(*linkHttpAddr, nil)
}
