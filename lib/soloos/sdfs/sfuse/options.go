package sfuse

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Options struct {
	NameSpaceID int64
	MountPoint  string

	NodefsOptsTimeoutMs           uint32
	TimeDurationNodefsOptsTimeout time.Duration `json:"-"`
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

	options.TimeDurationNodefsOptsTimeout = time.Millisecond * time.Duration(options.NodefsOptsTimeoutMs)

	return options, nil
}
