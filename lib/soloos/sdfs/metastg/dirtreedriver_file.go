package metastg

import (
	"soloos/log"
	"strings"
)

func (p *DirTreeDriver) OpenFile(fsInodePath string) error {
	var (
		err   error
		paths []string
	)
	paths = strings.Split(fsInodePath, "/")
	for path := range paths {
		log.Error(path)
	}

	return err
}
