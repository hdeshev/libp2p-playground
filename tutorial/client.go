package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()
	serverAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/44033/p2p/QmTcJz3rEPXdBmjUDrrCArwkXywAhiRWdEm8Q2DNiUUexk")
	if err != nil {
		panic(fmt.Errorf("multiaddr %v", err))
	}

	node, err := libp2p.New(libp2p.Ping(false))
	defer node.Close()
	if err != nil {
		panic(fmt.Errorf("new node %v", err))
	}
	fmt.Println("Listen addresses", node.Addrs())

	peerInfo, err := peerstore.AddrInfoFromP2pAddr(serverAddr)
	if err != nil {
		panic(fmt.Errorf("peer addr %v", err))
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(peerInfo)
	fmt.Println("libp2p node address", addrs[0])

	if err := node.Connect(ctx, *peerInfo); err != nil {
		panic(fmt.Errorf("connect %v", err))
	}

	fmt.Println("sending pings to ", peerInfo)
	pingService := &ping.PingService{Host: node}
	ch := pingService.Ping(ctx, peerInfo.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("pinged", addrs[0], "in", res.RTT)
	}
}
