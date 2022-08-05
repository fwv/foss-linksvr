package main

import (
	"flag"
	"linksvr/internal/pkg/bucket"
	"linksvr/internal/pkg/config"
	"linksvr/internal/pkg/object"
	"linksvr/internal/pkg/osd"
	"linksvr/internal/pkg/server"
	"linksvr/internal/pkg/version"
	"linksvr/pkg/zlog"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	flag.Parse()
	done := make(chan bool)
	// http server
	http.HandleFunc("/object/", object.Handler)
	http.HandleFunc("/bucket/", bucket.Handler)
	http.HandleFunc("/version/", version.Handler)
	zlog.Info("linksvr http server start to serving", zap.String("addr", *config.LINK_HTTP_ADD))
	go func() {
		http.ListenAndServe(*config.LINK_HTTP_ADD, nil)
	}()

	// grpc server
	linkServer := server.NewLinkServer(osd.Selector)
	go func() {
		if err := linkServer.Serve(); err != nil {
			zlog.Error("linksvr grpc serve failed", zap.Error(err))
			close(done)
		}
		zlog.Info("linksvr grpc server start to serving", zap.String("addr", *config.LINK_GRPC_ADDR))
	}()

	// waiting for osd service register
	registerCh := make(chan bool)
	osd.Selector.Init(registerCh)
	<-registerCh
	<-done
}
