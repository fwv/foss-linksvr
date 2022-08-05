package object

import (
	"context"
	"linksvr/internal/pkg/bucket"
	"linksvr/internal/pkg/osd"
	"linksvr/pkg/zlog"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Handler /object/{obejectName}
// curl -v 127.0.0.1:4000/object/mysh1 --upload-file send.sh -H "Digest:SHA-256=vqQo/S8aBTe+q3Qc9b+L1UruvKpPV5TA1+Ou20REGHQ=" -H "Authentication:fengwei"
func Handler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authentication")
	if auth == "" {
		// w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m := r.Method
	if m == http.MethodGet {
		get(w, r)
	}

	if m == http.MethodPut {
		put(w, r)
	}

	if m == http.MethodDelete {
		zlog.Info("handler object Delete")
		w.WriteHeader(http.StatusOK)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	zlog.Info("handler object Get")
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	zlog.Info("start get object", zap.Any("obeject name", objectName))

	// parse verison
	openID := r.Header.Get("Authentication")
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	osdSvc, err := osd.Selector.ChooseOSDbyHash(bucketName)
	if err != nil {
		zlog.Error("choose osd server failed", zap.Any("bucket", bucketName), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ds := bucket.GetMongoBucketDataSource()
	// check bucketID
	b, err := ds.FindBucket(context.TODO(), openID, bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if b == nil {
		zlog.Error("user's bucket not find", zap.Any("open id", openID), zap.Any("bucket name", bucketName))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can't find this bucket, please check if the bucket id you input is right."))
		return
	}
	if b.BucketID == "" {
		zlog.Error("bucketID is nil", zap.Any("open id", openID), zap.Any("bucket name", bucketName))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	versionParam := r.URL.Query().Get("version")
	var version int64
	if versionParam == "" {
		version = 0
	} else {
		version, err = strconv.ParseInt(versionParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	// bucketFlag := strings.Join([]string{bucketName, "-", b.BucketID}, "")
	err = osdSvc.DownloadFileFromStream(context.Background(), w, objectName, bucketName, version)
	if err != nil {
		zlog.Error("failed to download file from osdsvc")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func put(w http.ResponseWriter, r *http.Request) {
	zlog.Info("handler object Put")
	// openID := r.Header.Get("Authentication")
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	// zlog.Info("start put object", zap.Any("obeject name", objectName))
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" || objectName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check bucketID
	// ds := bucket.GetMongoBucketDataSource()
	// b, _ := ds.FindBucket(context.TODO(), openID, bucketName)
	// if b == nil {
	// 	zlog.Error("user's bucket not find", zap.Any("open id", openID), zap.Any("bucket name", bucketName))
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("can't find this bucket, please check if the bucket id you input is right."))
	// 	return
	// }

	// osdSvc, _ := osd.Selector.ChooseOSDbyHash(bucketName)
	osdSvc, _ := osd.Selector.OsdClients[1]
	// if err != nil {
	// 	zlog.Error("choose osd server failed", zap.Any("bucket", bucketName), zap.Error(err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// use hash to identify object
	// err = osdSvc.UploadFileFromStream(context.Background(), r.Body, url.PathEscape(hash))
	// bucketFlag := strings.Join([]string{bucketName, "-", b.BucketID}, "")
	if err := osdSvc.UploadFileFromStream(context.Background(), r.Body, objectName, bucketName); err != nil {
		zlog.Error("", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
