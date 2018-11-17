package metastg

import (
	_ "github.com/mattn/go-sqlite3"
)

func (p *MetaStg) InstallSqlite3Schema() error {
	var (
		sql string
		err error
	)
	sql = commonSchemaSql()
	_, err = p.DBConn.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
