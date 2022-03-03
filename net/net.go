/*
 * @Date: 2022.03.02 10:05
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:05
 */

package net

import "net"

type ConnCipher func(net.Conn) net.Conn

type listener struct {
	net.Listener
	ConnCipher
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	return l.ConnCipher(c), err
}

type PacketConnCipher func(net.PacketConn) net.PacketConn
