package gosqlx

import (
	"github.com/jmoiron/sqlx"
)

const (
	updateCount = 1000
)

func update1(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate1(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update2(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate2(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update3(sxd *SqlxDemo) int64 {
	stmtupdate, err := sxd.db.Preparex(updateOneSql)
	if err != nil {
		return 0
	}
	defer func() {
		_ = stmtupdate.Close()
	}()
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate3(stmtupdate, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update4(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	tx, err := sxd.db.Beginx()
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

func update5(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate5(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update6(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate6(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update7(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	stmtupdateNamed, err := sxd.db.PrepareNamed(updateOneNamedSql)
	if err != nil {
		return succCount
	}
	defer func() {
		_ = stmtupdateNamed.Close()
	}()
	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate7(stmtupdateNamed, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func update8(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	txx, err := sxd.db.Beginx()
	if err != nil {
		return succCount
	}

	defer func() {
		if p := recover(); p != nil {
			//log.Println("recover from panic to rollback")
			_ = txx.Rollback()
			panic(p)
		} else if err != nil {
			//log.Println("error happens, about to rollback, err: ", err)
			_ = txx.Rollback()
		} else {
			_ = txx.Commit()
			//log.Println("succ to commit!")
		}
	}()

	for i := 1; i <= updateCount; i++ {
		if n, err := pureupdate8(txx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func pureupdate1(sxd *SqlxDemo, i int) (int64, error) {
	result, err := sxd.db.Exec(updateOneSql, i)
	if err != nil {
		return 0, err
	}

	if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func pureupdate2(sxd *SqlxDemo, i int) (int64, error) {
	stmtupdate, err := sxd.db.Preparex(updateOneSql)
	if err != nil {
		return 0, err
	}
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

func pureupdate3(stmtupdate *sqlx.Stmt, i int) (int64, error) {
	if result, err := stmtupdate.Exec(i); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func pureupdate4(tx *sqlx.Tx, i int) (int64, error) {
	if result, err := tx.Exec(updateOneSql, i); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func pureupdate5(sxd *SqlxDemo, i int) (int64, error) {
	if result, err := sxd.db.NamedExec(updateOneNamedSql, map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func pureupdate6(sxd *SqlxDemo, i int) (int64, error) {
	if stmtupdate, err := sxd.db.PrepareNamed(updateOneNamedSql); err != nil {
		return 0, err
	} else {
		defer func() {
			_ = stmtupdate.Close()
		}()

		if result, err := stmtupdate.Exec(map[string]interface{}{
			"id": i,
		}); err != nil {
			return 0, err
		} else if n, err := result.RowsAffected(); err != nil {
			return 0, err
		} else {
			return n, nil
		}
	}
}

func pureupdate7(stmtupdate *sqlx.NamedStmt, i int) (int64, error) {
	if result, err := stmtupdate.Exec(map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func pureupdate8(tx *sqlx.Tx, i int) (int64, error) {
	if result, err := tx.NamedExec(updateOneNamedSql, map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}
