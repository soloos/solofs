package main

import (
	"encoding/json"
	"io/ioutil"
)

type Options struct {
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
	Mode                  string

	SNetDriverListenAddr string
	SNetDriverServeAddr  string
	ServeAddr            string
	ListenAddr           string

	DataNodePeerID      string
	DataNodeLocalFSRoot string

	NameNodePeerID string

	PProfListenAddr string
	DBDriver        string
	Dsn             string
}

func LoadOptionsFile(optionsFilePath string) (Options, error) {
	var (
		err     error
		content []byte
		options Options
	)

	content, err = ioutil.ReadFile(optionsFilePath)
	if err != nil {
		return options, err
	}

	err = json.Unmarshal(content, &options)
	if err != nil {
		return options, err
	}

	return options, nil
}
