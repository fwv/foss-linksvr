package osd

import (
	"bytes"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"linksvr/internal/pkg/config"
	"linksvr/internal/pkg/ssclient"
	"linksvr/pkg/zlog"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/foss/osdsvr/pkg/proto/osdpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Selector = NewOSDSelector()
)

type OSDSelector struct {
	mu         sync.RWMutex
	OsdClients map[int64]*ssclient.OsdSerivce
}

func NewOSDSelector() *OSDSelector {
	s := &OSDSelector{
		OsdClients: make(map[int64]*ssclient.OsdSerivce),
	}
	return s
}

func (r *OSDSelector) Init(done chan bool) error {
	zlog.Info("waiting for osd service to register", zap.Any("need osd num", *config.OSD_NUM))
	for {
		if r.GetOSDSize() >= *config.OSD_NUM {
			break
		}
	}
	done <- true
	return nil
}

func (r *OSDSelector) GetOSDSize() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.OsdClients)
}

func (r *OSDSelector) AddOSDService(osdAddr string, osdID int64) error {
	if osdAddr == "" || osdID == 0 {
		return errors.New("bad request")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// if r.osdClients[osdID] != nil {
	// 	zlog.Error("osd service register repeatedlly", zap.Any("osd addr", osdAddr), zap.Int64("osd id", osdID))
	// 	return errors.New("osd service register repeatedlly")
	// }
	conn, err := grpc.Dial(osdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	// defer conn.Close()
	c := osdpb.NewOsdServiceClient(conn)
	r.OsdClients[osdID] = ssclient.NewOsdService(c)
	zlog.Info("osd service register successfully", zap.Any("osd addr", osdAddr), zap.Int64("osd id", osdID))
	return nil
}

func (r *OSDSelector) ChooseOSDbyRandom() *ssclient.OsdSerivce {
	len := r.GetOSDSize()
	return r.OsdClients[int64(rand.Intn(len))]
}

func (r *OSDSelector) ChooseOSDbyHash(source string) (*ssclient.OsdSerivce, error) {
	// use sha1 convert to int
	hash := sha1.New()
	hash.Write([]byte(source))
	hashBytes := hash.Sum(nil)

	// conversion to base32
	base32str := strings.ToLower(base32.HexEncoding.EncodeToString(hashBytes))
	// no, err := strconv.ParseInt(base32str, 10, 32)
	no := BytesToInt([]byte(base32str))
	// r.mu.Lock()
	// defer r.mu.Unlock()
	size := len(r.OsdClients)
	if size == 0 {
		return nil, errors.New("there is no avaliable osd service")
	}
	if no < 0 {
		no = -no
	}
	idx := no%size + 1
	service := r.OsdClients[int64(idx)]
	return service, nil
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}
