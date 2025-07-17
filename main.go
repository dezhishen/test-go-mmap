// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/dezhishen/test-go-mmap/shared"
	"github.com/edsrzf/mmap-go"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func main() {
	runPlugin()
}

func runTest() {
}
func runPlugin() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	currentOS := runtime.GOOS
	var pluginSuffix string
	if currentOS == "windows" {
		pluginSuffix = ".exe"
	} else {
		pluginSuffix = ""
	}
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("./plugin/mmap_operator" + pluginSuffix),
		Logger:          logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC},
	})
	defer client.Kill()
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("mmap")
	if err != nil {
		log.Fatal(err)
	}
	operator := raw.(shared.MmapOperator)
	log.Println("成功获取插件对象")
	// 创建文件
	filename := "test.data"
	contents := []string{"test1", "test2"}
	defer os.Remove(filename)
	// 主程序读取内容
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Panic(err)
	}
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		log.Panic(err)
	}
	defer m.Unmap()
	// z
	for index, content := range contents {
		log.Printf("第%d次操作", index)
		log.Println("开始通过插件写入")
		err = operator.Write(filename, []byte(content))
		if err != nil {
			log.Panic(err)
		}
		log.Println("宿主程序读取")
		log.Println(string(m))
		log.Println("插件中读取")
		err = operator.Read(filename)
		if err != nil {
			log.Panic(err)
		}
	}
}
