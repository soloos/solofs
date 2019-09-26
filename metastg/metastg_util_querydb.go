package metastg

import (
	"database/sql"
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
)

func FetchAndUpdateMaxID(dbConn *solodbapi.Connection, key string, delta solofsapitypes.FsINodeID) (solofsapitypes.FsINodeID, error) {
	var (
		sess         solodbapi.Session
		sqlRows      *sql.Rows
		isNeedInsert bool
		maxid        solofsapitypes.FsINodeID
		err          error
	)

	err = dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("maxid").
		From("b_maxid").
		Where("mkey=?", key).Rows()
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
			Columns("mkey", "maxid").
			Values(key, maxid).
			Exec()
		if err != nil {
			goto QUERY_DONE
		}
	} else {
		maxid += delta
		_, err = sess.Update("b_maxid").
			Set("maxid", maxid).
			Where("mkey=?", key).
			Exec()
	}

QUERY_DONE:
	return maxid, err
}
