package moyai

import (
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/paroxity/portal/socket"
	"github.com/sirupsen/logrus"

	proxypacket "github.com/paroxity/portal/socket/packet"
)

var s *socket.Client
var sMut sync.Mutex
var playerInfo chan *proxypacket.PlayerInfoResponse
var latencies chan *proxypacket.UpdatePlayerLatency

func init() {
	playerInfo = make(chan *proxypacket.PlayerInfoResponse)
	latencies = make(chan *proxypacket.UpdatePlayerLatency)

	addr, err := net.ResolveTCPAddr("tcp", "localhost:19131")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	sock := socket.NewClient(conn, logrus.New())
	sock.WritePacket(&proxypacket.AuthRequest{
		Protocol: 1,
		Secret:   "abc123",
		Name:     "syn.hcf",
	})
	sock.WritePacket(&proxypacket.RegisterServer{
		Address: "127.0.0.1:19134",
	})
	sock.Authenticate("syn.hcf")

	sMut.Lock()
	defer sMut.Unlock()
	s = sock

	go func() {
		for {
			sMut.Lock()
			pk, err := s.ReadPacket()
			sMut.Unlock()
			logrus.Info(pk)
			if err != nil {
				logrus.Info(err)
			}
			if pir, ok := pk.(*proxypacket.PlayerInfoResponse); ok {
				playerInfo <- pir
			} else if upl, ok := pk.(*proxypacket.UpdatePlayerLatency); ok {
				latencies <- upl
			}
		}
	}()
}

func SearchInfo(id uuid.UUID) (*proxypacket.PlayerInfoResponse, bool) {
	logrus.Info("ALLAH")
	for inf := range playerInfo {
		logrus.Info(inf)
		if inf.PlayerUUID != id {
			playerInfo <- inf
		} else {
			return inf, true
		}
	}
	return nil, false
}

func SearchLatency(id uuid.UUID) (*proxypacket.UpdatePlayerLatency, bool) {
	for inf := range latencies {
		if inf.PlayerUUID != id {
			latencies <- inf
		} else {
			return inf, true
		}
	}
	return nil, false
}

func PlayerInfo() chan *proxypacket.PlayerInfoResponse {
	return playerInfo
}

func Socket() *socket.Client {
	sMut.Lock()
	return s
}

func SocketUnlock() {
	sMut.Unlock()
}
