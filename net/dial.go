/*
 * @Date: 2022.03.02 10:10
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:10
 */

package net

import "net"

func Listen(network, address string, cipher ConnCipher) (net.Listener, error) {
	l, err := net.Listen(network, address)
	return &listener{Listener: l, ConnCipher: cipher}, err
}

func Dial(network, address string, cipher ConnCipher) (net.Conn, error) {
	c, err := net.Dial(network, address)
	return cipher(c), err
}

func ListenPacket(network, address string, cipher PacketConnCipher) (net.PacketConn, error) {
	c, err := net.ListenPacket(network, address)
	return cipher(c), err
}
