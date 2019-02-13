package metastg

import (
	"soloos/common/util"
	"soloos/common/util/offheap"
)

func MakeMetaStgForTest(offheapDriver *offheap.OffheapDriver, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
