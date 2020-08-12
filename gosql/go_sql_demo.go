package gosql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	datasourceName = "root:1234@tcp(127.0.0.1:3306)/ys?timeout=90s&collation=utf8mb4_bin"
	//datasourceName   = "root:Jingle@100@tcp(10.21.248.251:3306)/ys?timeout=90s&collation=utf8mb4_bin"
	truncateSql      = "truncate table `square_num`"
	insertSql        = "INSERT INTO `square_num`(num, square_num) VALUES (?, ?)"
	deleteSql        = "DELETE FROM `square_num` WHERE `num` = ?"
	selectByNumSql   = "SELECT square_num FROM square_num WHERE num = ?"
	updateByNumEqSql = "UPDATE square_num SET square_num = ? WHERE num = ?"
	updateByNumLtSql = "UPDATE square_num SET square_num = square_num - 1 WHERE num < 3"
	updateByNumGtSql = "UPDATE square_num SET square_num = square_num + 5 WHERE num > 2"
)

type SqlDemo struct {
	db *sql.DB
}

type SquareNum struct {
	Num       int `db:"num"`
	SquareNum int `db:"square_num"`
}

func NewSqlDemo() (*SqlDemo, error) {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Println("fail to open mysql driver, ", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		log.Println("fail to ping mysql, ", err)
		return nil, err
	}

	sd := SqlDemo{
		db: db,
	}
	return &sd, nil
}

func (sd *SqlDemo) InsertDemo() error {
	stmtInsert, err := sd.db.Prepare(insertSql)
	if err != nil {
		log.Println("failed to prepare insert, ", err)
		return err
	}
	defer func() {
		_ = stmtInsert.Close()
	}()

	var num int64
	for i := 1; i <= 10; i++ {
		result, err := stmtInsert.Exec(i, i*i)
		if err != nil {
			log.Println("failed to exec stmtInsert, ", err)
			continue
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Println("failed to get lastInsertId, ", err)
			continue
		}
		fmt.Println("lastInsertID is: ", lastInsertID)

		n, err := result.RowsAffected()
		if err != nil {
			log.Println("failed to get RowsAffected, ", err)
			continue
		}
		num += n
	}
	fmt.Println("insert finish, affected rows: ", num)
	fmt.Println("================================================")

	return nil
}

func (sd *SqlDemo) DeleteDemo() error {
	stmtDelete, err := sd.db.Prepare(deleteSql)
	if err != nil {
		log.Println("failed to prepare stmtDelete, ", err)
		return err
	}
	defer func() {
		_ = stmtDelete.Close()
	}()

	var num int64
	for i := 1; i <= 5; i++ {
		result, err := stmtDelete.Exec(i)
		if err != nil {
			log.Println("failed to exec delete, ", err)
			continue
		}
		n, err := result.RowsAffected()
		if err != nil {
			log.Println("failed to get rows affected, ", err)
			continue
		}
		num += n
	}
	fmt.Printf("delete finish, affected rows: %d\n", num)
	fmt.Println("================================================")

	return nil
}

func (sd *SqlDemo) TruncateDemo() error {
	result, err := sd.db.Exec(truncateSql)
	if err != nil {
		log.Println("failed to truncate table square_num, err: ", err)
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rows affected, err: ", err)
		return err
	}
	fmt.Printf("truncate table finish, rows affected=%d\n", n)
	fmt.Println("================================================")

	return nil
}

func (sd *SqlDemo) QueryDemo() error {
	stmtQuery, err := sd.db.Prepare(selectByNumSql)
	if err != nil {
		log.Println("failed to prepare stmtQuery, ", err)
		return err
	}
	defer func() {
		_ = stmtQuery.Close()
	}()

	var value int
	for i := 1; i <= 10; i++ {
		err := stmtQuery.QueryRow(i).Scan(&value)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("num=%d, query result=no rows\n", i)
			} else {
				log.Println(fmt.Sprintf("failed to query row with num = %d, ", i), err)
				continue
			}
		} else {
			fmt.Printf("num=%d, query result=%d\n", i, value)
		}
	}
	fmt.Println("================================================")

	return nil
}

func (sd *SqlDemo) UpdateDemo() error {
	stmtUpdate, err := sd.db.Prepare(updateByNumEqSql)
	if err != nil {
		log.Println("failed to prepare update, ", err)
		return err
	}
	defer func() {
		_ = stmtUpdate.Close()
	}()

	var num int64
	for i := 1; i <= 10; i++ {
		result, err := stmtUpdate.Exec(i*2, i)
		if err != nil {
			log.Println(fmt.Sprintf("failed to exec update, num=%d, ", i), err)
			continue
		}
		n, err := result.RowsAffected()
		if err != nil {
			log.Println("failed to get rows affected, ", err)
			continue
		}
		num += n
	}
	fmt.Printf("update finish, affected rows=%d\n", num)
	fmt.Println("================================================")

	return nil
}

func (sd *SqlDemo) TransactionDemo() error {
	tx, err := sd.db.Begin()
	if err != nil {
		log.Println("failed to open transaction, ", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Println("rollback, ", err)
			_ = tx.Rollback()
		} else {
			if err = tx.Commit(); err != nil {
				log.Println("failed to commit, err: ", err)
			} else {
				log.Println("succ to commit")
			}
		}
	}()

	result, err := tx.Exec(updateByNumLtSql)
	if err != nil {
		log.Println("failed to update square_num=square_num-1, ", err)
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rows affected, ", err)
		return err
	}
	fmt.Println("update square_num=square_num-1 affected rows: ", n)

	result, err = tx.Exec(updateByNumGtSql)
	if err != nil {
		log.Println("failed to update square_num=square_num+5, ", err)
		return err
	}
	n, err = result.RowsAffected()
	if err != nil {
		log.Println("failed to get rows affected, ", err)
		return err
	}
	fmt.Println("update square_num=square_num+5 affected rows: ", n)

	func() {
		panic("running error happens!")
	}()

	return nil
}
