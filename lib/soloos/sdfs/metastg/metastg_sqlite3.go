package metastg

import (
	_ "github.com/mattn/go-sqlite3"
)

func (p *MetaStg) InstallSqlite3Schema() error {
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

	return nil
}
