package moyai

import (
	"net"

	"github.com/google/uuid"
	"github.com/paroxity/portal/socket"
	"github.com/sirupsen/logrus"

	portal "github.com/paroxity/portal/socket/packet"
)

var (
	sock  *socket.Client
	xuids = map[uuid.UUID]chan PlayerInformation{}
)

type PlayerInformation struct {
	XUID    string
	Address string
}

func NewProxySocket() *socket.Client {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:19131")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	sock = socket.NewClient(conn, logrus.New())
	_ = sock.WritePacket(&portal.AuthRequest{
		Protocol: 1,
		Secret:   "abc123",
		Name:     "syn.hcf",
	})
	_ = sock.WritePacket(&portal.RegisterServer{
		Address: "127.0.0.1:19134",
	})

	sock.Authenticate("syn.hcf")

	go func() {
		for {
			pk, err := sock.ReadPacket()
			if err != nil {
				panic(err)
			}

			if i, ok := pk.(*portal.PlayerInfoResponse); ok {
				xuids[i.PlayerUUID] <- PlayerInformation{XUID: i.XUID, Address: i.Address}
			}
		}
	}()
	return sock
}

func SearchInfo(id uuid.UUID) PlayerInformation {
	ch := make(chan PlayerInformation)
	xuids[id] = ch

	defer delete(xuids, id)
	defer close(ch)

	_ = sock.WritePacket(&portal.PlayerInfoRequest{
		PlayerUUID: id,
	})

	return <-ch
}

func Socket() (*socket.Client, bool) {
	return sock, sock != nil
}
