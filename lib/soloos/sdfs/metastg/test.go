package metastg

import (
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func MakeMetaStgForTest(soloOSEnv *soloosbase.SoloOSEnv, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
