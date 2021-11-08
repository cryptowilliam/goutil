package socks5internal

/**


Version identifier/method selection request message:
	  +----+----------+----------+
	  |VER | NMETHODS | METHODS  |
	  +----+----------+----------+
	  | 1  |    1     | 1 to 255 |
	  +----+----------+----------+
VER field is set to X'05' for this ver of the protocol.
NMETHODS field contains the number of method identifier octets that appear in the METHODS field, in other words, NMETHODS is the length of METHODS.
METHODS field is a list of authentication methods supported by the client. Each method occupies 1 byte.



Version identifier/method selection reply message (without any auth method):
	  +----+--------+
	  |VER | METHOD |
	  +----+--------+
	  | 1  |   1    |
	  +----+--------+
VER field is set to X'05' for this ver of the protocol.
METHOD: set it as 0 if no auth required




Command request message:
	  +----+-----+-------+------+----------+----------+
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
VER: socks proxy version, it is 5 in socks5
CMD: Connect: X'01', Bind: X'02', UdpAssociate: X'03'
RSV: reserve, should be set as 0
ATYP: address type, IP V4 address: X’01’, DOMAINNAME: X’03’, IP V6 address: X’04’
DST.ADDR：dest address, could be ipv4/ipv6/domain-name, if it is domain-name, it ends with 0?
DST.PORT：dest port




Command reply message:
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
REP:
X’00’ succeeded
X’01’ general socks server failure
X’02’ connection not allowed by ruleset
X’03’ Network unreachable
X’04’ Host unreachable
X’05’ Connection refused
X’06’ TTL expired
X’07’ Command not supported
X’08’ Address type not supported
X’09’ to X’FF’ unassigned
*/

const (
	Version = uint8(5)

	AddrIPv4       AddressType = 1
	AddrDomainName AddressType = 3 // domain name
	AddrIPv6       AddressType = 4

	CommandConnect      Command = 1
	CommandBind         Command = 2
	CommandUdpAssociate Command = 3
	CommandMin                  = CommandConnect
	CommandMax                  = CommandUdpAssociate

	ReplySuccess            Reply = 0
	ReplyFailed             Reply = 1
	ReplyUnreached          Reply = 3
	ReplyNoSuchHost         Reply = 4
	ReplyConnectDenied      Reply = 5
	ReplyTtlOver            Reply = 6
	ReplyUnsupportedCommand Reply = 7
	ReplyUnsupportedAddress Reply = 8
)

type (
	Command     byte
	AddressType byte
	Reply       byte

	VimsReq struct {
		Ver      byte
		NMethods byte
		Methods  []byte // 1-255 bytes
	}

	VimsResp struct {
	}

	CommandReq struct {
		Ver      byte
		Cmd      Command
		Reverse  byte
		AddrType AddressType
		DstAddr  string
		DstPort  uint16
	}

	CommandResp struct {
		Ver      byte
		Rep      Reply
		Reverse  byte
		AddrType AddressType
		DstAddr  string
		DstPort  uint16
	}
)

func (cmd Command) String() string {
	switch cmd {
	case CommandConnect:
		return "CommandConnect"
	case CommandBind:
		return "CommandBind"
	case CommandUdpAssociate:
		return "CommandUdpAssociate"
	default:
		return "Unknown"
	}
}

func (at AddressType) String() string {
	switch at {
	case AddrIPv4:
		return "AddrIPv4"
	case AddrDomainName:
		return "AddrDomainName"
	case AddrIPv6:
		return "AddrIPv6"
	default:
		return "Unknown"
	}
}

func (r Reply) String() string {
	switch r {
	case ReplySuccess:
		return "ReplySuccess"
	case ReplyFailed:
		return "ReplyFailed"
	case ReplyUnreached:
		return "ReplyUnreached"
	case ReplyNoSuchHost:
		return "ReplyNoSuchHost"
	case ReplyConnectDenied:
		return "ReplyConnectDenied"
	case ReplyTtlOver:
		return "ReplyTtlOver"
	case ReplyUnsupportedCommand:
		return "ReplyUnsupportedCommand"
	case ReplyUnsupportedAddress:
		return "ReplyUnsupportedAddress"
	default:
		return "Unknown"
	}
}
