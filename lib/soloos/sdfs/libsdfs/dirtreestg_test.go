package libsdfs

import (
	"soloos/common/fsapi"
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgDirTreeStgBase(t *testing.T) {
	var (
		client      Client
		fsINode     types.FsINode
		netBlockCap = types.DefaultNetBlockCap
		memBlockCap = types.DefaultMemBlockCap
		rawfs       fsapi.RawFileSystem
		code        fsapitypes.Status
		err         error
	)
	MakeClientForTest(&client)
	rawfs = client.GetRawFileSystem()

	code = rawfs.SimpleMkdir(&fsINode, nil, types.RootFsINodeID, 0777, "test", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		assert.Equal(t, code, types.FS_EEXIST)
	}

	util.Ignore(fsINode)
	fsINode, err = rawfs.SimpleOpenFile("/test/hi", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = rawfs.SimpleOpenFile("/test/hi2", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = rawfs.SimpleOpenFile("/test/hi3", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = rawfs.SimpleOpenFile("/test/hi4", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = rawfs.SimpleOpenFile("/test/hi5", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	err = rawfs.DeleteFsINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = rawfs.RenameWithFullPath("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = rawfs.RenameWithFullPath("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = rawfs.ListFsINodeByParentPath("/test", true,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount), uint64(0)
		},
		func(fsINode types.FsINode) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = rawfs.SimpleOpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err, types.ErrObjectNotExists)
}
