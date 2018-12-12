package metastg

import (
	"soloos/util"
	"soloos/util/offheap"
	"testing"
)

func InitMetaStgForTest(t *testing.T, offheapDriver *offheap.OffheapDriver, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
