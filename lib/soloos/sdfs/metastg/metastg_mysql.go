package metastg

import (
	_ "github.com/go-sql-driver/mysql"
)

func (p *MetaStg) InstallMysqlSchema() error {
	var (
		sqls []string
		err  error
	)
	sqls = commonSchemaSqls()
	for _, sql := range sqls {
		_, err = p.dbConn.Exec(sql)
		if err != nil {
			return err
		}
	}

	sqls = baseDataSqls()
	for _, sql := range sqls {
		p.dbConn.Exec(sql)
	}

	return nil
}
