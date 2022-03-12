package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"deshev.com/libp2p-time-service/net"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	listenPort := flag.Uint("p", 0, "listen port")
	server := flag.String("s", "", "server address")
	flag.Parse()

	fmt.Println("server: ", *server)
	fmt.Println("listenPort: ", *listenPort)
	fmt.Println("Ctrl-C to exit...")

	if *listenPort == 0 && *server == "" {
		*listenPort = 9000
	}

	ctx := context.Background()
	if *server == "" {
		defer runServer(ctx, *listenPort)()
	} else {
		defer connectToServer(ctx, *server)()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

type cleanup func()

func runServer(ctx context.Context, listenPort uint) cleanup {
	server, err := net.NewServer(listenPort)
	if err != nil {
		return func() {}
	}

	peerInfo := &peerstore.AddrInfo{
		ID:    server.ID(),
		Addrs: server.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(peerInfo)
	fmt.Println("server address", addrs[0])

	return func() {
		fmt.Println("closing server")
		server.Close()
	}
}

func connectToServer(ctx context.Context, serverAddress string) cleanup {
	client, err := net.NewClient()
	if err != nil {
		return func() {}
	}

	addrInfo, err := client.Connect(ctx, serverAddress)
	if err != nil {
		fmt.Printf("Connection error: %v", err)
		return func() {}
	}

	err = client.StartRequest(ctx, addrInfo.ID)
	if err != nil {
		fmt.Printf("Request error: %v", err)
		return func() {}
	}

	return func() {
		fmt.Println("closing client")
		client.Close()
	}
}
