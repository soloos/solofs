package main

import (
	"encoding/json"
	"io/ioutil"
	"soloos/sfuse"
)

type Options struct {
	NameNodeSRPCServerAddr string
	DefaultNetBlockCap     int
	DefaultMemBlockCap     int
	DefaultMemBlocksLimit  int32
	DBDriver               string
	Dsn                    string
	PProfListenAddr        string
	SFuseOptions           sfuse.Options
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
