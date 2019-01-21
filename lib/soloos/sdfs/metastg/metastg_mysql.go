package metastg

import (
	"soloos/log"

	_ "github.com/go-sql-driver/mysql"
)

func (p *MetaStg) installMysqlSchema() error {
	var (
		sqls []string
		err  error
	)

	sqls = prepareNetINodesSqls()
	for _, sql := range sqls {
		_, err = p.dbConn.Exec(sql)
		if err != nil {
			log.Error(err, sql)
			return err
		}
	}

	return nil
}
