package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func main() {
	listenAddrs := libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0")
	node, err := libp2p.New(listenAddrs, libp2p.Ping(false))
	defer node.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Listen addresses", node.Addrs())

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := &peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(peerInfo)
	fmt.Println("peer addrs", addrs[0])

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("Received signal, shutting down...")
}
