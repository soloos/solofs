package metastg

import (
	"soloos/util"
	"soloos/util/offheap"
)

func MakeMetaStgForTest(offheapDriver *offheap.OffheapDriver, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
