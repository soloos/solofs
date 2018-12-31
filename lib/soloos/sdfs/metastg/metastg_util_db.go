package metastg

import (
	"database/sql"

	"github.com/gocraft/dbr"
)

func (p *MetaStg) FetchAndUpdateMaxID(key string, delta int64) (int64, error) {
	var (
		sess         *dbr.Session
		sqlRows      *sql.Rows
		isNeedInsert bool
		maxid        int64
		err          error
	)

	sess = p.dbConn.NewSession(nil)
	sqlRows, err = sess.Select("maxid").
		From("b_maxid").
		Where("key=?", key).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	isNeedInsert = true
	if sqlRows.Next() {
		err = sqlRows.Scan(&maxid)
		if err != nil {
			goto QUERY_DONE
		}
		isNeedInsert = false
	}

	if sqlRows != nil {
		sqlRows.Close()
	}

	if isNeedInsert {
		_, err = sess.InsertInto("b_maxid").
			Columns("key", "maxid").
			Values(key, maxid).
			Exec()
		if err != nil {
			goto QUERY_DONE
		}
	} else {
		maxid += delta
		_, err = sess.Update("b_maxid").
			Set("maxid", maxid).
			Where("key=?", key).
			Exec()
	}

QUERY_DONE:
	return maxid, err
}
