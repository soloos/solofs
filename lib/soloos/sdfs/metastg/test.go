package metastg

import (
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
)

func MakeMetaStgForTest(soloOSEnv *soloosbase.SoloOSEnv, metaStg *MetaStg) {
	var err error
	err = metaStg.Init(soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
}
