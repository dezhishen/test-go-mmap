package shared

import (
	"net/rpc"
)

// RPCClient is an implementation of KV that talks over RPC.
type MmapOperatorRPCClient struct{ client *rpc.Client }

func (m *MmapOperatorRPCClient) Write(filename string, content []byte) error {
	// We don't expect a response, so we can just use interface{}
	var resp interface{}

	// The args are just going to be a map. A struct could be better.
	return m.client.Call("Plugin.Write", map[string]interface{}{
		"filename": filename,
		"content":  content,
	}, &resp)
}

func (m *MmapOperatorRPCClient) Read(filename string) error {
	var resp interface{}
	err := m.client.Call("Plugin.Read", filename, &resp)
	return err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type MmapOperatorRPCServer struct {
	// This is the real implementation
	Impl MmapOperator
}

func (m *MmapOperatorRPCServer) Write(args map[string]interface{}, resp *interface{}) error {
	return m.Impl.Write(args["filename"].(string), args["content"].([]byte))
}

func (m *MmapOperatorRPCServer) Read(filename string, resp *interface{}) error {
	err := m.Impl.Read(filename)
	return err
}
