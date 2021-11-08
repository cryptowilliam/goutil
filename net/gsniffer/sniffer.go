package gsniffer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/net/gsocks5/socks5internal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

/*
read from net.Conn and detect what protocol it is.

references:
https://github.com/soheilhy/cmux/blob/master/matchers.go
*/

type ConnInfo struct {
	Protocol     string // "http1","http2","websocket", "socks5"
	HTTPHeader   http.Header
	HTTPVer      string
	Socks5Target string // "www.google.com:443" "example.com:80"
}

func (ci *ConnInfo) checkWebsocket() {
	if ci.Protocol == "http" {
		if ci.HTTPHeader.Get("Upgrade") == "websocket" {
			ci.Protocol = "websocket"
		}
	}
}

// support:
// "http1","http2","websocket", "socks5"
// WARNING: data read by sniffler can't be read from conn again,
// but sniffer put what it read into readBuf, please handle readBuf manually
func Sniff(conn net.Conn) (r *ConnInfo, readBuf bytes.Buffer, ok bool) {
	/*
		you can't use reader := bufio.NewReader(conn), this will delete what sniffer read,
		then data will loss in this stream.
		io.Reader is treated like a stream. Because of this you cannot read it twice.
		Imagine the an incoming TCP connection. You cannot rewind the whats coming in.
		But you can use the io.TeeReader to duplicate the stream

		Technically, on one reader, you cannot read multiple times.
		Even if you create different references but
		when you read once it will be same object referred by all references.
		so what you can do is read the content and store it in one variable.
		Then use that variable as many times as you want.
	*/
	reader := io.TeeReader(conn, &readBuf)
	r = new(ConnInfo)

	req, err := tryGetHTTP1Request(reader)
	if err == nil {
		r.Protocol = "http"
		r.HTTPVer = "1"
		r.HTTPHeader = req.Header
		r.checkWebsocket()
		return r, readBuf, true
	}

	fields, err := tryGetHTTP2Request(reader)
	if err == nil {
		r.Protocol = "http"
		r.HTTPVer = "2"
		for _, s2 := range fields {
			r.HTTPHeader.Add(s2[0], s2[1])
		}
		r.checkWebsocket()
		return r, readBuf, true
	}

	return nil, readBuf, false
}

func tryGetHTTP1Request(r io.Reader) (*http.Request, error) {
	req, err := http.ReadRequest(bufio.NewReader(r))
	if err != nil {
		return nil, err
	}

	return req, nil
}

func hasHTTP2Preface(r io.Reader) bool {
	var b [len(http2.ClientPreface)]byte
	last := 0

	for {
		n, err := r.Read(b[last:])
		if err != nil {
			return false
		}

		last += n
		eq := string(b[:last]) == http2.ClientPreface[:last]
		if last == len(http2.ClientPreface) {
			return eq
		}
		if !eq {
			return false
		}
	}
}

func tryGetHTTP2Request(r io.Reader) (fields [][2]string, err error) {
	if !hasHTTP2Preface(r) {
		return nil, gerrors.Errorf("not HTTP 2")
	}

	w := ioutil.Discard
	done := false
	framer := http2.NewFramer(w, r)
	hdec := hpack.NewDecoder(uint32(4<<10), func(hf hpack.HeaderField) {
		fields = append(fields, [2]string{hf.Name, hf.Value})
	})
	for {
		f, err := framer.ReadFrame()
		if err != nil {
			return nil, err
		}

		switch f := f.(type) {
		case *http2.SettingsFrame:
			// Sender acknoweldged the SETTINGS frame. No need to write
			// SETTINGS again.
			if f.IsAck() {
				break
			}
			if err := framer.WriteSettings(); err != nil {
				return nil, err
			}
		case *http2.ContinuationFrame:
			if _, err := hdec.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			done = done || f.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0
		case *http2.HeadersFrame:
			if _, err := hdec.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			done = done || f.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0
		}

		if done {
			return fields, nil
		}
	}
}

func tryGetSocks5Request(r io.Reader) (socks5target string, err error) {
	errNotSocks5 := gerrors.Errorf("not socks5 protocol")
	buf := make([]byte, 4)

	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}

	if buf[0] != socks5internal.Version {
		return "", errNotSocks5
	}

	if socks5internal.Command(buf[1]) < socks5internal.CommandMin || socks5internal.Command(buf[1]) > socks5internal.CommandMax {
		return "", gerrors.Errorf("invalid command %d", buf[1])
	}

	addrType := socks5internal.AddressType(buf[3])

	switch addrType {
	case socks5internal.AddrIPv4:
		buf = make([]byte, 4+2)
		if _, err := io.LimitReader(r, 4+2).Read(buf); err != nil {
			return "", err
		}

		ipAddr := net.IPv4(buf[0], buf[1], buf[2], buf[3])
		port := uint16(buf[4])<<8 | uint16(buf[5])
		return fmt.Sprintf("%s:%d", ipAddr.String(), port), nil
	case socks5internal.AddrDomainName:
		bufLen := make([]byte, 1)
		_, err := r.Read(bufLen) // Read 1 byte.
		if err != nil {
			return "", err
		}
		domainLen := bufLen[0]

		buf = make([]byte, domainLen+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", err
		}
		domain := string(buf[:domainLen])
		port := uint16(buf[domainLen])<<8 | uint16(buf[domainLen+1])
		return fmt.Sprintf("%s:%d", domain, port), nil
	default:
		return "", gerrors.Errorf("invalid or unsupported address type(%d)", buf[3])
	}
}
