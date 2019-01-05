package main

import (
	"encoding/json"
	"io/ioutil"
)

type Options struct {
	NameNodeSRPCServerAddr string
	MemBlockChunkSize      int
	MemBlockChunksLimit    int32
	DBDriver               string
	Dsn                    string
	PProfListenAddr        string
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
