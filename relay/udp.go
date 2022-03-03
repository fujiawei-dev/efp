/*
 * @Date: 2022.03.02 11:14
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:14
 */

package relay

import (
	"net"
	"sync"
	"time"
)

const UDPBufferSize = 64 * 1024
const UDPTimeout = time.Minute * 3

// NATMap Packet NAT table
type NATMap struct {
	pcs     map[string]net.PacketConn
	timeout time.Duration
	sync.RWMutex
}

func NewNATMap(timeout time.Duration) *NATMap {
	return &NATMap{
		pcs:     make(map[string]net.PacketConn),
		timeout: timeout,
	}
}

func (m *NATMap) Get(key string) net.PacketConn {
	m.RLock()
	defer m.RUnlock()
	return m.pcs[key]
}

func (m *NATMap) Set(key string, pc net.PacketConn) {
	m.Lock()
	defer m.Unlock()
	m.pcs[key] = pc
}

func (m *NATMap) Del(key string) net.PacketConn {
	m.Lock()
	defer m.Unlock()

	if pc, ok := m.pcs[key]; ok {
		delete(m.pcs, key)
		return pc
	}

	return nil
}

func (m *NATMap) Add(peer net.Addr, dst, src net.PacketConn) {
	m.Set(peer.String(), src)

	go func() {
		timedCopy(dst, peer, src, m.timeout)
		if pc := m.Del(peer.String()); pc != nil {
			pc.Close()
		}
	}()
}

// copy from src to dst with addr with read timeout
func timedCopy(dst net.PacketConn, addr net.Addr, src net.PacketConn, timeout time.Duration) error {
	buf := make([]byte, UDPBufferSize)

	for {
		src.SetReadDeadline(time.Now().Add(timeout))
		n, _, err := src.ReadFrom(buf)
		if err != nil {
			return err
		}

		_, err = dst.WriteTo(buf[:n], addr)
		if err != nil {
			return err
		}
	}
}
