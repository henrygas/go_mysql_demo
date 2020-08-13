package gosqlx

import (
	"github.com/jmoiron/sqlx"
)

const (
	deleteCount = 1000
)

func delete1(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete1(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete2(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete2(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete3(sxd *SqlxDemo) int64 {
	stmtdelete, err := sxd.db.Preparex(deleteOneSql)
	if err != nil {
		return 0
	}
	defer func() {
		_ = stmtdelete.Close()
	}()
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete3(stmtdelete, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete4(sxd *SqlxDemo) int64 {
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

	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}

	return succCount
}

func delete5(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete5(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete6(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete6(sxd, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete7(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	stmtdeleteNamed, err := sxd.db.PrepareNamed(deleteOneNamedSql)
	if err != nil {
		return succCount
	}
	defer func() {
		_ = stmtdeleteNamed.Close()
	}()
	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete7(stmtdeleteNamed, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func delete8(sxd *SqlxDemo) int64 {
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

	for i := 1; i <= deleteCount; i++ {
		if n, err := puredelete8(txx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	return succCount
}

func puredelete1(sxd *SqlxDemo, i int) (int64, error) {
	result, err := sxd.db.Exec(deleteOneSql, i)
	if err != nil {
		return 0, err
	}

	if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete2(sxd *SqlxDemo, i int) (int64, error) {
	stmtdelete, err := sxd.db.Preparex(deleteOneSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmtdelete.Close()
	}()

	if result, err := stmtdelete.Exec(i); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete3(stmtdelete *sqlx.Stmt, i int) (int64, error) {
	if result, err := stmtdelete.Exec(i); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete4(tx *sqlx.Tx, i int) (int64, error) {
	if result, err := tx.Exec(deleteOneSql, i); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete5(sxd *SqlxDemo, i int) (int64, error) {
	if result, err := sxd.db.NamedExec(deleteOneNamedSql, map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete6(sxd *SqlxDemo, i int) (int64, error) {
	if stmtdelete, err := sxd.db.PrepareNamed(deleteOneNamedSql); err != nil {
		return 0, err
	} else {
		defer func() {
			_ = stmtdelete.Close()
		}()

		if result, err := stmtdelete.Exec(map[string]interface{}{
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

func puredelete7(stmtdelete *sqlx.NamedStmt, i int) (int64, error) {
	if result, err := stmtdelete.Exec(map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func puredelete8(tx *sqlx.Tx, i int) (int64, error) {
	if result, err := tx.NamedExec(deleteOneNamedSql, map[string]interface{}{
		"id": i,
	}); err != nil {
		return 0, err
	} else if n, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}
