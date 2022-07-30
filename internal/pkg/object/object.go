package object

import (
	"context"
	"linksvr/internal/pkg/ssclient"
	"linksvr/pkg/zlog"
	"log"
	"net/http"
	"strings"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Handler /object/{obejectName}
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		zlog.Info("handler object Get")
		w.WriteHeader(http.StatusOK)
	}

	if m == http.MethodPut {
		put(w, r)
	}

	if m == http.MethodDelete {
		zlog.Info("handler object Delete")
		w.WriteHeader(http.StatusOK)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	zlog.Info("handler object Put")
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	zlog.Info("start put object", zap.Any("obeject name", objectName))

	// todo: connection reuse
	conn, err := grpc.Dial(":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := osdpb.NewOsdServiceClient(conn)
	osdSvc := ssclient.NewOsdService(c)

	err = osdSvc.UploadFileFromStream(context.Background(), r.Body, objectName)
	if err != nil {
		zlog.Error("", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
