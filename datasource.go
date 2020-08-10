package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var datasourceName = "root:1234@tcp(127.0.0.1:3306)/ys?timeout=90s&collation=utf8mb4_bin"

func testOpen() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("fail to open mysql driver, ", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("fail to ping mysql, ", err)
	} else {
		log.Println("succ to ping mysql!")
	}
}

func testInsert() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("failed to open mysql, ", err)
	}
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO square_num(num, square_num) VALUES (?, ?)")
	if err != nil {
		log.Fatal("failed to prepare insert, ", err)
	}
	defer stmtInsert.Close()

	for i := 1; i <= 10; i++ {
		result, err := stmtInsert.Exec(i, i*i)
		if err != nil {
			log.Fatal("failed to exec stmtInsert, ", err)
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Fatal("failed to get lastInsertId, ", err)
		}
		fmt.Println("lastInsertID is: ", lastInsertID)

		n, err := result.RowsAffected()
		if err != nil {
			log.Fatal("failed to get RowsAffected, ", err)
		}
		fmt.Println("insert finish, affected rows: ", n)
	}
}

func testDelete() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("failed to open mysql, ", err)
	}
	defer db.Close()

	stmtDelete, err := db.Prepare("DELETE FROM square_num WHERE num = ?")
	if err != nil {
		log.Fatal("failed to prepare stmtDelete, ", err)
	}
	defer stmtDelete.Close()

	for i := 1; i <= 10; i++ {
		result, err := stmtDelete.Exec(i)
		if err != nil {
			log.Fatal("failed to exec delete, ", err)
		}
		n, err := result.RowsAffected()
		if err != nil {
			log.Fatal("failed to get rows affected, ", err)
		}
		fmt.Printf("delete finish, affected rows: %d\n", n)
	}
}

func testQuery() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("failed to open mysql, ", err)
	}
	defer db.Close()

	stmtQuery, err := db.Prepare("SELECT square_num FROM square_num WHERE num = ?")
	if err != nil {
		log.Fatal("failed to prepare stmtQuery, ", err)
	}
	defer stmtQuery.Close()

	var value int
	for i := 1; i <= 10; i++ {
		err := stmtQuery.QueryRow(i).Scan(&value)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("num=%d, query result=no rows\n", i)
			} else {
				log.Fatal(fmt.Sprintf("failed to query row with num = %d, ", i), err)
			}
		} else {
			fmt.Printf("num=%d, query result=%d\n", i, value)
		}
	}
}

func testUpdate() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("failed to open mysql, ", err)
	}
	defer db.Close()

	stmtUpdate, err := db.Prepare("UPDATE square_num SET square_num = ? WHERE num = ?")
	if err != nil {
		log.Fatal("failed to prepare update, ", err)
	}
	defer stmtUpdate.Close()

	for i := 1; i <= 10; i++ {
		result, err := stmtUpdate.Exec(i*2, i)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to exec update, num=%d, ", i), err)
		}
		n, err := result.RowsAffected()
		if err != nil {
			log.Fatal("failed to get rows affected, ", err)
		}
		fmt.Printf("num=%d, update affected rows=%d\n", i, n)
	}
}

func testTransaction() {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal("failed to open mysql, ", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("failed to open transaction, ", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil {
			log.Fatal("failed to rollback transaction, ", err)
		}
	}(tx)

	result, err := tx.Exec("UPDATE square_num SET square_num = square_num - 1 WHERE num < 11")
	if err != nil {
		log.Fatal("failed to update square_num=square_num-1, ", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Fatal("failed to get rows affected, ", err)
	}
	fmt.Println("update square_num=square_num-1 affected rows: ", n)

	result, err = tx.Exec("UPDATE square_num SET square_num = square_num + 5 WHERE num < 11")
	if err != nil {
		log.Fatal("failed to update square_num=square_num+5, ", err)
	}
	n, err = result.RowsAffected()
	if err != nil {
		log.Fatal("failed to get rows affected, ", err)
	}
	fmt.Println("update square_num=square_num+5 affected rows: ", n)

	func() {
		panic("running error happens!")
	}()

	if err := tx.Commit(); err != nil {
		log.Fatal("failed to commit transaction, ", err)
	} else {
		fmt.Println("succ to commit transaction!")
	}
}

func main() {
	//testOpen()
	//testInsert()
	//testDelete()
	//testQuery()
	//testUpdate()
	//testTransaction()
}
