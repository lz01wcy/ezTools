package ezNetworking

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

func ListenUDP(port int, sendBufferSize int, recvBufferSize int) (*net.UDPConn, error) {
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error {
		if sendBufferSize > 0 {
			err := setIOBufferSize(c, sendBufferSize, syscall.SO_SNDBUF)
			if err != nil {
				return err
			}
		}
		if recvBufferSize > 0 {
			err := setIOBufferSize(c, recvBufferSize, syscall.SO_RCVBUF)
			if err != nil {
				return err
			}
		}
		return nil
	}}
	l, err := lc.ListenPacket(context.Background(), "udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return l.(*net.UDPConn), nil
}
func ListenTCP(port int, sendBufferSize int, recvBufferSize int) (*net.TCPListener, error) {
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error {
		if sendBufferSize > 0 {
			err := setIOBufferSize(c, sendBufferSize, syscall.SO_SNDBUF)
			if err != nil {
				return err
			}
		}
		if recvBufferSize > 0 {
			err := setIOBufferSize(c, recvBufferSize, syscall.SO_RCVBUF)
			if err != nil {
				return err
			}
		}
		return nil
	}}
	l, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return l.(*net.TCPListener), nil
}
