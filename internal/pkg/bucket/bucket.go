package bucket

import (
	"context"
	"linksvr/pkg/zlog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Bucket struct {
	OpenID     string "bson:`openid`"
	BucketName string "bson:`bucketname`"
	BucketID   string "bson:`bucketid`"
}

// Handler /bucket/{obejectName}
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
	// zlog.Info("handler verison Get")
	// openID := r.Header.Get("Authentication")
	// objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	// zlog.Info("start get object", zap.Any("obeject name", objectName))
	// ds := metadata.GetMongoMetaDataSource()
	// metas, err := ds.FindAllVersionMetaData(context.Background(), openID, objectName)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// data, err := json.Marshal(metas)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func put(w http.ResponseWriter, r *http.Request) {
	zlog.Info("handler verison Get")
	openID := r.Header.Get("Authentication")
	bucketName := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucketName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	zlog.Info("start put bucket", zap.Any("open id", openID), zap.Any("bucket name", bucketName))

	ds := GetMongoBucketDataSource()
	bucket, _ := ds.FindBucket(context.TODO(), openID, bucketName)
	if bucket == nil {
		bucketID := uuid.New().String()
		b := &Bucket{
			OpenID:     openID,
			BucketName: bucketName,
			BucketID:   bucketID,
		}
		if err := ds.InsertBucket(context.TODO(), b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		zlog.Info("bucket put sucessfully", zap.Any("open id", openID), zap.Any("bucket name", bucketName), zap.Any("bucket id", bucketID))
	}
	w.WriteHeader(http.StatusOK)
}
