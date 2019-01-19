package libsdfs

import (
	"soloos/sdfs/types"
	"soloos/util"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
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

	err = client.DirTreeDriver.Mkdir(client.DirTreeDriver.AllocFsINodeID(),
		&fuse.MkdirIn{
			InHeader: fuse.InHeader{
				Length: 0,
				Opcode: 0,
				Unique: 0,
				NodeId: types.RootFsINodeID,
				Context: fuse.Context{
					Owner: fuse.Owner{
						Uid: 0,
						Gid: 0,
					},
					Pid: 0,
				},
				Padding: 0,
			},
			Mode:  0,
			Umask: 0,
		}, "test", &fuse.EntryOut{})
	if err != types.ErrObjectExists {
		assert.NoError(t, err)
	}

	util.Ignore(fsINode)
	fsINode, err = client.DirTreeDriver.OpenFile("/test/hi", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.DirTreeDriver.OpenFile("/test/hi2", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.DirTreeDriver.OpenFile("/test/hi3", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.DirTreeDriver.OpenFile("/test/hi4", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	fsINode, err = client.DirTreeDriver.OpenFile("/test/hi5", netBlockCap, memBlockCap)
	assert.NoError(t, err)
	err = client.DirTreeDriver.DeleteINodeByPath("/test/hi4")
	assert.NoError(t, err)

	err = client.DirTreeDriver.Rename("/test/hi5", "/testhi5")
	assert.NoError(t, err)
	err = client.DirTreeDriver.Rename("/testhi5", "/test/hi5")
	assert.NoError(t, err)

	err = client.DirTreeDriver.ListFsINodeByParentPath("/test",
		func(resultCount int) bool {
			return true
		},
		func(fsINode types.FsINode) bool {
			return true
		})
	assert.NoError(t, err)

	_, err = client.DirTreeDriver.OpenFile("/noexists/hi5", netBlockCap, memBlockCap)
	assert.Equal(t, err, types.ErrObjectNotExists)
}
