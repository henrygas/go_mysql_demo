package gosqlx

import (
	"github.com/jmoiron/sqlx"
	"strconv"
)

const (
	insertCount = 1000
	//deleteCount = 1000
	//selectCount = 1000
	//updateCount = 1000
)

func insert1(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert1(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func insert2(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert2(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func insert3(sxd *SqlxDemo) int64 {
	stmtInsert, err := sxd.db.Preparex(insertOneSql)
	if err != nil {
		return 0
	}
	defer func() {
		_ = stmtInsert.Close()
	}()
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert3(stmtInsert, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func insert4(sxd *SqlxDemo) int64 {
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

	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}

	return succCount
}

func insert5(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		n, err := pureInsert5(sxd, i)
		if err != nil {
			continue
		}
		succCount += n
	}
	return succCount
}

func insert6(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert6(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func insert7(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	stmtInsertNamed, err := sxd.db.PrepareNamed(insertOneNamedSql)
	if err != nil {
		return succCount
	}
	defer func() {
		_ = stmtInsertNamed.Close()
	}()
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert7(stmtInsertNamed, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func insert8(sxd *SqlxDemo) int64 {
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

	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert8(txx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func pureInsert1(sxd *SqlxDemo, i int) (int64, error) {
	result, err := sxd.db.Exec(insertOneSql, i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert2(sxd *SqlxDemo, i int) (int64, error) {
	stmtInsert, err := sxd.db.Preparex(insertOneSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmtInsert.Close()
	}()

	result, err := stmtInsert.Exec(i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert3(stmtInsert *sqlx.Stmt, i int) (int64, error) {
	result, err := stmtInsert.Exec(i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert4(tx *sqlx.Tx, i int) (int64, error) {
	result, err := tx.Exec(insertOneSql, i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert5(sxd *SqlxDemo, i int) (int64, error) {
	result, err := sxd.db.NamedExec(insertOneNamedSql, map[string]interface{}{
		"name": "江南" + strconv.Itoa(i),
		"age":  i % 20,
	})
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert6(sxd *SqlxDemo, i int) (int64, error) {
	stmtInsert, err := sxd.db.PrepareNamed(insertOneNamedSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmtInsert.Close()
	}()

	result, err := stmtInsert.Exec(map[string]interface{}{
		"name": "江南" + strconv.Itoa(i),
		"age":  i % 20,
	})
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert7(stmtInsert *sqlx.NamedStmt, i int) (int64, error) {
	result, err := stmtInsert.Exec(map[string]interface{}{
		"name": "江南" + strconv.Itoa(i),
		"age":  i % 20,
	})
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func pureInsert8(tx *sqlx.Tx, i int) (int64, error) {
	result, err := tx.NamedExec(insertOneNamedSql, map[string]interface{}{
		"name": "江南" + strconv.Itoa(i),
		"age":  i % 20,
	})
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}
