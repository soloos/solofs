package main

import (
	"encoding/json"
	"io/ioutil"
)

type Options struct {
	Mode                string
	DataNodePeerIDStr   string
	ListenAddr          string
	DataNodeLocalFsRoot string
	NameNodePeerIDStr   string
	NameNodeAddr        string
	PProfListenAddr     string
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
