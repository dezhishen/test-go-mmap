// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/dezhishen/test-go-mmap/shared"
	"github.com/edsrzf/mmap-go"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type MmapOperatorImpl struct{}

func (MmapOperatorImpl) Write(filename string, content []byte) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// 扩展文件
	if err := f.Truncate(int64(len(content))); err != nil {
		return err
	}
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	defer m.Unmap()
	copy(m, []byte(content))
	return m.Flush()
}

func (MmapOperatorImpl) Read(filename string) error {
	// 主程序读取内容
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	m, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer m.Unmap()
	log.Println(string(m))
	return nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"mmap": &shared.MmapOperatorPlugin{Impl: &MmapOperatorImpl{}},
		},
	})
}
