package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"mmap": &MmapOperatorPlugin{},
}

type MmapOperator interface {
	Write(filename string, content []byte) error
	Read(filename string) error
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type MmapOperatorPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl MmapOperator
}

func (p *MmapOperatorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &MmapOperatorRPCServer{Impl: p.Impl}, nil
}

func (*MmapOperatorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &MmapOperatorRPCClient{client: c}, nil
}
