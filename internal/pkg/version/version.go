package version

import (
	"linksvr/pkg/zlog"
	"net/http"
)

// Handler /version/{obejectName}
// curl -v 127.0.0.1:4000/version/mysh1  -H "Authentication:fengwei"
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

	// if m == http.MethodPut {
	// 	put(w, r)
	// }

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
