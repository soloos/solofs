package libsdfs

import (
	"soloos/sdfs/types"
	"soloos/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgDirTreeDriverBase(t *testing.T) {
	var (
		client      Client
		fsINode     types.FsINode
		netBlockCap = types.DefaultNetBlockCap
		memBlockCap = types.DefaultMemBlockCap
		err         error
	)
	MakeClientForTest(&client)

	fsINode, err = client.MetaStg.DirTreeDriver.Mkdir("/test")
	util.Ignore(fsINode)

	fsINode, err = client.MetaStg.DirTreeDriver.OpenFile("/test/hi", netBlockCap, memBlockCap)
	fsINode, err = client.MetaStg.DirTreeDriver.OpenFile("/test/hi2", netBlockCap, memBlockCap)
	fsINode, err = client.MetaStg.DirTreeDriver.OpenFile("/test/hi3", netBlockCap, memBlockCap)
	fsINode, err = client.MetaStg.DirTreeDriver.OpenFile("/test/hi4", netBlockCap, memBlockCap)
	fsINode, err = client.MetaStg.DirTreeDriver.OpenFile("/test/hi5", netBlockCap, memBlockCap)
	err = client.MetaStg.DirTreeDriver.DeleteINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = client.MetaStg.DirTreeDriver.Rename("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = client.MetaStg.DirTreeDriver.Rename("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = client.MetaStg.DirTreeDriver.ListFsINodeByParentPath("/test",
		func(resultCount int) bool {
			return true
		},
		func(fsINode types.FsINode) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = client.MetaStg.DirTreeDriver.OpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err, types.ErrObjectNotExists)
}
