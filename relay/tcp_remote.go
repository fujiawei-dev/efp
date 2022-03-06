/*
 * @Date: 2022.03.02 10:17
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:17
 */

package relay

import (
	"net"

	netp "github.com/fujiawei-dev/efp/net"
	"github.com/fujiawei-dev/efp/net/socks"
)

func NewTCPRemoteProxyServer(addr string, cipher netp.ConnCipher) {
	proxy := TCPRemoteProxyServer{
		Addr:   addr,
		cipher: cipher,
	}

	if err := proxy.ListenAndServe(); err != nil {
		logf("TCP proxy error, %v", err)
	}
}

type TCPRemoteProxyServer struct {
	Addr string

	cipher netp.ConnCipher
}

func (p *TCPRemoteProxyServer) ListenAndServe() error {
	listener, err := netp.Listen("tcp", p.Addr, p.cipher)
	if err != nil {
		logf("failed to listen on %s, %v", p.Addr, err)
		return err
	}

	logf("listening TCP on %s", p.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logf("failed to accept, %s", err)
			continue
		}

		go p.relay(conn)
	}
}

func (p *TCPRemoteProxyServer) relay(leftConn net.Conn) {
	defer leftConn.Close()

	targetAddr, err := socks.ReadAddr(leftConn)
	if err != nil {
		logf("failed to read address, %v", err)
		return
	}

	logf("proxy %s <-> %s", leftConn.RemoteAddr(), targetAddr)

	rightConn, err := net.Dial("tcp", targetAddr.String())
	if err != nil {
		logf("failed to connect to target, %s", err)
		return
	}

	defer rightConn.Close()

	if _, _, err = relay(leftConn, rightConn); err != nil {
		logf("relay error, %v", err)
		return
	}
}
