package gicmp

import (
	"net"
	"time"
)

// https://hamy.io/post/000c/how-to-find-the-correct-mtu-and-mru-of-your-link/

type (
	Pong struct {
		RTT time.Duration
	}
)

func Ping(target string, payloadSize *uint16, timeout time.Duration) (*Pong, error) {
	targetAddr, err := net.ResolveIPAddr("ip", target)
	if err != nil {
		return nil, err
	}

	p, err := ping.New("0.0.0.0", "::")
	if err != nil {
		return nil, err
	}
	defer p.Close()

	if payloadSize != nil && *payloadSize > 0 {
		p.SetPayloadSize(*payloadSize)
	}

	RTT, err := p.Ping(targetAddr, timeout)
	if err != nil {
		return nil, err
	}
	return &Pong{RTT: RTT}, nil
}

// TODO
func DetectMTU(target string) (int, error) {
	return 0, nil
}
