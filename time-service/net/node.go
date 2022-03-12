package net

import (
	"context"
	"fmt"

	"deshev.com/libp2p-time-service/protocol"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type Node struct {
	host.Host
	*protocol.TimeProtocol
}

func NewServer(listenPort uint) (*Node, error) {
	listenAddrOption := libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort))
	host, err := libp2p.New(listenAddrOption)
	if err != nil {
		return nil, fmt.Errorf("new_server: %w", err)
	}
	protocol := protocol.NewProtocol(host)
	return &Node{Host: host, TimeProtocol: protocol}, nil
}

func NewClient() (*Node, error) {
	host, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("new_client: %w", err)
	}
	protocol := protocol.NewProtocol(host)
	return &Node{Host: host, TimeProtocol: protocol}, nil
}

func (n *Node) Connect(ctx context.Context, serverAddress string) (*peerstore.AddrInfo, error) {
	addr, err := multiaddr.NewMultiaddr(serverAddress)
	if err != nil {
		return nil, fmt.Errorf("client_address: %w", err)
	}

	peerInfo, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil, fmt.Errorf("client_connect: %w", err)
	}

	err = n.Host.Connect(ctx, *peerInfo)
	return peerInfo, err
}
