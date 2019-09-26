package metastg

import (
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func MakeMetaStgForTest(soloosEnv *soloosbase.SoloosEnv, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(soloosEnv, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
