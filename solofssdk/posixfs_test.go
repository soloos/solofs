package solofssdk

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgPosixFsBase(t *testing.T) {
	var (
		client      Client
		fsINodeMeta solofstypes.FsINodeMeta
		netBlockCap = 1024 * 1024 * 8
		memBlockCap = 1024 * 1024 * 2
		code        fsapi.Status
		err         error
	)
	MakeClientForTest(&client)
	var posixFs = &client.PosixFs

	code = posixFs.SimpleMkdir(&fsINodeMeta, nil, solofstypes.RootFsINodeID, 0777, "test", 0, 0, solofstypes.FS_RDEV)
	if code != fsapi.OK {
		util.AssertTrue(code == solofstypes.FS_EEXIST)
	}

	util.Ignore(fsINodeMeta)
	fsINodeMeta, err = posixFs.SimpleOpenFile("/test/hi", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFs.SimpleOpenFile("/test/hi2", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFs.SimpleOpenFile("/test/hi3", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFs.SimpleOpenFile("/test/hi4", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINodeMeta, err = posixFs.SimpleOpenFile("/test/hi5", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	err = posixFs.DeleteFsINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = posixFs.RenameWithFullPath("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = posixFs.RenameWithFullPath("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = posixFs.ListFsINodeByParentPath("/test", true,
		func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount), uint64(0)
		},
		func(fsINodeMeta solofstypes.FsINodeMeta) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = posixFs.SimpleOpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err.Error(), solofstypes.ErrObjectNotExists.Error())
}
