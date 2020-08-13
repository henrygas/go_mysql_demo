package gosqlx

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// test insert
func TestInsert(t *testing.T) {
	sxd, err := NewSqlxDemo()
	if err != nil {
		log.Fatal("failed to NewSqlxDemo(), err: ", err)
	}

	insertFuncs := [8]func(sxd *SqlxDemo) int64{insert1, insert2, insert3, insert4, insert5, insert6, insert7, insert8}
	for funcSeq := 1; funcSeq <= len(insertFuncs); funcSeq++ {
		insertMethod(insertFuncs[funcSeq-1], funcSeq, sxd)
	}
}

// test delete
func TestDelete(t *testing.T) {
	sxd, err := NewSqlxDemo()
	if err != nil {
		log.Println("failed to NewSqlxDemo(), err: ", err)
		return
	}

	deleteFuncs := [8]func(sxd *SqlxDemo) int64{delete1, delete2, delete3, delete4, delete5, delete6, delete7, delete8}
	for funcSeq := 1; funcSeq <= len(deleteFuncs); funcSeq++ {
		deleteMethod(deleteFuncs[funcSeq-1], funcSeq, sxd)
	}
}

// test update
func TestUpdate(t *testing.T) {
	sxd, err := NewSqlxDemo()
	if err != nil {
		log.Println("failed to NewSqlxDemo(), err: ", err)
		return
	}

	updateFuncs := [8]func(sxd *SqlxDemo) int64{update1, update2, update3, update4, update5, update6, update7, update8}
	for funcSeq := 1; funcSeq <= len(updateFuncs); funcSeq++ {
		updateMethod(updateFuncs[funcSeq-1], funcSeq, sxd)
	}
}

// test select
func TestSelect(t *testing.T) {
	sxd, err := NewSqlxDemo()
	if err != nil {
		log.Println("failed to NewSqlxDemo(), err: ", err)
		return
	}

	selectFuncs := [6]func(sxd *SqlxDemo) int64{select1, select2, select3, select4, select5, select6}
	for funcSeq := 1; funcSeq <= len(selectFuncs); funcSeq++ {
		selectMethod(selectFuncs[funcSeq-1], funcSeq, sxd)
	}
}

func pureTruncate(sxd *SqlxDemo) {
	_, _ = sxd.db.Exec(truncateSql)
}

func insertMethod(insertFunc func(sxd *SqlxDemo) int64, funcSeq int, sxd *SqlxDemo) {
	pureTruncate(sxd)
	timeStart := time.Now()
	succCount := insertFunc(sxd)
	timeEnd := time.Now()
	fmt.Printf("insert%d planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		funcSeq, insertCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
}

func deleteMethod(deleteFunc func(sxd *SqlxDemo) int64, funcSeq int, sxd *SqlxDemo) {
	initData(sxd, deleteCount)
	timeStart := time.Now()
	succCount := deleteFunc(sxd)
	timeEnd := time.Now()
	fmt.Printf("delete%d planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		funcSeq, deleteCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/deleteCount)
}

func updateMethod(updateFunc func(sxd *SqlxDemo) int64, funcSeq int, sxd *SqlxDemo) {
	initData(sxd, updateCount)
	timeStart := time.Now()
	succCount := updateFunc(sxd)
	timeEnd := time.Now()
	fmt.Printf("update%d planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		funcSeq, updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)
}

func selectMethod(selectFunc func(sxd *SqlxDemo) int64, funcSeq int, sxd *SqlxDemo) {
	initData(sxd, selectCount)
	timeStart := time.Now()
	succCount := selectFunc(sxd)
	timeEnd := time.Now()
	fmt.Printf("update%d planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
		funcSeq, updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)
}

func initData(sxd *SqlxDemo, count int) error {
	pureTruncate(sxd)
	stmtInsert, err := sxd.db.Preparex(insertOneSql)
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
