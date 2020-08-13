package gosql

import "database/sql"

const (
	insertCount = 1000
)

func Insert1(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert1(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Insert2(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert2(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Insert3(sd *SqlDemo) int64 {
	var succCount int64 = 0
	if stmtInsert, err := sd.db.Prepare(insertSql); err != nil {
		return succCount
	} else {
		defer func() {
			_ = stmtInsert.Close()
		}()
		for i := 1; i <= insertCount; i++ {
			if n, err := pureInsert3(stmtInsert, i); err != nil {
				continue
			} else {
				succCount += n
			}
		}
		return succCount
	}
}

func Insert4(sd *SqlDemo) int64 {
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

	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}

	return succCount
}

func pureInsert1(sd *SqlDemo, i int) (int64, error) {
	result, err := sd.db.Exec(insertSql, i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	sd.db.Exec("select name, age from user  where id = :id", sql.Named("id", 1))

	return n, nil
}

func pureInsert2(sd *SqlDemo, i int) (int64, error) {
	stmtInsert, err := sd.db.Prepare(insertSql)
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

func pureInsert3(stmtInsert *sql.Stmt, i int) (int64, error) {
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

func pureInsert4(tx *sql.Tx, i int) (int64, error) {
	result, err := tx.Exec(insertSql, i, i*i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}
