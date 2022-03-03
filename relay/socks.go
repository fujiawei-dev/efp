/*
 * @Date: 2022.03.02 10:58
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:58
 */

package relay

import (
	"net"

	netp "efp/net"
	"efp/net/socks"
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
