/*
 * @Date: 2022.03.02 11:12
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:12
 */

package relay

import (
	"net"

	netp "efp/net"
	"efp/net/socks"
)

// NewUDPTunnel Listen on localAddr for UDP packets, encrypt and send to remoteAddr to reach targetAddr.
func NewUDPTunnel(localAddr, remoteAddr, targetAddr string, cipher netp.PacketConnCipher) {
	remoteServerAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		logf("UDP remoteAddr error, %v", err)
		return
	}

	targetSOCKSAddr := socks.ParseAddr(targetAddr)
	if targetSOCKSAddr == nil {
		logf("UDP targetAddr error, %v", err)
		return
	}

	c, err := net.ListenPacket("udp", localAddr)
	if err != nil {
		logf("UDP local listen error, %v", err)
		return
	}
	defer c.Close()

	nm := NewNATMap(UDPTimeout)
	buf := make([]byte, UDPBufferSize)
	copy(buf, targetSOCKSAddr)

	logf("UDP tunnel %s <-> %s <-> %s", localAddr, remoteAddr, targetAddr)

	for {
		n, raddr, err := c.ReadFrom(buf[len(targetSOCKSAddr):])
		if err != nil {
			logf("UDP local read error, %v", err)
			continue
		}

		pc := nm.Get(raddr.String())
		if pc == nil {
			pc, err = net.ListenPacket("udp", "")
			if err != nil {
				logf("UDP local listen error, %v", err)
				continue
			}

			pc = cipher(pc)
			nm.Add(raddr, c, pc)
		}

		if _, err = pc.WriteTo(buf[:len(targetSOCKSAddr)+n], remoteServerAddr); err != nil {
			logf("UDP local write error, %v", err)
			continue
		}
	}
}
