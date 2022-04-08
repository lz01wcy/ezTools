package ezGRPCConnPool

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

type Pool struct {
	crtPath      string
	serverName   string
	url          string
	connChannels chan *grpc.ClientConn
}

func NewPool(crtPath string, serverName string, url string, count int) (*Pool, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be greater than 0")
	}
	rs := &Pool{
		crtPath,
		serverName,
		url,
		make(chan *grpc.ClientConn, count),
	}
	for i := 0; i < count; i++ {
		conn, err := rs.newConn()
		if err != nil {
			rs.Free()
			return nil, err
		}
		rs.connChannels <- conn
	}
	return rs, nil
}
func (p *Pool) Use(f func(conn *grpc.ClientConn, err error)) {
	conn := <-p.connChannels
	var err error
	if conn.GetState() != connectivity.Idle {
		_ = conn.Close()
		conn, err = p.newConn()
		if err != nil {
			f(nil, err)
			return
		}
	}
	f(conn, err)
	p.connChannels <- conn
}
func (p *Pool) Free() {
	for {
		select {
		case conn := <-p.connChannels:
			_ = conn.Close()
			break
		default:
			return
		}
	}
}
func (p *Pool) newConn() (*grpc.ClientConn, error) {
	crt, err := credentials.NewClientTLSFromFile(p.crtPath, p.serverName)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(p.url, grpc.WithTransportCredentials(crt))
}
