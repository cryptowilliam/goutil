package grpcs

// Notice:
// If a function of the server-side registered Receiver returns an error,
// the result of the output parameter will not be transmitted to client when the function is called by RPC.

// Registered Receiver member function sample:
// (r *Recv) Method(in InputParam, out *OutputParam) error

import (
	"github.com/cryptowilliam/goutil/net/gnet"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync/atomic"
)

type (
	Server struct {
		rpcSvr        *rpc.Server
		netSvr        net.Listener
		paramChecker  ParamChecker
		onRequestUser OnRequest
		rpcType       RpcType
	}

	Svr Server

	InitChecker func() ParamChecker
	OnRequest   func(in Request, out *Reply) error
)

var (
	jsonRpcRegisterDone = int32(0)
)

func Listen(rpcType RpcType, network, address string) (*Server, error) {
	netSer, err := gnet.Listen(network, address)
	if err != nil {
		return nil, err
	}
	res := &Server{
		netSvr:  netSer,
		rpcSvr:  rpc.NewServer(),
		rpcType: rpcType,
	}
	return res, nil
}

func (s *Svr) OnRequestInternal(in Request, out *Reply) error {
	/*if err := s.paramChecker.VerifyIn(in.Func, in); err != nil {
		return err
	}*/
	if err := s.onRequestUser(in, out); err != nil {
		return err
	}
	if err := s.paramChecker.VerifyOut(in.Func, out, false); err != nil {
		return err
	}
	return nil
}

// Before every 'onReq' call, rpc server will check input and out put param with 'checker'.
func (s *Server) Run(checker ParamChecker, onReq OnRequest) error {
	s.paramChecker = checker
	s.onRequestUser = onReq

	// Register
	if s.rpcType == RpcTypeJSON {
		// In json rpc, types can be registered only once,
		// no matter how many different network servers started.
		if atomic.LoadInt32(&jsonRpcRegisterDone) == 0 {
			atomic.StoreInt32(&jsonRpcRegisterDone, 1)
			if err := rpc.Register((*Svr)(s)); err != nil {
				return err
			}
		}
	} else if s.rpcType == RpcTypeGOB {
		if err := s.rpcSvr.Register((*Svr)(s)); err != nil {
			return err
		}
	}

	// Accept
	if s.rpcType == RpcTypeJSON {
		for {
			conn, e := s.netSvr.Accept()
			if e != nil {
				continue
			}
			go jsonrpc.ServeConn(conn)
		}
	} else if s.rpcType == RpcTypeGOB {
		s.rpcSvr.Accept(s.netSvr)
	}
	return nil
}

func (s *Server) Close() error {
	return s.netSvr.Close()
}
