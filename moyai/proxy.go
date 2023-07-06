package moyai

import (
	"net"

	"github.com/paroxity/portal/socket"
	"github.com/sirupsen/logrus"

	proxypacket "github.com/paroxity/portal/socket/packet"
)

var s *socket.Client

func init() {
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

	go func() {

	}()

	s = sock
}

func Socket() *socket.Client {
	return s
}
