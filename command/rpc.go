package command

import "flag"

// RPCAddrFlag returns a pointer to a string that will be populated
// when the given flagset is parsed with the RPC address of the sr6.
func RPCAddrFlag(f *flag.FlagSet) *string {
	defaultRPCAddr := "127.0.0.1:8300"
	return f.String("rpc-addr", defaultRPCAddr,
		"RPC address of the sr6 agent")
}
