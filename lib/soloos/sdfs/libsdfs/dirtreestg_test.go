package libsdfs

import (
	"soloos/sdfs/types"
	"soloos/util"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/stretchr/testify/assert"
)

func TestMetaStgDirTreeStgBase(t *testing.T) {
	var (
		client      Client
		fsINode     types.FsINode
		netBlockCap = types.DefaultNetBlockCap
		memBlockCap = types.DefaultMemBlockCap
		code        fuse.Status
		err         error
	)
	MakeClientForTest(&client)

	code = client.MemDirTreeStg.SimpleMkdir(&fsINode, nil, types.RootFsINodeID, 0777, "test", 0, 0, types.FS_RDEV)
	if code != fuse.OK {
		assert.Equal(t, code, types.FS_EEXIST)
	}

	util.Ignore(fsINode)
	fsINode, err = client.MemDirTreeStg.OpenFile("/test/hi", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.MemDirTreeStg.OpenFile("/test/hi2", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.MemDirTreeStg.OpenFile("/test/hi3", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.MemDirTreeStg.OpenFile("/test/hi4", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.MemDirTreeStg.OpenFile("/test/hi5", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	err = client.MemDirTreeStg.DeleteFsINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = client.MemDirTreeStg.RenameWithFullPath("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = client.MemDirTreeStg.RenameWithFullPath("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = client.MemDirTreeStg.ListFsINodeByParentPath("/test", true,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount), uint64(0)
		},
		func(fsINode types.FsINode) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = client.MemDirTreeStg.OpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err, types.ErrObjectNotExists)
}