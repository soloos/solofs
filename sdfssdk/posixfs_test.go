package sdfssdk

import (
	"soloos/common/fsapi"
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/util"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgPosixFSBase(t *testing.T) {
	var (
		client      Client
		fsINodeMeta sdfsapitypes.FsINodeMeta
		netBlockCap = types.DefaultNetBlockCap
		memBlockCap = types.DefaultMemBlockCap
		posixFS     fsapi.PosixFS
		code        fsapitypes.Status
		err         error
	)
	MakeClientForTest(&client)
	posixFS = client.GetPosixFS()

	code = posixFS.SimpleMkdir(&fsINodeMeta, nil, sdfsapitypes.RootFsINodeID, 0777, "test", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		assert.Equal(t, code, types.FS_EEXIST)
	}

	util.Ignore(fsINodeMeta)
	fsINodeMeta, err = posixFS.SimpleOpenFile("/test/hi", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFS.SimpleOpenFile("/test/hi2", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFS.SimpleOpenFile("/test/hi3", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFS.SimpleOpenFile("/test/hi4", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFS.SimpleOpenFile("/test/hi5", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	err = posixFS.DeleteFsINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = posixFS.RenameWithFullPath("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = posixFS.RenameWithFullPath("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = posixFS.ListFsINodeByParentPath("/test", true,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount), uint64(0)
		},
		func(fsINodeMeta sdfsapitypes.FsINodeMeta) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = posixFS.SimpleOpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err, sdfsapitypes.ErrObjectNotExists)
}
