package gosql

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// test insert
func TestInsert(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Fatal("failed to NewSqlDemo(), err: ", err)
	}

	inserts := [4]func(sd *SqlDemo) int64{Insert1, Insert2, Insert3, Insert4}
	for i := 0; i < len(inserts); i++ {
		pureTruncate(sd)
		timeStart := time.Now()
		succCount := inserts[i](sd)
		timeEnd := time.Now()
		fmt.Printf("insert%d playCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
			i+1, insertCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
	}
}

// test delete
func TestDelete(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	deletes := [4]func(sd *SqlDemo) int64{Delete1, Delete2, Delete3, Delete4}
	for i := 0; i < len(deletes); i++ {
		if err := initData(sd, deleteCount); err != nil {
			continue
		}
		timeStart := time.Now()
		succCount := deletes[i](sd)
		timeEnd := time.Now()
		fmt.Printf("delete%d playCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
			i+1, deleteCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
	}

}

// test update
func TestUpdate(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	updates := [4]func(sd *SqlDemo) int64{Update1, Update2, Update3, Update4}
	for i := 0; i < len(updates); i++ {
		if err := initData(sd, updateCount); err != nil {
			continue
		}

		timeStart := time.Now()
		succCount := updates[i](sd)
		timeEnd := time.Now()
		fmt.Printf("update%d planCount: %d, succCount: %d, costTime: %f ms, %f ms/op\n",
			i+1, updateCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/updateCount)
	}
}

// test select
func TestSelect(t *testing.T) {
	sd, err := NewSqlDemo()
	if err != nil {
		log.Println("failed to NewSqlDemo(), err: ", err)
		return
	}

	selects := [3]func(sd *SqlDemo) int64{Select1, Select2, Select3}
	for i := 0; i < len(selects); i++ {
		if err = initData(sd, selectCount); err != nil {
			continue
		}

		timeStart := time.Now()
		succCount := selects[i](sd)
		timeEnd := time.Now()
		fmt.Printf("select%d planCount: %d, succCount: %d, cost Time: %f ms, %f ms/op\n",
			i+1, selectCount, succCount, timeEnd.Sub(timeStart).Seconds()*1000, timeEnd.Sub(timeStart).Seconds()*1000/insertCount)
	}
}

func pureTruncate(sd *SqlDemo) {
	_, _ = sd.db.Exec(truncateSql)
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
