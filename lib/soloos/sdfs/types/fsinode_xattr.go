package types

import (
	"bytes"
	"encoding/gob"
)

type FsINodeXAttr map[string][]byte

func InitFsINodeXAttr() FsINodeXAttr {
	return FsINodeXAttr(make(map[string][]byte))
}

func SerializeFIXAttr(xattr FsINodeXAttr) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(xattr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserializeFIXAttr(data []byte, xattr *FsINodeXAttr) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(xattr)
}
