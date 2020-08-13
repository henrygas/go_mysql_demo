package gosql

import "database/sql"

const (
	deleteCount = 1000
)

func Delete1(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete1(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Delete2(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete2(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func Delete3(sd *SqlDemo) int64 {
	var succCount int64 = 0
	if stmtdelete, err := sd.db.Prepare(deleteSql); err != nil {
		return succCount
	} else {
		defer func() {
			_ = stmtdelete.Close()
		}()
		for i := 1; i <= deleteCount; i++ {
			if n, err := puredelete3(stmtdelete, i); err != nil {
				continue
			} else {
				succCount += n
			}
		}
		return succCount
	}
}

func Delete4(sd *SqlDemo) int64 {
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

	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}

	return succCount
}

func puredelete1(sd *SqlDemo, i int) (int64, error) {
	result, err := sd.db.Exec(deleteSql, i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func puredelete2(sd *SqlDemo, i int) (int64, error) {
	stmtdelete, err := sd.db.Prepare(deleteSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmtdelete.Close()
	}()

	result, err := stmtdelete.Exec(i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func puredelete3(stmtdelete *sql.Stmt, i int) (int64, error) {
	result, err := stmtdelete.Exec(i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func puredelete4(tx *sql.Tx, i int) (int64, error) {
	result, err := tx.Exec(deleteSql, i)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}
