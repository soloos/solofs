package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/sdfs/types"
)

func FetchAndUpdateMaxID(dbConn *sdbapi.Connection, key string, delta types.FsINodeID) (types.FsINodeID, error) {
	var (
		sess         sdbapi.Session
		sqlRows      *sql.Rows
		isNeedInsert bool
		maxid        types.FsINodeID
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
