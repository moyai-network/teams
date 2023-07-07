package moyai

import (
	"github.com/google/uuid"
	"github.com/paroxity/portal/socket"
	"github.com/sirupsen/logrus"
	"net"

	proxypacket "github.com/paroxity/portal/socket/packet"
)

var (
	sock  *socket.Client
	xuids = map[uuid.UUID]chan PlayerInformation{}
)

type PlayerInformation struct {
	XUID    string
	Address string
}

func init() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:19131")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	sock = socket.NewClient(conn, logrus.New())
	_ = sock.WritePacket(&proxypacket.AuthRequest{
		Protocol: 1,
		Secret:   "abc123",
		Name:     "syn.hcf",
	})
	_ = sock.WritePacket(&proxypacket.RegisterServer{
		Address: "127.0.0.1:19134",
	})

	sock.Authenticate("syn.hcf")

	go func() {
		for {
			pk, err := sock.ReadPacket()
			if err != nil {
				panic(err)
			}

			if i, ok := pk.(*proxypacket.PlayerInfoResponse); ok {
				xuids[i.PlayerUUID] <- PlayerInformation{XUID: i.XUID, Address: i.Address}
			}
		}
	}()
}

func SearchInfo(id uuid.UUID) PlayerInformation {
	ch := make(chan PlayerInformation)
	xuids[id] = ch

	defer delete(xuids, id)
	defer close(ch)

	_ = sock.WritePacket(&proxypacket.PlayerInfoRequest{
		PlayerUUID: id,
	})

	return <-ch
}

func Socket() *socket.Client {
	return sock
}
