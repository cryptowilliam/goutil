package grpcs

type (
	RpcType string
)

var (
	RpcTypeGOB    = RpcType("gob")
	RpcTypeJSON   = RpcType("json")
	RpcTypeHPROSE = RpcType("hprose")
)
