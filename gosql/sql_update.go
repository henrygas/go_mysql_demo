package gosql

import (
	"database/sql"
)

const (
	updateCount = 1000
)

func Update1(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate1(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Update2(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate2(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Update3(sd *SqlDemo) int64 {
	var succCount int64 = 0
	if stmtupdate, err := sd.db.Prepare(updateSql); err != nil {
		return succCount
	} else {
		defer func() {
			_ = stmtupdate.Close()
		}()
		for i := 1; i <= updateCount; i++ {
			if n, err := pureupdate3(stmtupdate, i); err != nil {
				continue
			} else {
				succCount += n
			}
		}
		return succCount
	}
}

func Update4(sd *SqlDemo) int64 {
	var succCount int64 = 0

	tx, err := sd.db.Begin()
	if err != nil {
		return succCount
	}

	defer func() {
		if p := recover(); p != nil {
			//log.Println("recover from panic to rollback")
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			//log.Println("error happens, about to rollback, err: ", err)
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
			//log.Println("succ to commit!")
		}
	}()

	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}

	return succCount
}

func pureupdate1(sd *SqlDemo, i int) (int64, error) {
	result, err := sd.db.Exec(updateSql, i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureupdate2(sd *SqlDemo, i int) (int64, error) {
	if stmtupdate, err := sd.db.Prepare(updateSql); err != nil {
		return 0, err
	} else {
		defer func() {
			_ = stmtupdate.Close()
		}()

		if result, err := stmtupdate.Exec(i); err != nil {
			return 0, err
		} else if n, err := result.RowsAffected(); err != nil {
			return 0, err
		} else {
			return n, nil
		}
	}
}

func pureupdate3(stmtupdate *sql.Stmt, i int) (int64, error) {
	result, err := stmtupdate.Exec(i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureupdate4(tx *sql.Tx, i int) (int64, error) {
	result, err := tx.Exec(updateSql, i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}
