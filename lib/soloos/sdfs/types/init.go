package types

import "encoding/gob"

func init() {
	gob.Register(FsINodeXAttr{})
}
