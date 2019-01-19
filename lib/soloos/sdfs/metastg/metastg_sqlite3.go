package metastg

import (
	"soloos/dbcli"
	"soloos/log"

	_ "github.com/mattn/go-sqlite3"
)

func InstallSqlite3Schema(dbConn *dbcli.Connection) error {
	var (
		sqls []string
		err  error
	)
	sqls = prepareNetINodesSqls()
	for _, sql := range sqls {
		_, err = dbConn.Exec(sql)
		if err != nil {
			log.Error(err, sql)
			return err
		}
	}

	sqls = prepareDirTreeSqls()
	for _, sql := range sqls {
		_, err = dbConn.Exec(sql)
		if err != nil {
			log.Error(err, sql)
		}
	}

	return nil
}
