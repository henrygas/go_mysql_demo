package gosql

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	insertCount = 1000
	deleteCount = 1000
	selectCount = 1000
	updateCount = 1000
)

// test insert
func TestInsert(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Fatal("failed to NewSqlDemo(), err: ", err)
	}

	pureTruncate(sd)
	// insert1
	timeStart := time.Now()
	var succCount int64 = 0
	for i := 1; i <= insertCount; i++ {
		n, err := pureInsert1(sd, i)
		if err != nil {
			continue
		}
		succCount += n
	}
	timeEnd := time.Now()
	fmt.Printf("insert1 playCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		insertCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// insert2
	pureTruncate(sd)
	timeStart = time.Now()
	succCount = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert2(sd, i); err != nil {
			continue
		} else {
			succCount += n
		}

	}
	timeEnd = time.Now()
	fmt.Printf("insert2 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// insert3
	pureTruncate(sd)
	timeStart = time.Now()
	stmtInsert, err := sd.db.Prepare(insertSql)
	if err != nil {
		return
	}
	defer func() {
		_ = stmtInsert.Close()
	}()
	succCount = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert3(stmtInsert, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd = time.Now()
	fmt.Printf("insert3 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// insert4
	pureTruncate(sd)
	timeStart = time.Now()
	tx, err := sd.db.Begin()
	if err != nil {
		return
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

	succCount = 0
	for i := 1; i <= insertCount; i++ {
		if n, err := pureInsert4(tx, i); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd = time.Now()
	fmt.Printf("insert4 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
}

// test delete
func TestDelete(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	err = initData(sd, deleteCount)
	if err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// delete1
	timeStart := time.Now()
	var succCount int64 = 0
	for i := 1; i <= deleteCount; i++ {
		result, err := sd.db.Exec(deleteSql, i)
		if err != nil {
			continue
		}

		if n, err := result.RowsAffected(); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd := time.Now()
	fmt.Printf("delete1 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// reinitData
	err = initData(sd, deleteCount)
	if err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// delete2
	timeStart = time.Now()
	succCount = 0
	for i := 1; i <= deleteCount; i++ {
		stmt, err := sd.db.Prepare(deleteSql)
		if err != nil {
			continue
		}
		result, err := stmt.Exec(i)
		if err != nil {
			continue
		}

		if n, err := result.RowsAffected(); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd = time.Now()
	fmt.Printf("delete2 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// reinitData
	err = initData(sd, deleteCount)
	if err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// delete3
	timeStart = time.Now()
	succCount = 0
	if stmt, err := sd.db.Prepare(deleteSql); err == nil {
		for i := 1; i <= deleteCount; i++ {
			result, err := stmt.Exec(i)
			if err != nil {
				continue
			}

			if n, err := result.RowsAffected(); err != nil {
				continue
			} else {
				succCount += n
			}
		}
	} else {
		log.Println("failed to prepare deleteSql, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("delete3 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// reinitData
	err = initData(sd, deleteCount)
	if err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// delete4
	timeStart = time.Now()
	succCount = 0
	if tx, err := sd.db.Begin(); err == nil {
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()

		for i := 1; i <= deleteCount; i++ {
			if result, err := tx.Exec(deleteSql, i); err != nil {
				continue
			} else if n, err := result.RowsAffected(); err != nil {
				continue
			} else {
				succCount += n
			}
		}
	} else {
		log.Println("failed to begin tx, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("delete4 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
}

// test update
func TestUpdate(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	if err := initData(sd, updateCount); err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// update1
	timeStart := time.Now()
	var succCount int64 = 0
	for i := 1; i <= updateCount; i++ {
		if result, err := sd.db.Exec(updateByNumEqSql, i*i+1, i); err != nil {
			continue
		} else if n, err := result.RowsAffected(); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd := time.Now()
	fmt.Printf("update1 planCount: %d, succCount: %d, costTime: %f ms, %f ms/op\n",
		updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)

	// reinitData
	if err := initData(sd, updateCount); err != nil {
		return
	}

	// update2
	timeStart = time.Now()
	succCount = 0
	for i := 1; i <= updateCount; i++ {
		if stmt, err := sd.db.Prepare(updateByNumEqSql); err != nil {
			continue
		} else if result, err := stmt.Exec(i*i+1, i); err != nil {
			continue
		} else if n, err := result.RowsAffected(); err != nil {
			continue
		} else {
			succCount += n
		}
	}
	timeEnd = time.Now()
	fmt.Printf("update2 planCount: %d, succCount: %d, costTime: %f ms, %f ms/op\n",
		updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)

	// reinitData
	if err := initData(sd, updateCount); err != nil {
		return
	}

	// update3
	timeStart = time.Now()
	succCount = 0
	if stmt, err := sd.db.Prepare(updateByNumEqSql); err != nil {
		log.Println("failed to prepare updateByNumEqSql, err: ", err)
	} else {
		for i := 1; i <= updateCount; i++ {
			if result, err := stmt.Exec(i*i+1, i); err != nil {
				continue
			} else if n, err := result.RowsAffected(); err != nil {
				continue
			} else {
				succCount += n
			}
		}
	}
	timeEnd = time.Now()

	fmt.Printf("update3 planCount: %d, succCount: %d, costTime: %f ms, %f ms/op\n",
		updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)

	// reinitData
	if err := initData(sd, updateCount); err != nil {
		return
	}

	// update4
	timeStart = time.Now()
	succCount = 0
	if tx, err := sd.db.Begin(); err != nil {
		log.Println("failed to begin tx, err: ", err)
	} else {
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Println("has found panic: ", p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Println("has found error: ", err)
			} else {
				_ = tx.Commit()
			}
		}()
		if stmt, err := tx.Prepare(updateByNumEqSql); err != nil {
			log.Println("failed to prepare updateByNumEqSql, err: ", err)
		} else {
			for i := 1; i <= updateCount; i++ {
				if result, err := stmt.Exec(i*i+1, i); err != nil {
					continue
				} else if n, err := result.RowsAffected(); err != nil {
					continue
				} else {
					succCount += n
				}
			}
		}
	}
	timeEnd = time.Now()
	fmt.Printf("update4 planCount: %d, succCount: %d, costTime: %f ms, %f ms/op\n",
		updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)
}

// test select
func TestSelect(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	err = initData(sd, selectCount)
	if err != nil {
		log.Println("failed to initData(), err: ", err)
		return
	}

	// select1
	timeStart := time.Now()
	var squareNum int
	succCount := 0
	for i := 1; i <= selectCount; i++ {
		rows, err := sd.db.Query(selectByNumSql, i)
		if err != nil {
			continue
		}
		for rows.Next() {
			if err = rows.Scan(&squareNum); err != nil {
				break
			} else {
				succCount++
			}
		}
	}
	timeEnd := time.Now()
	fmt.Printf("select1 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// select2
	timeStart = time.Now()
	succCount = 0
	for i := 1; i <= selectCount; i++ {
		stmt, err := sd.db.Prepare(selectByNumSql)
		if err != nil {
			continue
		}

		rows, err := stmt.Query(i)
		if err != nil {
			_ = stmt.Close()
			continue
		}
		var squareNum int
		for rows.Next() {
			if err = rows.Scan(&squareNum); err != nil {
				break
			} else {
				succCount++
			}
		}
	}
	timeEnd = time.Now()
	fmt.Printf("select2 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// select3
	timeStart = time.Now()
	succCount = 0
	if stmt, err := sd.db.Prepare(selectByNumSql); err == nil {
		var squareNum int
		for i := 1; i <= selectCount; i++ {
			rows, err := stmt.Query(i)
			if err != nil {
				continue
			}
			for rows.Next() {
				if err = rows.Scan(&squareNum); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	} else {
		log.Println("failed to prepare selectSql, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("select3 planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
}

func pureTruncate(sd *SqlDemo) {
	_, _ = sd.db.Exec(truncateSql)
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

func initData(sd *SqlDemo, count int) error {
	pureTruncate(sd)
	stmtInsert, err := sd.db.Prepare(insertSql)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmtInsert.Close()
	}()
	for i := 1; i <= count; i++ {
		if _, err := pureInsert3(stmtInsert, i); err != nil {
			continue
		}
	}
	return nil
}
