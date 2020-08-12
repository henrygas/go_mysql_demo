package gosql

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	insertCount = 10000
	deleteCount = 10000
	selectCount = 10000
	updateCount = 10000
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
	for i := 1; i <= insertCount; i++ {
		pureInsert1(sd, i)
	}
	timeEnd := time.Now()
	fmt.Printf("insert1 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// insert2
	pureTruncate(sd)
	timeStart = time.Now()
	for i := 1; i <= insertCount; i++ {
		pureInsert2(sd, i)
	}
	timeEnd = time.Now()
	fmt.Printf("insert2 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

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
	for i := 1; i <= insertCount; i++ {
		pureInsert3(stmtInsert, i)
	}
	timeEnd = time.Now()
	fmt.Printf("insert3 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

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

	for i := 1; i <= insertCount; i++ {
		pureInsert4(tx, i)
	}
	timeEnd = time.Now()
	fmt.Printf("insert4 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
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
	for i := 1; i <= deleteCount; i++ {
		_, err = sd.db.Exec(deleteSql, i)
		if err != nil {
			continue
		}
	}
	timeEnd := time.Now()
	fmt.Printf("delete1 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// delete2
	timeStart = time.Now()
	for i := 1; i <= deleteCount; i++ {
		stmt, err := sd.db.Prepare(deleteSql)
		if err != nil {
			continue
		}
		_, err = stmt.Exec(i)
		if err != nil {
			continue
		}
	}
	timeEnd = time.Now()
	fmt.Printf("delete2 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// delete3
	timeStart = time.Now()
	if stmt, err := sd.db.Prepare(deleteSql); err == nil {
		for i := 1; i <= deleteCount; i++ {
			_, err = stmt.Exec(i)
			if err != nil {
				continue
			}
		}
	} else {
		log.Println("failed to prepare deleteSql, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("delete3 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// delete4
	timeStart = time.Now()
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
			_, err = tx.Exec(deleteSql, i)
			if err != nil {
				continue
			}
		}
	} else {
		log.Println("failed to begin tx, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("delete4 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
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
	for i := 1; i <= selectCount; i++ {
		if i == 153 {
			fmt.Println("break here!")
		}
		rows, err := sd.db.Query(selectByNumSql, i)
		if err != nil {
			continue
		}
		for rows.Next() {
			err = rows.Scan(&squareNum)
			if err != nil {
				break
			}
		}
	}
	timeEnd := time.Now()
	fmt.Printf("select1 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// select2
	timeStart = time.Now()
	for i := 1; i <= selectCount; i++ {
		stmt, err := sd.db.Prepare(selectByNumSql)
		if err != nil {
			continue
		}
		_, err = stmt.Query(i)
		if err != nil {
			continue
		}
	}
	timeEnd = time.Now()
	fmt.Printf("select2 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)

	// select3
	timeStart = time.Now()
	if stmt, err := sd.db.Prepare(selectByNumSql); err == nil {
		for i := 1; i <= selectCount; i++ {
			_, err = stmt.Query(i)
			if err != nil {
				continue
			}
		}
	} else {
		log.Println("failed to prepare selectSql, err: ", err)
	}
	timeEnd = time.Now()
	fmt.Printf("select3 cost Time: %f ms, %f ms/op\n",
		timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
}

func pureTruncate(sd *SqlDemo) {
	_, _ = sd.db.Exec(truncateSql)
}

func pureInsert1(sd *SqlDemo, i int) {
	_, _ = sd.db.Exec(insertSql, i, i*i)
}

func pureInsert2(sd *SqlDemo, i int) {
	stmtInsert, err := sd.db.Prepare(insertSql)
	if err != nil {
		return
	}
	defer func() {
		_ = stmtInsert.Close()
	}()

	_, err = stmtInsert.Exec(i, i*i)
	if err != nil {
		return
	}
}

func pureInsert3(stmtInsert *sql.Stmt, i int) {
	_, err := stmtInsert.Exec(i, i*i)
	if err != nil {
		return
	}
}

func pureInsert4(tx *sql.Tx, i int) {
	_, err := tx.Exec(insertSql, i, i*i)
	if err != nil {
		return
	}
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
		pureInsert3(stmtInsert, i)
	}
	return nil
}
