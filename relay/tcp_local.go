/*
 * @Date: 2022.03.02 10:16
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:16
 */

package relay

import (
	"net"

	netp "github.com/fujiawei-dev/efp/net"
	"github.com/fujiawei-dev/efp/net/socks"
)

// NewSOCKS5ProxyClient Create a SOCKS5 server listening on address and proxy to server.
func NewSOCKS5ProxyClient(clientAddr, serverAddr string, cipher netp.ConnCipher) {
	logf("SOCKS proxy %s <-> %s", clientAddr, serverAddr)

	client := TCPLocalProxyClient{
		Addr:       clientAddr,
		RemoteAddr: serverAddr,
		cipher:     cipher,
		targetAddrExtractor: func(c net.Conn) (socks.Addr, error) {
			return socks.Handshake(c)
		},
	}

	if err := client.ListenAndServe(); err != nil {
		logf("SOCKS proxy error, %v", err)
	}
}

// NewTCPTunnel Create a TCP tunnel from clientAddr to remoteAddr via server (serverAddr).
func NewTCPTunnel(clientAddr, serverAddr, remoteAddr string, cipher netp.ConnCipher) {
	targetAddr := socks.ParseAddr(remoteAddr)
	if targetAddr == nil {
		logf("invalid target address %q", remoteAddr)
		return
	}

	logf("TCP tunnel %s <-> %s <-> %s", clientAddr, serverAddr, remoteAddr)

	client := TCPLocalProxyClient{
		Addr:       clientAddr,
		RemoteAddr: serverAddr,
		cipher:     cipher,
		targetAddrExtractor: func(c net.Conn) (socks.Addr, error) {
			return targetAddr, nil
		},
	}

	if err := client.ListenAndServe(); err != nil {
		logf("SOCKS proxy error, %v", err)
	}
}

type TCPLocalProxyClient struct {
	Addr       string
	RemoteAddr string

	cipher              netp.ConnCipher
	targetAddrExtractor func(net.Conn) (socks.Addr, error)
}

// ListenAndServe Listen on Addr and proxy to RemoteAddr to reach target from targetAddrExtractor.
func (p *TCPLocalProxyClient) ListenAndServe() (err error) {
	listener, err := net.Listen("tcp", p.Addr)
	if err != nil {
		logf("failed to listen on %s: %v", p.Addr, err)
		return
	}

	logf("listening TCP on %s", p.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logf("failed to accept, %v", err)
			continue
		}

		targetAddr, err := p.targetAddrExtractor(conn)
		if err != nil {
			logf("failed to get target address, %v", err)
			continue
		}

		go p.relay(conn, targetAddr)
	}
}

func (p *TCPLocalProxyClient) relay(leftConn net.Conn, targetAddr socks.Addr) {
	defer leftConn.Close()

	logf("relay %s <-> %s <-> %s", leftConn.RemoteAddr(), p.RemoteAddr, targetAddr)

	rightConn, err := netp.Dial("tcp", p.RemoteAddr, p.cipher)
	if err != nil {
		logf("failed to connect to server %v, %v", p.RemoteAddr, err)
		return
	}

	defer rightConn.Close()

	if _, err = rightConn.Write(targetAddr); err != nil {
		logf("failed to send target address, %v", err)
		return
	}

	if _, _, err = relay(leftConn, rightConn); err != nil {
		logf("relay error, %v", err)
		return
	}
}
