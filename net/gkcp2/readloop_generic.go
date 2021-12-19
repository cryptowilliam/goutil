// +build !linux

package kcp2

func (s *UDPSession) readLoop() {
	s.defaultReadLoop()
}

func (l *Listener) monitor() {
	l.defaultMonitor()
}
